package paypay_batch

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"strings"
	"sync"

	"github.com/huydq/test/internal/infrastructure/adapter/storage"
	payinFileModel "github.com/huydq/test/internal/infrastructure/persistence/paypay_batch/paypay_payin_file/dto"
)

type CSVToMySQL struct {
	S3Downloader *storage.S3ServiceConfig
	DataUC       *DataImportUsecase
}

type csvTask struct {
	Bucket      string
	Key         string
	FileType    int
	PayinFileID int
}

func NewCSVToMySQL(s3 *storage.S3ServiceConfig, dataUC *DataImportUsecase) *CSVToMySQL {
	return &CSVToMySQL{
		S3Downloader: s3,
		DataUC:       dataUC,
	}
}

func (svc *CSVToMySQL) RunBatch(ctx context.Context, bucket string, readers, filesLoadPerStream, lineOfDataPerStream int) error {
	taskCh := make(chan csvTask)
	var wg sync.WaitGroup
	sem := make(chan struct{}, readers)

	// Spawn readers
	for range readers {
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
		payinFile, err := svc.DataUC.PayinFileUsecase.FindByFilename(ctx, key)
		if err != nil || payinFile == nil {
			// Log the error or handle it as needed
			fmt.Printf("Failed to find PayinFile for key %s: %v\n", key, err)
			continue
		}

		taskCh <- csvTask{
			Bucket:      bucket,
			Key:         key,
			FileType:    payinFile.PayinFileType, // Retrieve FileType from PayinFile
			PayinFileID: payinFile.ID,            // Retrieve ID from PayinFile
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

	if payinFile.ImportStatus == payinFileModel.StatusSuccess {
		return
	}

	r, err := svc.S3Downloader.DownloadStream(ctx, task.Bucket, task.Key)
	if err != nil {
		svc.DataUC.PayinFileUsecase.UpdateImportStatus(ctx, payinFile.ID, payinFileModel.StatusFailed)
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
		case payinFileModel.PayinFileTypePaymentDetail:
			svc.DataUC.InsertPayinDetailBatch(ctx, payinFile.ID, records)
		case payinFileModel.PayinFileTypePaymentSummary:
			svc.DataUC.InsertPayinSummaryBatch(ctx, payinFile.ID, records)
		case payinFileModel.PayinFileTypePaymentTransaction:
			svc.DataUC.InsertPayinTransactionBatch(ctx, payinFile.ID, records)
		default:
			// Handle unknown file type if necessary
		}
	}
	svc.DataUC.PayinFileUsecase.UpdateImportStatus(ctx, payinFile.ID, payinFileModel.StatusSuccess)
}
