package ImportPaypayPayinData

import (
	"context"
	"log"
	"time"

	"github.com/huydq/test/cmd/service"
	"github.com/huydq/test/internal/infrastructure/adapter/storage"
	payinDetailRepo "github.com/huydq/test/internal/infrastructure/persistence/paypay_batch/paypay_payin_detail"
	payinFileRepo "github.com/huydq/test/internal/infrastructure/persistence/paypay_batch/paypay_payin_file"
	payinSummaryRepo "github.com/huydq/test/internal/infrastructure/persistence/paypay_batch/paypay_payin_summary"
	payinTransactionRepo "github.com/huydq/test/internal/infrastructure/persistence/paypay_batch/paypay_payin_transaction"
	"github.com/huydq/test/internal/pkg/config"
	usecase "github.com/huydq/test/internal/usecase/paypay_batch"
)

func Execute(reader, fileLoadPerStream, lineOfDataReadePerStream int) {
	log.Println("======= Start ImportPaypayPayinData Shell =======")
	defer log.Println("======= Stop ImportPaypayPayinData Shell =======")

	appConfig := config.GetConfig()

	// Init DB
	batchService, err := service.NewBatchService()
	if err != nil {
		log.Fatalf("Failed to initialize DB: %v", err)
	}
	defer batchService.Close()

	// S3 Downloader
	s3Downloader, err := storage.NewS3Service(storage.S3Config{
		Bucket:          appConfig.S3Bucket,
		Region:          appConfig.S3Region,
		AccessKeyID:     appConfig.AwsAccessKeyID,
		SecretAccessKey: appConfig.AwsSecretAccessKey,
	})
	if err != nil {
		log.Fatalf("Failed to initialize S3Downloader: %v", err)
	}

	// Initialize repositories
	payinDetailRepo := payinDetailRepo.NewPayinDetailRepository(batchService.DB)
	payinSummaryRepo := payinSummaryRepo.NewPayinSummaryRepository(batchService.DB)
	payinTransactionRepo := payinTransactionRepo.NewPayinTransactionRepository(batchService.DB)

	// Initialize PayinFileRepository
	payinFileRepo := payinFileRepo.NewPayinFileRepository(batchService.DB)

	// Initialize PayinFileUsecase with the repository
	payinFileUsecase := usecase.NewPayinFileUsecase(payinFileRepo)

	// Repositories and Usecases
	dataUC := usecase.NewDataImportUsecase(payinDetailRepo, payinSummaryRepo, payinTransactionRepo, payinFileUsecase)

	// Create ImportPaypayPayinData Service
	papaySvc := usecase.NewCSVToMySQL(s3Downloader, dataUC)

	start := time.Now()

	ctx := context.Background()
	err = papaySvc.RunBatch(ctx, appConfig.S3Bucket, reader, fileLoadPerStream, lineOfDataReadePerStream)
	if err != nil {
		log.Fatalf("RunJob error: %v", err)
	}

	log.Printf("ImportPaypayPayinData Job completed in %s", time.Since(start))
}
