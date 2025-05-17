package repository

import (
	"context"

	model "github.com/huydq/test/internal/domain/model/payin"
)

type PayinFileRepository interface {
	// Create a new PayinFile record
	Create(ctx context.Context, file *model.PayinFile) (*model.PayinFile, error)

	// UpdateStatus updates the status of a PayinFile record
	UpdateStatus(ctx context.Context, file *model.PayinFile) error

	// FindByFilename checks if a file exists by its filename and returns its PayinFile model
	FindByFilename(ctx context.Context, filename string) (*model.PayinFile, error)

	// GetByID retrieves a PayinFile record by its ID
	GetByID(ctx context.Context, id int) (*model.PayinFile, error)
}
