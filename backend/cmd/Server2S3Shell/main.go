package Server2S3Shell

import (
	"log"
	"strconv"
	"time"

	"github.com/vnlab/makeshop-payment/cmd/service"
	"github.com/vnlab/makeshop-payment/src/infrastructure/config"
	dbRepo "github.com/vnlab/makeshop-payment/src/infrastructure/persistence/repositories"
	"github.com/vnlab/makeshop-payment/src/infrastructure/remote"
	"github.com/vnlab/makeshop-payment/src/infrastructure/storage"
	"github.com/vnlab/makeshop-payment/src/usecase"
)

func Execute(workers, fileLoadSizePerStream int) {
	log.Println("======= Start ServerToS3 Shell =======")
	defer log.Println("======= Stop ServerToS3 Shell =======")

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

	s3Uploader, err := storage.NewS3Handler(storage.S3Config{
		Bucket:          appConfig.S3Bucket,
		Region:          appConfig.S3Region,
		AccessKeyID:     appConfig.AwsAccessKeyID,
		SecretAccessKey: appConfig.AwsSecretAccessKey,
	})
	if err != nil {
		log.Fatalf("Failed to init S3Uploader: %v", err)
	}

	// Usecase wiring
	fileRepo := dbRepo.NewPayinFileRepository(batchService.DB)
	fileGroupRepo := dbRepo.NewPayinFileGroupRepository(batchService.DB)
	fileUC := usecase.NewPayinFileUsecase(fileRepo)
	fileGroupUC := usecase.NewPayinFileGroupUsecase(fileGroupRepo)

	svc := service.NewServerToS3(batchService.DB, sshClient, s3Uploader, fileGroupUC, fileUC)

	job := service.Job{
		RemoteDir:  appConfig.RemoteDir,
		LocalDir:   appConfig.LocalDir,
		ProviderID: appConfig.ProviderID,
	}

	start := time.Now()
	if err := svc.RunJob(job, workers, fileLoadSizePerStream); err != nil {
		log.Fatalf("Error running job: %v", err)
	}
	log.Printf("Job completed in %s", time.Since(start))
}
