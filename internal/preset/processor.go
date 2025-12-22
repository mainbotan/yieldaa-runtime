package preset

import (
	"fmt"
	"hash/crc32"
	"os"
	"sync"
	"time"
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

				content, err := os.ReadFile(file.Path)
				if err != nil {
					errors <- fmt.Errorf("%s: read: %w", file.Path, err)
					progress.CompleteJob()
					continue
				}

				contentHash := crc32.ChecksumIEEE(content)

				// Atomic check and store
				if _, alreadyProcessed := seenHashes.LoadOrStore(contentHash, true); alreadyProcessed {
					progress.CompleteJob()
					continue
				}

				result := ProcessEntity(file)

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
