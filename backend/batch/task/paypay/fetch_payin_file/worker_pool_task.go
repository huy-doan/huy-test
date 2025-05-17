package task

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/huydq/test/internal/pkg/logger"
)

// ProcessRemoteFileFunc defines the function signature for processing remote files
// It takes a context and RemoteFileInfo, processes the file, and returns error if any
type ProcessRemoteFileFunc func(ctx context.Context, fileInfo RemoteFileInfo) error

// WorkerPoolTask handles concurrent processing of remote files
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
* ProcessRemoteFiles processes remote files concurrently using a worker pool.
* It processes each file using the provided process function.
*
* @param ctx The context for the operation.
* @param files Channel of RemoteFileInfo to process.
* @param processFunc Function to process each file.
* @return int The total number of successfully processed files.
*/
func (t *WorkerPoolTask) ProcessRemoteFiles(
	ctx context.Context,
	files <-chan RemoteFileInfo,
	processFunc ProcessRemoteFileFunc,
) int {
	// Setup worker pool
	var wg sync.WaitGroup
	sem := make(chan struct{}, t.MaxWorkers)
	
	// Process each file
	for fileInfo := range files {
		wg.Add(1)
		sem <- struct{}{}
		
		// Process file in a goroutine
		go func(fi RemoteFileInfo) {
			defer func() {
				<-sem
				wg.Done()
			}()
			
			// Process the file
			err := processFunc(ctx, fi)
			if err != nil {
				t.Logger.Error("Error processing file", map[string]any{
					"remotePath": fi.RemotePath,
					"error":      err.Error(),
				})
				return
			}
			
			// Increment processed count atomically
			atomic.AddInt32(t.ProcessedCount, 1)
		}(fileInfo)
	}
	
	// Wait for all goroutines to finish
	wg.Wait()
	
	return int(*t.ProcessedCount)
}
