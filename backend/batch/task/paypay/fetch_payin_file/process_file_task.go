package task

import (
	"bytes"
	"context"
	"log"
	"path/filepath"
	"sync"

	remoteAdapter "github.com/huydq/test/batch/infrastructure/adapter/remote"
	storageAdapter "github.com/huydq/test/batch/infrastructure/adapter/storage"
	payinUsecase "github.com/huydq/test/batch/usecase/payin"
	payinModel "github.com/huydq/test/internal/domain/model/payin"
	payinObject "github.com/huydq/test/internal/domain/object/payin"
	"github.com/huydq/test/internal/pkg/database"
	"gorm.io/gorm"
)

// ProcessFileTask represents a task for processing a file, including downloading, uploading to S3, and updating statuses.
type ProcessFileTask struct {
	FileUC     *payinUsecase.PayinFileUsecase // Use case for handling file-related operations
	S3Uploader storageAdapter.S3Service       // S3 uploader configuration
	SSHClient  remoteAdapter.SSHService       // SSH client for remote file access
	RemotePath string                         // Path to the remote file
	LocalPath  string                         // Local path for processing
	GroupID    int                            // Group ID associated with the file
}

var fileLocks sync.Map // Map to store file locks for concurrency control

// NewProcessFileTask initializes a new ProcessFileTask instance.
func NewProcessFileTask(
	fileUC *payinUsecase.PayinFileUsecase,
	s3 storageAdapter.S3Service,
	ssh remoteAdapter.SSHService,
	remotePath string,
	localPath string,
	groupID int,
) *ProcessFileTask {
	return &ProcessFileTask{
		FileUC:     fileUC,
		S3Uploader: s3,
		SSHClient:  ssh,
		RemotePath: remotePath,
		LocalPath:  localPath,
		GroupID:    groupID,
	}
}

// Do executes the file processing task, including downloading, uploading, and updating statuses.
func (t *ProcessFileTask) Do(ctx context.Context) error {
	fileName := filepath.Base(t.RemotePath) // Extract the file name from the remote path

	tx, err := database.GetTxOrDB(ctx)
	if err != nil {
		return err
	}

	// Ensure only one task processes the same file at a time using a mutex.
	lockValue, _ := fileLocks.LoadOrStore(fileName, &sync.Mutex{})
	fileLock := lockValue.(*sync.Mutex)

	fileLock.Lock()
	defer fileLock.Unlock()

	// Check if the file has already been processed.
	exists, err := t.FileUC.FileExistsAndDownloaded(ctx, fileName)
	if err != nil {
		log.Printf("[ProcessFileTask] FileExistsAndDownloaded error: %v", err)
	}
	if exists {
		log.Printf("[ProcessFileTask] File already processed: %s", fileName)
		return nil
	}

	// Create a new file record in the database.
	fileModel := &payinModel.PayinFile{
		PaymentProviderID: int(payinObject.PayinFileProviderId),
		PayinFileGroupID:  &t.GroupID,
		FileName:          fileName,
		FileContentKey:    t.RemotePath,
		DownloadStatus:    payinObject.StatusPending,
		UploadStatus:      payinObject.StatusPending,
		ImportStatus:      payinObject.StatusPending,
	}

	var created *payinModel.PayinFile
	err = tx.Transaction(func(tx *gorm.DB) error {
		var err error
		created, err = t.FileUC.CreateFile(ctx, fileModel)
		if err != nil {
			log.Printf("[handleTask] CreateFile error: %v", err)
			return err
		}
		return nil
	})
	if err != nil {
		log.Printf("[handleTask] Transaction error: %v", err)
		return err
	}

	log.Printf("[ProcessFileTask] Created file record: %s (ID: %d)", fileName, created.ID)

	// Read the file content from the remote server.
	log.Printf("[ProcessFileTask] Streaming and uploading: %s", t.RemotePath)
	data, err := t.SSHClient.ReadRemoteFile(t.RemotePath)
	if err != nil {
		log.Printf("[ERROR] ReadRemoteFile failed for %s: %v", t.RemotePath, err)
		_ = t.FileUC.UpdateDownloadStatus(ctx, created, payinObject.StatusFailed)
		return err
	}
	reader := bytes.NewReader(data)
	size := int64(len(data))
	_ = t.FileUC.UpdateDownloadStatus(ctx, created, payinObject.StatusSuccess)

	// Upload the file content to S3.
	key := t.S3Uploader.GetS3KeyFromRemotePath(t.RemotePath, t.LocalPath)
	err = t.S3Uploader.UploadStreamWithContentLength(ctx, key, reader, size)
	if err != nil {
		log.Printf("[ERROR] UploadStream failed for %s: %v", t.RemotePath, err)
		_ = t.FileUC.UpdateUploadStatus(ctx, created, payinObject.StatusFailed)
		return err
	}
	_ = t.FileUC.UpdateUploadStatus(ctx, created, payinObject.StatusSuccess)
	log.Printf("[ProcessFileTask] Successfully processed file: %s", t.RemotePath)
	return nil
}
