package service

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"strings"
	"sync"

	"github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/infrastructure/storage"
	"github.com/huydq/test/src/usecase"
)

type CSVToMySQL struct {
	S3Downloader *storage.S3Handler
	DataUC       *usecase.DataImportUsecase
}

type csvTask struct {
	Bucket      string
	Key         string
	FileType    int
	PayinFileID int
}

func NewCSVToMySQL(s3 *storage.S3Handler, dataUC *usecase.DataImportUsecase) *CSVToMySQL {
	return &CSVToMySQL{
		S3Downloader: s3,
		DataUC:       dataUC,
	}
}

func (svc *CSVToMySQL) RunJob(ctx context.Context, bucket string, workers, fileType, payinFileID int) error {
	taskCh := make(chan csvTask)
	var wg sync.WaitGroup
	sem := make(chan struct{}, workers)

	// Spawn workers
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range taskCh {
				sem <- struct{}{}
				svc.processFile(ctx, t)
				<-sem
			}
		}()
	}

	// Stream S3 Keys
	keys, err := svc.S3Downloader.StreamKeys(ctx, bucket)
	if err != nil {
		close(taskCh)
		wg.Wait()
		return err
	}
	for key := range keys {
		taskCh <- csvTask{
			Bucket:      bucket,
			Key:         key,
			FileType:    fileType,
			PayinFileID: payinFileID,
		}
	}

	close(taskCh)
	wg.Wait()
	return nil
}

func (svc *CSVToMySQL) processFile(ctx context.Context, task csvTask) {
	payinFile, err := svc.DataUC.PayinFileUsecase.FindByFilename(ctx, task.Key)
	if err != nil || payinFile == nil {
		return
	}

	if payinFile.ImportStatus == models.StatusSuccess {
		return
	}

	r, err := svc.S3Downloader.DownloadStream(ctx, task.Bucket, task.Key)
	if err != nil {
		svc.updateImportStatus(ctx, payinFile, models.StatusFailed)
		return
	}
	defer r.Close()

	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	var headers []string
	var records []map[string]string

	for scanner.Scan() {
		line := scanner.Text()
		reader := csv.NewReader(strings.NewReader(line))
		record, err := reader.Read()
		if err != nil {
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

	if len(records) > 0 {
		switch task.FileType {
		case models.PayinFileTypePaymentDetail:
			svc.DataUC.InsertPayinDetailBatch(ctx, payinFile.ID, records)
		case models.PayinFileTypePaymentReport:
			svc.DataUC.InsertPayinSummaryBatch(ctx, payinFile.ID, records)
		default:
			svc.DataUC.InsertPayinTransactionBatch(ctx, payinFile, records)
		}
	}

	svc.updateImportStatus(ctx, payinFile, models.StatusSuccess)
}

func (svc *CSVToMySQL) updateImportStatus(ctx context.Context, payinFile *models.PayinFile, status int) {
	payinFile.ImportStatus = status
	err := svc.DataUC.PayinFileUsecase.UpdateImportStatus(ctx, payinFile.ID, status)
	if err != nil {
		// Log the error or handle it as needed
		fmt.Printf("Failed to update import status: %v\n", err)
	}
}
