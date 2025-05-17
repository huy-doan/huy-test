package application

import (
	"context"
	"log"
	"time"

	paypayService "github.com/huydq/test/batch/domain/service/paypay"
	csvService "github.com/huydq/test/batch/domain/service/shared/csv"
	"github.com/huydq/test/batch/infrastructure/adapter/storage"
	storageImpl "github.com/huydq/test/batch/infrastructure/adapter/storage"
	"github.com/huydq/test/batch/infrastructure/container"
	payinPersistence "github.com/huydq/test/batch/infrastructure/persistence/payin"
	paypayPersistence "github.com/huydq/test/batch/infrastructure/persistence/paypay"
	task "github.com/huydq/test/batch/task/paypay/import_payin_file"
	payinUsecase "github.com/huydq/test/batch/usecase/payin"
	paypayUsecase "github.com/huydq/test/batch/usecase/paypay"
	"github.com/huydq/test/internal/pkg/database"
)

// Execute runs the import process for PayPay payin data
func Execute(readers, filesLoadPerStream, lineOfDataReadPerStream int) {
	log.Println("======= Start ImportPaypayPayinData Shell =======")
	defer log.Println("======= Stop ImportPaypayPayinData Shell =======")

	// Initialize batch container and services
	batchService, err := container.NewBatchContainer()
	if err != nil {
		log.Fatalf("Failed to initialize DB: %v", err)
	}
	defer batchService.Close()

	appConfig := batchService.AppConfig
	logger := batchService.Logger

	s3Config := storage.S3Config{
		Bucket:          appConfig.S3Bucket,
		Region:          appConfig.S3Region,
		AccessKeyID:     appConfig.AwsAccessKeyID,
		SecretAccessKey: appConfig.AwsSecretAccessKey,
	}
	s3Client, err := storageImpl.NewS3Client(s3Config)
	if err != nil {
		logger.Error("Failed to initialize S3Downloader", map[string]any{
			"error": err.Error(),
		})
		return
	}

	// Setup context with DB
	ctx := context.Background()
	ctx, dbSetErr := database.SetDB(ctx, batchService.DB)
	if dbSetErr != nil {
		logger.Error("Failed to set DB in context:", map[string]any{
			"error": dbSetErr.Error(),
		})
		return
	}

	// Initialize repositories
	payinFileRepo := payinPersistence.NewPayinFileRepository(batchService.DB)
	paypayPayinDetailRepo := paypayPersistence.NewPayinDetailRepository(batchService.DB)
	paypayPayinSummaryRepo := paypayPersistence.NewPayinSummaryRepository(batchService.DB)
	paypayPayinTransactionRepo := paypayPersistence.NewPayinTransactionRepository(batchService.DB)

	// Initialize usecases
	payinFileUC := payinUsecase.NewPayinFileUsecase(payinFileRepo)
	detailUC := paypayUsecase.NewPayinDetailUsecase(paypayPayinDetailRepo, logger)
	summaryUC := paypayUsecase.NewPayinSummaryUsecase(paypayPayinSummaryRepo, logger)
	transactionUC := paypayUsecase.NewPayinTransactionUsecase(paypayPayinTransactionRepo, logger)

	// Initialize domain services
	csvReaderService := csvService.NewCsvReaderService()
	validateFieldsService := paypayService.NewValidateCSVFieldsService()
	multiSectionImportService := paypayService.NewMultiSectionCSVImportService(
		summaryUC.ProcessAndInsertSummaries,
		detailUC.ProcessAndInsertDetails,
	)

	// Initialize tasks
	filterTask := task.NewFilterS3KeysTask(nil)
	targetFolders := filterTask.BuildTargetFolders(
		appConfig.RemoteDir,
		appConfig.TopUpSummaryDetailsPath,
		appConfig.TopUpReportPath,
	)
	filterTask.TargetFolders = targetFolders

	zipProcessor := task.NewProcessZipFileTask(
		s3Client,
		payinFileUC,
		csvReaderService,
		validateFieldsService,
		multiSectionImportService,
		transactionUC,
		appConfig.S3Bucket,
		appConfig.RemoteDir,
		appConfig.TopUpReportPath,
		appConfig.TopUpSummaryDetailsPath,
		logger,
	)
	
	// Initialize worker pool task
	workerPoolTask := task.NewWorkerPoolTask(readers, logger)

	// Start the import process
	start := time.Now()

	// Stream keys from S3
	keys, err := s3Client.StreamKeys(ctx, appConfig.S3Bucket)
	if err != nil {
		logger.Error("Failed to stream keys from S3:", map[string]any{
			"error": err.Error(),
		})
		return
	}

	// Process files using worker pool with explicitly typed functions
	processedCount := workerPoolTask.ProcessS3Keys(
		ctx,
		keys,
		filterTask.Do,                // S3KeyFilterFunc - filters keys based on folder and extension
		zipProcessor.Do,              // ZipFileProcessFunc - processes zip files containing CSV data
	)
	
	log.Printf("ImportPaypayPayinData job completed in %s, processed %d files", time.Since(start), processedCount)
}
