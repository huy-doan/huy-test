package UploadPaypayCSVToS3

import (
	"log"
	"strconv"
	"time"

	"github.com/huydq/test/cmd/service"
	"github.com/huydq/test/internal/infrastructure/adapter/remote"
	"github.com/huydq/test/internal/infrastructure/adapter/storage"
	fileRepo "github.com/huydq/test/internal/infrastructure/persistence/paypay_batch/paypay_payin_file"
	fileGroupRepo "github.com/huydq/test/internal/infrastructure/persistence/paypay_batch/paypay_payin_filegroup"

	"github.com/huydq/test/internal/pkg/config"
	paypayUsecase "github.com/huydq/test/internal/usecase/paypay_batch"
)

func Execute(workers, fileLoadSizePerStream int) {
	log.Println("======= Start UploadPaypayCSVToS3 Shell =======")
	defer log.Println("======= Stop UploadPaypayCSVToS3 Shell =======")

	appConfig := config.GetConfig()

	batchService, err := service.NewBatchService()
	if err != nil {
		log.Fatalf("Failed to initialize DB: %v", err)
	}
	defer batchService.Close()

	port, err := strconv.Atoi(appConfig.SSHPort)
	if err != nil {
		log.Fatalf("Invalid SSH port: %v", err)
	}
	sshClient := remote.NewSSHClient(remote.SSHConfig{
		User:     appConfig.SSHUser,
		Host:     appConfig.SSHHost,
		Port:     port,
		Password: appConfig.SSHPassword,
	})

	s3Uploader, err := storage.NewS3Service(storage.S3Config{
		Bucket:          appConfig.S3Bucket,
		Region:          appConfig.S3Region,
		AccessKeyID:     appConfig.AwsAccessKeyID,
		SecretAccessKey: appConfig.AwsSecretAccessKey,
	})
	if err != nil {
		log.Fatalf("Failed to init S3Uploader: %v", err)
	}

	// Usecase wiring
	fileRepo := fileRepo.NewPayinFileRepository(batchService.DB)
	fileGroupRepo := fileGroupRepo.NewPayinFileGroupRepository(batchService.DB)
	fileUC := paypayUsecase.NewPayinFileUsecase(fileRepo)
	fileGroupUC := paypayUsecase.NewPayinFileGroupUsecase(fileGroupRepo)

	papaySvc := paypayUsecase.NewPaypayCSVToS3(batchService.DB, sshClient, s3Uploader, fileGroupUC, fileUC)

	job := paypayUsecase.Job{
		RemoteDir:  appConfig.RemoteDir,
		LocalDir:   appConfig.LocalDir,
		ProviderID: appConfig.ProviderID,
	}

	start := time.Now()

	if err := papaySvc.RunBatch(job, workers, fileLoadSizePerStream); err != nil {
		log.Fatalf("Error running job: %v", err)
	}
	log.Printf("Job completed in %s", time.Since(start))
}
