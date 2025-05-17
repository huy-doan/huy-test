package task

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/huydq/test/internal/pkg/logger"
)

// S3KeyFilterFunc defines the function signature for filtering S3 keys
// It takes an S3 key string and returns true if the key should be processed,
// false otherwise
type S3KeyFilterFunc func(s3Key string) bool

// ZipFileProcessFunc defines the function signature for processing ZIP files
// It takes a context and an S3 key string, processes the file, and returns:
// - bool: true if the file was successfully processed
// - error: any error that occurred during processing
type ZipFileProcessFunc func(ctx context.Context, s3Key string) (bool, error)

// WorkerPoolTask handles the concurrent processing of files using a worker pool
type WorkerPoolTask struct {
	MaxWorkers     int           // Maximum number of concurrent workers
	Logger         logger.Logger // Logger for recording events
	ProcessedCount *int32        // Count of successfully processed files (atomic)
}

// NewWorkerPoolTask creates a new instance of WorkerPoolTask
func NewWorkerPoolTask(maxWorkers int, logger logger.Logger) *WorkerPoolTask {
	var count int32 = 0
	return &WorkerPoolTask{
		MaxWorkers:     maxWorkers,
		Logger:         logger,
		ProcessedCount: &count,
	}
}

/**
 * ProcessS3Keys processes S3 keys concurrently using a worker pool.
 * It filters the keys using the provided filter function and processes
 * each selected key using the provided process function.
 *
 * @param ctx The context for the operation.
 * @param keys Channel of S3 keys to process.
 * @param filterFunc Function to determine which keys should be processed.
 * @param processFunc Function to process each selected key.
 * @return int The total number of successfully processed files.
 */
func (t *WorkerPoolTask) ProcessS3Keys(
	ctx context.Context,
	keys <-chan string,
	filterFunc S3KeyFilterFunc,
	processFunc ZipFileProcessFunc,
) int {
	// Setup worker pool
	var wg sync.WaitGroup
	sem := make(chan struct{}, t.MaxWorkers)
	
	// Process each key
	for key := range keys {
		// Filter keys
		if !filterFunc(key) {
			continue
		}
		
		wg.Add(1)
		sem <- struct{}{}
		
		// Process file in a goroutine
		go func(s3Key string) {
			defer func() {
				<-sem
				wg.Done()
			}()
			
			// Process the file
			processed, err := processFunc(ctx, s3Key)
			if err != nil {
				t.Logger.Error("Error processing file", map[string]any{
					"key":   s3Key,
					"error": err.Error(),
				})
				return
			}
			
			if processed {
				// Increment processed count atomically
				atomic.AddInt32(t.ProcessedCount, 1)
			}
		}(key)
	}
	
	// Wait for all goroutines to finish
	wg.Wait()
	
	return int(*t.ProcessedCount)
}
