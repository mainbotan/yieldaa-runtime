package preset

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"os"
	"sync"
	"time"

	"github.com/ghodss/yaml"
)

func ProcessEntities(files []EntityFile, maxWorkers int) ([]ProcessedEntity, []error) {
	if len(files) == 0 {
		return []ProcessedEntity{}, nil
	}

	if maxWorkers <= 0 {
		maxWorkers = 10
	}
	if maxWorkers > len(files) {
		maxWorkers = len(files)
	}

	progress := NewProgressTracker(len(files))
	jobs := make(chan EntityFile, len(files))
	results := make(chan ProcessedEntity, len(files))
	errors := make(chan error, len(files))

	var wg sync.WaitGroup
	var seenHashes sync.Map // thread-safe access

	// Workers
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for file := range jobs {
				progress.StartJob()

				// content hash before entity validation
				content, err := os.ReadFile(file.Path)
				if err != nil {
					errors <- fmt.Errorf("%s: read: %w", file.Path, err)
					progress.CompleteJob()
					continue
				}

				contentHash := crc32.ChecksumIEEE(content)

				// checking hash unique
				if _, alreadyProcessed := seenHashes.Load(contentHash); alreadyProcessed {
					// skip duplicate
					progress.CompleteJob()
					continue
				}

				// set as processed
				seenHashes.Store(contentHash, true)

				// total file process
				result := processEntityWithContent(file, content, contentHash)

				if result.FatalError != nil {
					errors <- fmt.Errorf("%s: %w", file.Path, result.FatalError)
				} else {
					results <- result
				}

				progress.CompleteJob()
			}
		}(i)
	}

	// Progress monitor
	stopProgress := make(chan bool)
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-stopProgress:
				return
			case <-ticker.C:
				progress.print()
			}
		}
	}()

	// Feed jobs
	go func() {
		for _, file := range files {
			jobs <- file
		}
		close(jobs)
	}()

	// Wait completion
	go func() {
		wg.Wait()
		close(results)
		close(errors)
		stopProgress <- true
		progress.Finish()
	}()

	// Collect results
	var processed []ProcessedEntity
	var fatalErrors []error

	for result := range results {
		processed = append(processed, result)
	}
	for err := range errors {
		fatalErrors = append(fatalErrors, err)
	}

	return processed, fatalErrors
}

// processEntityWithContent
func processEntityWithContent(file EntityFile, content []byte, contentHash uint32) ProcessedEntity {
	result := ProcessedEntity{
		File:        file,
		ContentHash: contentHash,
	}

	// YAML → JSON
	jsonData, err := yaml.YAMLToJSON(content)
	if err != nil {
		result.FatalError = fmt.Errorf("YAML→JSON: %w", err)
		return result
	}
	result.JSONData = jsonData

	// validation round 1
	var parsed map[string]any
	if err := json.Unmarshal(jsonData, &parsed); err != nil {
		result.FatalError = fmt.Errorf("invalid JSON: %w", err)
		return result
	}
	result.ParsedData = parsed

	// validation round 2
	if errs := validateStructure(parsed); len(errs) > 0 {
		result.Errors = append(result.Errors, errs...)
	}

	// validation round 3
	if errs := validateFieldsDirectly(parsed); len(errs) > 0 {
		result.Errors = append(result.Errors, errs...)
	}

	return result
}
