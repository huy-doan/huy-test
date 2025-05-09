package CSV2MysqlShell

import (
	"context"
	"log"
	"time"

	"github.com/huydq/test/cmd/service"
	"github.com/huydq/test/src/infrastructure/config"
	dbRepo "github.com/huydq/test/src/infrastructure/persistence/repositories"
	"github.com/huydq/test/src/infrastructure/storage"
	"github.com/huydq/test/src/usecase"
)

func Execute(reader, fileLoadPerStream, lineOfDataReadePerStream int) {
	log.Println("======= Start CSV2Mysql Shell =======")
	defer log.Println("======= Stop CSV2Mysql Shell =======")

	appConfig := config.GetConfig()

	// Init DB
	batchService, err := service.NewBatchService()
	if err != nil {
		log.Fatalf("Failed to initialize DB: %v", err)
	}
	defer batchService.Close()

	// S3 Downloader
	s3Downloader, err := storage.NewS3Handler(storage.S3Config{
		Bucket:          appConfig.S3Bucket,
		Region:          appConfig.S3Region,
		AccessKeyID:     appConfig.AwsAccessKeyID,
		SecretAccessKey: appConfig.AwsSecretAccessKey,
	})
	if err != nil {
		log.Fatalf("Failed to initialize S3Downloader: %v", err)
	}

	// Initialize repositories
	payinDetailRepo := dbRepo.NewPayinDetailRepository(batchService.DB)
	payinSummaryRepo := dbRepo.NewPayinSummaryRepository(batchService.DB)
	payinTransactionRepo := dbRepo.NewPayinTransactionRepository(batchService.DB)

	// Initialize PayinFileRepository
	payinFileRepo := dbRepo.NewPayinFileRepository(batchService.DB)

	// Initialize PayinFileUsecase with the repository
	payinFileUsecase := usecase.NewPayinFileUsecase(payinFileRepo)

	// Repositories and Usecases
	dataUC := usecase.NewDataImportUsecase(payinDetailRepo, payinSummaryRepo, payinTransactionRepo, payinFileUsecase)

	// Create Csv2Mysql Service
	csv2Mysql := service.NewCSVToMySQL(s3Downloader, dataUC)

	start := time.Now()

	ctx := context.Background()
	err = csv2Mysql.RunJob(ctx, appConfig.S3Bucket, reader, fileLoadPerStream, lineOfDataReadePerStream) // Updated arguments
	if err != nil {
		log.Fatalf("RunJob error: %v", err)
	}

	log.Printf("CSV2Mysql Job completed in %s", time.Since(start))
}
