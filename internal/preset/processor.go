package preset

import (
	"fmt"
	"os"
	"sync"
	"time"
)

func ProcessEntities(files []EntityFile, maxWorkers int) ([]ProcessedEntity, []error) {
	if len(files) == 0 {
		return []ProcessedEntity{}, nil
	}

	if maxWorkers <= 0 {
		maxWorkers = DefaultWorkers
	}
	if maxWorkers > len(files) {
		maxWorkers = len(files)
	}

	progress := NewProgressTracker(len(files))
	jobs := make(chan EntityFile, len(files))
	results := make(chan ProcessedEntity, len(files))
	errors := make(chan error, len(files))

	var wg sync.WaitGroup
	var seenHashes sync.Map   // thread-safe для хешей контента
	var seenKeys sync.Map     // thread-safe для проверки ключей сущностей
	var keyConflicts []string // для сбора конфликтов

	var conflictsMu sync.Mutex // мьютекс для keyConflicts

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

				// xxHash64 вместо CRC32
				contentHash := calculateContentHash(content)

				// Atomic check and store
				if _, alreadyProcessed := seenHashes.LoadOrStore(contentHash, true); alreadyProcessed {
					progress.CompleteJob()
					continue
				}

				result := ProcessEntity(file, content)

				// check entity key
				if result.ParsedData != nil && result.FatalError == nil {
					key := EntityKey(result.ParsedData)
					if key != "" {
						if existingFile, exists := seenKeys.LoadOrStore(key, file.Path); exists {
							result.Errors = append(result.Errors,
								fmt.Sprintf("entity key conflict: '%s' already defined in '%s'",
									key, existingFile.(string)))

							conflictsMu.Lock()
							keyConflicts = append(keyConflicts,
								fmt.Sprintf("  %s:\n    • %s\n    • %s",
									key, existingFile.(string), file.Path))
							conflictsMu.Unlock()
						}
					}
				}

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
