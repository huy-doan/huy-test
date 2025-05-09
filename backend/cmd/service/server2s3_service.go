package service

import (
	"context"
	"path/filepath"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/infrastructure/remote"
	"github.com/huydq/test/src/infrastructure/storage"
	"github.com/huydq/test/src/usecase"
)

type ServerToS3 struct {
	FileUC      *usecase.PayinFileUsecase
	FileGroupUC *usecase.PayinFileGroupUsecase
	S3Uploader  *storage.S3Handler
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

func NewServerToS3(db *gorm.DB, ssh *remote.SSHClient, s3 *storage.S3Handler, fileGroupUC *usecase.PayinFileGroupUsecase, fileUC *usecase.PayinFileUsecase) *ServerToS3 {
	return &ServerToS3{
		FileUC:      fileUC,
		FileGroupUC: fileGroupUC,
		S3Uploader:  s3,
		SSHClient:   ssh,
	}
}

func (svc *ServerToS3) RunJob(job Job, workers int, pageSize int) error {
	ctx := context.Background()

	group := models.PayinFileGroup{
		FileGroupName:     time.Now().Format("20060102_150405"),
		PaymentProviderID: job.ProviderID,
		ImportTargetDate:  time.Now(),
	}
	if err := svc.FileGroupUC.CreateGroup(ctx, &group); err != nil {
		return err
	}

	taskCh := make(chan workerTask)
	var wg sync.WaitGroup
	sem := make(chan struct{}, workers)

	// Spawn workers
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range taskCh {
				sem <- struct{}{}
				svc.handleTask(ctx, t)
				<-sem
			}
		}()
	}

	// Stream and dispatch tasks in a paginated manner
	stream, err := svc.SSHClient.StreamFilesPaginated(job.RemoteDir, pageSize)
	if err != nil {
		close(taskCh)
		wg.Wait()
		return err
	}

	for remotePath := range stream {
		base := filepath.Base(remotePath)
		localPath := filepath.Join(job.LocalDir, base)

		taskCh <- workerTask{
			RemotePath: remotePath,
			LocalPath:  localPath,
			GroupID:    group.ID,
		}
	}

	close(taskCh)
	wg.Wait()
	return nil
}

func (svc *ServerToS3) handleTask(ctx context.Context, t workerTask) {
	fileName := filepath.Base(t.RemotePath)

	// File-specific locking
	lock, _ := fileLocks.LoadOrStore(fileName, &sync.Mutex{})
	fileLock := lock.(*sync.Mutex)

	fileLock.Lock()
	defer fileLock.Unlock()

	// Check if the file already exists
	exists, _ := svc.FileUC.FileExistsAndDownloaded(ctx, fileName)
	if exists {
		return
	}

	// Create a new file record
	fileModel := models.PayinFile{
		PaymentProviderID: models.PaymentProviderID,
		PayinFileGroupID:  &t.GroupID,
		FileName:          fileName,
		FileContentKey:    fileName,
		DownloadStatus:    models.StatusPending,
		UploadStatus:      models.StatusPending,
	}
	if err := svc.FileUC.CreateFile(ctx, &fileModel); err != nil {
		return
	}

	// Download the file
	if err := svc.SSHClient.Download(t.RemotePath, t.LocalPath); err != nil {
		svc.FileUC.UpdateDownloadStatus(ctx, fileModel.ID, models.StatusFailed)
		return
	}
	svc.FileUC.UpdateDownloadStatus(ctx, fileModel.ID, models.StatusSuccess)

	// Upload the file
	if err := svc.S3Uploader.Upload(ctx, t.LocalPath); err != nil {
		svc.FileUC.UpdateUploadStatus(ctx, fileModel.ID, models.StatusFailed)
		return
	}
	svc.FileUC.UpdateUploadStatus(ctx, fileModel.ID, models.StatusSuccess)
}
