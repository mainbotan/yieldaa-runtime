package preset

import (
	"fmt"
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

	// Workers
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for file := range jobs {
				progress.StartJob()
				result := ProcessEntity(file)
				progress.CompleteJob()

				if result.FatalError != nil {
					errors <- fmt.Errorf("%s: %w", file.Path, result.FatalError)
				} else {
					results <- result
				}
			}
		}(i)
	}

	// Progress monitor
	stopProgress := make(chan bool)
	go func() {
		for {
			select {
			case <-stopProgress:
				return
			default:
				progress.print()
				time.Sleep(100 * time.Millisecond)
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
