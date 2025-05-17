package task

import (
	"bufio"
	"context"
	"encoding/csv"
	"log"
	"strings"

	storage "github.com/huydq/test/batch/infrastructure/adapter/storage"
	payinUsecase "github.com/huydq/test/batch/usecase/payin"
	payinObject "github.com/huydq/test/internal/domain/object/payin"
	"gorm.io/gorm/logger"
)

type ProcessCSVFileTask struct {
	S3Downloader storage.S3Service
	PayinFileUC  *payinUsecase.PayinFileUsecase
	Bucket       string
	Key          string
}

func NewProcessCSVFileTask(
	s3 storage.S3Service,
	payinFileUC *payinUsecase.PayinFileUsecase,
	bucket string,
	key string,
	logger logger.Interface,
) *ProcessCSVFileTask {
	return &ProcessCSVFileTask{
		S3Downloader: s3,
		PayinFileUC:  payinFileUC,
		Bucket:       bucket,
		Key:          key,
	}
}

func (t *ProcessCSVFileTask) Do(ctx context.Context) ([]map[string]string, error) {

	payinFile, err := t.PayinFileUC.FindByFilename(ctx, t.Key)
	if err != nil {
		return nil, err
	}

	if payinFile == nil {
		return nil, nil
	}

	if payinFile.ImportStatus == payinObject.StatusSuccess {
		return nil, nil
	}

	// Download file from S3
	r, err := t.S3Downloader.DownloadStream(ctx, t.Bucket, t.Key)
	if err != nil {
		t.PayinFileUC.UpdateImportStatus(ctx, payinFile, payinObject.StatusFailed)
		return nil, err
	}
	defer r.Close()

	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024) // Increase buffer size

	var headers []string
	var records []map[string]string

	for scanner.Scan() {
		line := scanner.Text()
		reader := csv.NewReader(strings.NewReader(line))
		record, err := reader.Read()
		if err != nil {
			log.Printf("[ProcessCSVFileTask] Error reading CSV line: %v", err)
			continue
		}

		if headers == nil {
			headers = record
			continue
		}

		row := make(map[string]string)
		for i, value := range record {
			if i < len(headers) {
				row[headers[i]] = value
			}
		}
		records = append(records, row)
	}

	if scanner.Err() != nil {
		t.PayinFileUC.UpdateImportStatus(ctx, payinFile, payinObject.StatusFailed)
		return nil, scanner.Err()
	}

	return records, nil
}

func (t *ProcessCSVFileTask) MarkProcessed(ctx context.Context) error {
	payinFile, err := t.PayinFileUC.FindByFilename(ctx, t.Key)
	if err != nil || payinFile == nil {
		return err
	}

	return t.PayinFileUC.UpdateImportStatus(ctx, payinFile, payinObject.StatusSuccess)
}
