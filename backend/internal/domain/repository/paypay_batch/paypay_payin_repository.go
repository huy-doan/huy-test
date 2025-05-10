package repositories

import (
	"context"

	payinFileModel "github.com/huydq/test/internal/infrastructure/persistence/paypay_batch/paypay_payin_file/dto"
)

type PayinFileRepository interface {
	// Create a new PayinFile record
	Create(ctx context.Context, file *payinFileModel.PayinFile) (*payinFileModel.PayinFile, error)

	// UpdateStatus updates the status of a PayinFile record
	UpdateStatus(ctx context.Context, id int, field string, status int) (*payinFileModel.PayinFile, error)

	// FindByFilename checks if a file exists by its filename and returns its PayinFile model
	FindByFilename(ctx context.Context, filename string) (*payinFileModel.PayinFile, error)

	// GetByID retrieves a PayinFile record by its ID
	GetByID(ctx context.Context, id int) (*payinFileModel.PayinFile, error)
}
