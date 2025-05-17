package application

import (
	"context"
	"log"
	"strconv"
	"time"

	remoteAdapter "github.com/huydq/test/batch/infrastructure/adapter/remote"
	remoteImpl "github.com/huydq/test/batch/infrastructure/adapter/remote"
	storageAdapter "github.com/huydq/test/batch/infrastructure/adapter/storage"
	storageImpl "github.com/huydq/test/batch/infrastructure/adapter/storage"
	"github.com/huydq/test/batch/infrastructure/container"
	payinPersistence "github.com/huydq/test/batch/infrastructure/persistence/payin"
	task "github.com/huydq/test/batch/task/paypay/fetch_payin_file"
	payinUsecase "github.com/huydq/test/batch/usecase/payin"
	"github.com/huydq/test/internal/pkg/database"
)

// Execute runs the upload process for PayPay CSV files to S3
func Execute(workers int, fileLoadSizePerStream int, targetDate string) {
	log.Println("======= Start UploadPaypayCSVToS3 Shell =======")
	defer log.Println("======= Stop UploadPaypayCSVToS3 Shell =======")

	// Initialize batch container and services
	batchService, err := container.NewBatchContainer()
	if err != nil {
		log.Fatalf("Failed to initialize batch container: %v", err)
	}
	defer batchService.Close()

	appConfig := batchService.AppConfig
	logger := batchService.Logger

	// Parse SSH port
	port, err := strconv.Atoi(appConfig.SSHPort)
	if err != nil {
		logger.Error("Invalid SSH port:", map[string]any{
			"error": err.Error(),
		})
		return
	}

	// Initialize SSH client
	sshConfig := remoteAdapter.SSHConfig{
		User:     appConfig.SSHUser,
		Host:     appConfig.SSHHost,
		Port:     port,
		Password: appConfig.SSHPassword,
	}
	sshClient := remoteImpl.NewSSHClient(sshConfig)

	// Initialize S3 client
	s3Config := storageAdapter.S3Config{
		Bucket:          appConfig.S3Bucket,
		Region:          appConfig.S3Region,
		AccessKeyID:     appConfig.AwsAccessKeyID,
		SecretAccessKey: appConfig.AwsSecretAccessKey,
	}
	s3Client, err := storageImpl.NewS3Client(s3Config)
	if err != nil {
		logger.Error("Failed to init S3Uploader:", map[string]any{
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

	// Initialize repositories and usecases
	fileRepo := payinPersistence.NewPayinFileRepository(batchService.DB)
	fileGroupRepo := payinPersistence.NewPayinFileGroupRepository(batchService.DB)
	fileUC := payinUsecase.NewPayinFileUsecase(fileRepo)
	fileGroupUC := payinUsecase.NewPayinFileGroupUsecase(fileGroupRepo)

	// Start timing the process
	start := time.Now()

	// Create file group
	createGroupTask := task.NewCreateFileGroupTask(fileGroupUC, appConfig.ProviderID)
	groupID, err := createGroupTask.Do(ctx)
	if err != nil {
		logger.Error("Failed to create file group:", map[string]any{
			"error": err.Error(),
		})
		return
	}

	// Initialize worker pool
	workerPoolTask := task.NewWorkerPoolTask(workers, logger)

	// Initialize stream task
	streamTask := task.NewStreamRemoteFilesTask(sshClient, targetDate)
	
	// Stream remote files
	remoteFiles, err := streamTask.Do(appConfig.RemoteDir, fileLoadSizePerStream)
	if err != nil {
		logger.Error("Failed to stream remote files:", map[string]any{
			"error": err.Error(),
		})
		return
	}

	// Process function for each file
	processFile := func(ctx context.Context, fileInfo task.RemoteFileInfo) error {
		processTask := task.NewProcessFileTask(
			fileUC,
			s3Client,
			sshClient,
			fileInfo.RemotePath,
			fileInfo.LocalPath,
			groupID,
		)
		return processTask.Do(ctx)
	}

	// Process all files using worker pool
	processedCount := workerPoolTask.ProcessRemoteFiles(ctx, remoteFiles, processFile)

	log.Printf("Job completed in %s, total processed files: %d", time.Since(start), processedCount)
}
