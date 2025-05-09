package repositories

import (
	"context"

	"github.com/huydq/test/src/domain/models"
)

type PayinFileRepository interface {
	// Create a new PayinFile record
	Create(ctx context.Context, file *models.PayinFile) error

	// UpdateDownloadStatus updates the download status of a PayinFile record
	UpdateStatus(ctx context.Context, id int, field string, status int) error

	// FindIDByFilename checks if a file exists by its filename and returns its ID
	FindIDByFilename(ctx context.Context, filename string) (int, error)

	// GetDownloadStatusByID retrieves the download status of a PayinFile record by its ID
	GetDownloadStatusByID(ctx context.Context, id int) (int, error)

	// GetByID retrieves a PayinFile record by its ID
	GetByID(ctx context.Context, id int) (*models.PayinFile, error)
}
