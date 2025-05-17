package task

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"log"
	"strings"

	paypayService "github.com/huydq/test/batch/domain/service/paypay"
	csvService "github.com/huydq/test/batch/domain/service/shared/csv"
	storageService "github.com/huydq/test/batch/infrastructure/adapter/storage"
	payinUsecase "github.com/huydq/test/batch/usecase/payin"
	paypayUsecase "github.com/huydq/test/batch/usecase/paypay"
	model "github.com/huydq/test/internal/domain/model/payin"
	object "github.com/huydq/test/internal/domain/object/payin"
	"github.com/huydq/test/internal/pkg/logger"
)

// ProcessZipFileTask handles the processing of ZIP files containing CSV data
type ProcessZipFileTask struct {
	S3Client                 storageService.S3Service
	PayinFileUC              *payinUsecase.PayinFileUsecase
	CSVReaderService         *csvService.CsvReaderService
	ValidateFieldsService    *paypayService.ValidateCSVFieldsService
	MultiSectionImportService *paypayService.MultiSectionCSVImportService
	TransactionUC            *paypayUsecase.PayinTransactionUsecase
	S3Bucket                 string
	RemoteDir                string
	TopUpReportPath          string
	TopUpSummaryDetailsPath  string
	Logger                   logger.Logger
}

// NewProcessZipFileTask creates a new instance of ProcessZipFileTask
func NewProcessZipFileTask(
	s3Client storageService.S3Service,
	payinFileUC *payinUsecase.PayinFileUsecase,
	csvReaderService *csvService.CsvReaderService,
	validateFieldsService *paypayService.ValidateCSVFieldsService, 
	multiSectionImportService *paypayService.MultiSectionCSVImportService,
	transactionUC *paypayUsecase.PayinTransactionUsecase,
	s3Bucket string,
	remoteDir string,
	topUpReportPath string,
	topUpSummaryDetailsPath string,
	logger logger.Logger,
) *ProcessZipFileTask {
	return &ProcessZipFileTask{
		S3Client:                 s3Client,
		PayinFileUC:              payinFileUC,
		CSVReaderService:         csvReaderService,
		ValidateFieldsService:    validateFieldsService,
		MultiSectionImportService: multiSectionImportService,
		TransactionUC:            transactionUC,
		S3Bucket:                 s3Bucket,
		RemoteDir:                remoteDir,
		TopUpReportPath:          topUpReportPath,
		TopUpSummaryDetailsPath:  topUpSummaryDetailsPath,
		Logger:                   logger,
	}
}

/**
* Do processes a ZIP file from S3 containing CSV data.
* It extracts the CSV file, validates its headers, and processes the data.
* Returns true if the file was successfully processed, false otherwise.
*
* @param ctx The context for the operation.
* @param s3Key The S3 key of the ZIP file to process.
*
* @return bool True if the file was successfully processed, false otherwise.
* @return error Any error that occurred during processing.
*/
func (t *ProcessZipFileTask) Do(ctx context.Context, s3Key string) (bool, error) {
	t.Logger.Info("[Import] Starting import for:", map[string]any{
		"info": s3Key,
	})

	// Extract filename from key
	fileName := s3Key
	if idx := strings.LastIndex(s3Key, "/"); idx >= 0 {
		fileName = s3Key[idx+1:]
	}

	// Find payin file record in database
	payinFile, err := t.PayinFileUC.FindByFilename(ctx, fileName)
	if err != nil || payinFile == nil {
		t.Logger.Info("[Import] No PayinFile found for:", map[string]any{
			"info": fileName,
		})
		return false, err
	}

	// Download ZIP file from S3
	zipStream, err := t.S3Client.DownloadStream(ctx, t.S3Bucket, s3Key)
	if err != nil {
		t.PayinFileUC.UpdateImportStatus(ctx, payinFile, object.StatusFailed)
		t.Logger.Error("Failed to download zip", map[string]any{
			"error": err.Error(),
			"key":   s3Key,
		})
		return false, err
	}
	defer zipStream.Close()

	// Read ZIP file content
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, zipStream); err != nil {
		t.PayinFileUC.UpdateImportStatus(ctx, payinFile, object.StatusFailed)
		t.Logger.Error("Failed to read zip", map[string]any{
			"error": err.Error(),
			"key":   s3Key,
		})
		return false, err
	}

	// Create a ZIP reader
	zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil || len(zr.File) == 0 {
		log.Printf("[Import] Invalid zip or empty: %s", s3Key)
		t.PayinFileUC.UpdateImportStatus(ctx, payinFile, object.StatusFailed)
		t.Logger.Error("Invalid zip or empty", map[string]any{
			"error": err.Error(),
			"key":   s3Key,
		})
		return false, err
	}

	// Find CSV file in the ZIP archive
	var csvReader io.ReadCloser
	var csvKey string
	for _, f := range zr.File {
		if strings.HasSuffix(strings.ToLower(f.Name), ".csv") {
			rc, err := f.Open()
			if err != nil {
				continue
			}
			csvReader = rc
			csvKey = f.Name
			break
		}
	}
	if csvReader == nil {
		log.Printf("[Import] No csv in zip: %s", s3Key)
		t.PayinFileUC.UpdateImportStatus(ctx, payinFile, object.StatusFailed)
		t.Logger.Error("No csv in zip", map[string]any{
			"key": s3Key,
		})
		return false, nil
	}
	defer csvReader.Close()

	// Process based on file path and type
	joinRemotePath := func(base, sub string) string {
		sub = strings.TrimPrefix(sub, "/")
		return strings.TrimRight(base, "/") + "/" + sub
	}

	// Process Top-Up Report (multi-section import)
	if strings.HasPrefix(s3Key, strings.TrimLeft(joinRemotePath(t.RemoteDir, t.TopUpReportPath), "/")) {
		return t.processMultiSectionFile(ctx, s3Key, csvReader, payinFile)
	} 
	
	// Process Summary Details (transaction import)
	if strings.HasPrefix(s3Key, strings.TrimLeft(joinRemotePath(t.RemoteDir, t.TopUpSummaryDetailsPath), "/")) {
		return t.processTransactionFile(ctx, s3Key, csvReader, csvKey, payinFile, zr)
	}

	return false, nil
}

// processMultiSectionFile processes a multi-section CSV file
func (t *ProcessZipFileTask) processMultiSectionFile(ctx context.Context, key string, csvReader io.ReadCloser, payinFile *model.PayinFile) (bool, error) {
	log.Printf("[Import] Reading CSV for multi-section import: %s", key)
	
	var lines []string
	buf := new(bytes.Buffer)
	n, err := io.Copy(buf, csvReader)
	log.Printf("[Import] io.Copy read %d bytes for %s (err: %v)", n, key, err)
	if err != nil {
		t.PayinFileUC.UpdateImportStatus(ctx, payinFile, object.StatusFailed)
		t.Logger.Error("Failed to read csv in zip", map[string]any{
			"error": err.Error(),
			"key":   key,
		})
		return false, err
	}
	
	log.Printf("[Import] Buffer length after io.Copy: %d for %s", buf.Len(), key)
	for _, line := range strings.Split(buf.String(), "\n") {
		lines = append(lines, strings.TrimRight(line, "\r"))
	}
	
	if len(lines) > 0 {
		log.Printf("[Import] First 2 lines: %q", lines[:min(2, len(lines))])
	}
	
	log.Printf("[Import] Read %d lines from CSV for %s", len(lines), key)
	if len(lines) == 0 {
		log.Printf("[Import] No lines found in CSV for %s", key)
		return false, nil
	}
	
	log.Printf("[Import] Calling multiSectionImportService.ProcessLines for %s", key)
	log.Printf("[Import] Attempting multi-section import for %s", key)
	if err := t.MultiSectionImportService.ProcessLines(ctx, payinFile.ID, lines); err != nil {
		t.PayinFileUC.UpdateImportStatus(ctx, payinFile, object.StatusFailed)
		t.Logger.Error("Multi-section import error", map[string]any{
			"error": err.Error(),
			"key":   key,
		})
		return false, err
	}
	
	t.PayinFileUC.UpdateImportStatus(ctx, payinFile, object.StatusSuccess)
	log.Printf("[Import] Successfully processed file %s", key)
	log.Printf("[Import] Finished import attempt for: %s", key)
	return true, nil
}

// processTransactionFile processes a transaction CSV file
func (t *ProcessZipFileTask) processTransactionFile(ctx context.Context, key string, csvReader io.ReadCloser, csvKey string, payinFile *model.PayinFile, zipReader *zip.Reader) (bool, error) {
	log.Printf("[Import] Reading CSV for transaction import: %s", key)
	
	// Debug: peek at first 256 bytes of CSV
	peekBuf := make([]byte, 256)
	n, _ := csvReader.Read(peekBuf)
	log.Printf("[Import] First %d bytes of CSV for %s: %q", n, key, string(peekBuf[:n]))
	csvReader.Close() // close the current reader
	
	// Reopen the CSV file
	for _, f := range zipReader.File {
		if f.Name == csvKey {
			rc, err := f.Open()
			if err == nil {
				csvReader = rc
				defer csvReader.Close()
				break
			}
		}
	}

	// Read the CSV file using domain service
	records, err := t.CSVReaderService.ReadWithHeader(csvReader)
	if err != nil || len(records) == 0 {
		t.PayinFileUC.UpdateImportStatus(ctx, payinFile, object.StatusFailed)
		t.Logger.Error("CSV invalid or empty", map[string]any{
			"error": err.Error(),
			"key":   key,
		})
		return false, err
	}

	// Extract and validate headers
	var actualHeaders []string
	var rawHeaders []string
	if len(records) > 0 {
		for k := range records[0] {
			actualHeaders = append(actualHeaders, strings.TrimSpace(k))
		}
		// Try to get the raw headers from the first record if available
		for k := range records[0] {
			rawHeaders = append(rawHeaders, k)
		}
	}

	isValid, requiredHeaders := t.ValidateFieldsService.ValidateHeaders(actualHeaders, object.PayinFileTypePaymentTransaction)
	if !isValid {
		t.PayinFileUC.UpdateImportStatus(ctx, payinFile, object.StatusFailed)
		t.Logger.Error("CSV missing column", map[string]any{
			"key":             csvKey,
			"requiredHeaders": requiredHeaders,
			"actualHeaders":   actualHeaders,
			"rawHeaders":      rawHeaders,
		})
		return false, nil
	}
	
	log.Printf("[Import] Attempting transaction import for %s", key)
	insertErr := t.TransactionUC.ProcessAndInsertTransactions(ctx, payinFile.ID, records)
	if insertErr != nil {
		t.PayinFileUC.UpdateImportStatus(ctx, payinFile, object.StatusFailed)
		t.Logger.Error("Insert error", map[string]any{
			"error": insertErr.Error(),
			"key":   key,
		})
		return false, insertErr
	}
	
	t.PayinFileUC.UpdateImportStatus(ctx, payinFile, object.StatusSuccess)
	return true, nil
}
