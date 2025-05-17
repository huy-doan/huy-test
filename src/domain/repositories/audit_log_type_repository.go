package repositories

import (
	"context"

	"github.com/huydq/test/src/domain/models"
)

// AuditLogTypeRepository defines the interface for master audit log type operations
type AuditLogTypeRepository interface {
	GetByID(ctx context.Context, id int) (*models.AuditLogType, error)
	GetByCode(ctx context.Context, code string) (*models.AuditLogType, error)
	GetAll(ctx context.Context) ([]*models.AuditLogType, error)
}
