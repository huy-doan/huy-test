package task

import (
	"context"
	"fmt"
	"log"

	storage "github.com/huydq/test/batch/infrastructure/adapter/storage"
	payinUsecase "github.com/huydq/test/batch/usecase/payin"
	object "github.com/huydq/test/internal/domain/object/payin"
)

type PayinFileInfo struct {
	Bucket      string
	Key         string
	FileType    object.PayinFileType
	PayinFileID int
}

type StreamPayinFilesTask struct {
	S3Downloader storage.S3Service
	PayinFileUC  *payinUsecase.PayinFileUsecase
	Bucket       string
}

func NewStreamPayinFilesTask(
	s3 storage.S3Service,
	payinFileUC *payinUsecase.PayinFileUsecase,
	bucket string,
) *StreamPayinFilesTask {
	return &StreamPayinFilesTask{
		S3Downloader: s3,
		PayinFileUC:  payinFileUC,
		Bucket:       bucket,
	}
}

func (t *StreamPayinFilesTask) Do(ctx context.Context) (<-chan PayinFileInfo, error) {
	fileInfoCh := make(chan PayinFileInfo)

	keys, err := t.S3Downloader.StreamKeys(ctx, t.Bucket)
	if err != nil {
		close(fileInfoCh)
		return nil, fmt.Errorf("failed to stream keys from S3: %w", err)
	}

	go func() {
		defer close(fileInfoCh)

		for key := range keys {
			payinFile, err := t.PayinFileUC.FindByFilename(ctx, key)
			if err != nil || payinFile == nil {
				log.Printf("Failed to find PayinFile for key %s: %v", key, err)
				continue
			}

			fileInfoCh <- PayinFileInfo{
				Bucket:      t.Bucket,
				Key:         key,
				FileType:    payinFile.PayinFileType, // Retrieve FileType from PayinFile
				PayinFileID: payinFile.ID,            // Retrieve ID from PayinFile
			}
		}
	}()

	return fileInfoCh, nil
}
