package paypay_batch

import (
	"context"
	"log"
	"path/filepath"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/huydq/test/internal/infrastructure/adapter/remote"
	"github.com/huydq/test/internal/infrastructure/adapter/storage"
	payinFileModel "github.com/huydq/test/internal/infrastructure/persistence/paypay_batch/paypay_payin_file/dto"
	payinFileGroupModel "github.com/huydq/test/internal/infrastructure/persistence/paypay_batch/paypay_payin_filegroup/dto"
)

type PaypayCSVToS3 struct {
	FileUC      *PayinFileUsecase
	FileGroupUC *PayinFileGroupUsecase
	S3Uploader  *storage.S3ServiceConfig
	SSHClient   *remote.SSHClient
}

type Job struct {
	RemoteDir  string
	LocalDir   string
	ProviderID int
}

type workerTask struct {
	RemotePath string
	LocalPath  string
	GroupID    int
}

var fileLocks sync.Map // File-specific locks

func NewPaypayCSVToS3(db *gorm.DB, ssh *remote.SSHClient, s3 *storage.S3ServiceConfig, fileGroupUC *PayinFileGroupUsecase, fileUC *PayinFileUsecase) *PaypayCSVToS3 {
	return &PaypayCSVToS3{
		FileUC:      fileUC,
		FileGroupUC: fileGroupUC,
		S3Uploader:  s3,
		SSHClient:   ssh,
	}
}

func (svc *PaypayCSVToS3) RunBatch(job Job, workers int, pageSize int) error {
	ctx := context.Background()
	log.Println("[RunBatch] Starting job")

	group := payinFileGroupModel.PayinFileGroup{
		FileGroupName:     time.Now().Format("20060102_150405"),
		PaymentProviderID: job.ProviderID,
		ImportTargetDate:  time.Now(),
	}
	if err := svc.FileGroupUC.CreateGroup(ctx, &group); err != nil {
		log.Printf("[RunBatch] Failed to create group: %v", err)
		return err
	}
	log.Printf("[RunBatch] Created file group: %s (ID: %d)", group.FileGroupName, group.ID)

	taskCh := make(chan workerTask)
	var wg sync.WaitGroup
	sem := make(chan struct{}, workers)

	log.Printf("[RunBatch] Spawning %d workers", workers)
	for i := range workers {
		wg.Add(1)
		go func(workerID int) {
			log.Printf("[Worker-%d] Started", workerID)
			defer wg.Done()
			for t := range taskCh {
				log.Printf("[Worker-%d] Received task: %s", workerID, t.RemotePath)
				sem <- struct{}{}
				svc.processBatchTask(ctx, t)
				<-sem
			}
			log.Printf("[Worker-%d] Exiting", workerID)
		}(i)
	}

	stream, err := svc.SSHClient.StreamFilesPaginated(job.RemoteDir, pageSize)
	if err != nil {
		log.Printf("[RunBatch] Failed to stream files: %v", err)
		close(taskCh)
		wg.Wait()
		return err
	}

	log.Println("[RunBatch] Streaming files...")
	count := 0
	for remotePath := range stream {
		count++
		base := filepath.Base(remotePath)
		localPath := filepath.Join(job.LocalDir, base)

		log.Printf("[RunBatch] Queuing file for download: %s -> %s", remotePath, localPath)
		taskCh <- workerTask{
			RemotePath: remotePath,
			LocalPath:  localPath,
			GroupID:    group.ID,
		}
	}
	log.Printf("[RunBatch] Dispatched %d tasks", count)

	close(taskCh)
	wg.Wait()
	log.Println("[RunBatch] Job complete")
	return nil
}

func (svc *PaypayCSVToS3) processBatchTask(ctx context.Context, t workerTask) {
	fileName := filepath.Base(t.RemotePath)
	log.Printf("[handleTask] Starting task for file: %s", fileName)

	// File-specific locking
	lock, _ := fileLocks.LoadOrStore(fileName, &sync.Mutex{})
	fileLock := lock.(*sync.Mutex)

	fileLock.Lock()
	defer fileLock.Unlock()

	exists, err := svc.FileUC.FileExistsAndDownloaded(ctx, fileName)
	if err != nil {
		log.Printf("[handleTask] FileExistsAndDownloaded error: %v", err)
	}
	if exists {
		log.Printf("[handleTask] File already processed: %s", fileName)
		return
	}

	fileModel := payinFileModel.PayinFile{
		PaymentProviderID: payinFileModel.PaymentProviderID,
		PayinFileGroupID:  &t.GroupID,
		FileName:          fileName,
		FileContentKey:    fileName,
		DownloadStatus:    payinFileModel.StatusPending,
		UploadStatus:      payinFileModel.StatusPending,
		ImportStatus:      payinFileModel.StatusPending,
	}
	created, err := svc.FileUC.CreateFile(ctx, &fileModel)
	if err != nil {
		log.Printf("[handleTask] Failed to create file record: %v", err)
		return
	}
	log.Printf("[handleTask] Created file record: %s (ID: %d)", fileName, created.ID)

	// Download file
	log.Printf("[handleTask] Downloading: %s -> %s", t.RemotePath, t.LocalPath)
	if err := svc.SSHClient.Download(t.RemotePath, t.LocalPath); err != nil {
		log.Printf("[ERROR] Download failed for %s: %v", t.RemotePath, err)
		svc.FileUC.UpdateDownloadStatus(ctx, fileModel.ID, payinFileModel.StatusFailed)
		return
	}
	svc.FileUC.UpdateDownloadStatus(ctx, fileModel.ID, payinFileModel.StatusSuccess)

	// Upload to S3
	log.Printf("[handleTask] Uploading to S3: %s", t.LocalPath)
	if err := svc.S3Uploader.Upload(ctx, t.LocalPath); err != nil {
		log.Printf("[ERROR] Upload failed for %s: %v", fileName, err)
		svc.FileUC.UpdateUploadStatus(ctx, fileModel.ID, payinFileModel.StatusFailed)
		return
	}
	log.Printf("[handleTask] Uploaded to S3: %s", fileName)
	svc.FileUC.UpdateUploadStatus(ctx, fileModel.ID, payinFileModel.StatusSuccess)
}
