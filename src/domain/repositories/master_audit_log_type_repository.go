package repositories

import (
	"context"
	"github.com/vnlab/makeshop-payment/src/domain/models"
)

// MasterAuditLogTypeRepository defines the interface for master audit log type operations
type MasterAuditLogTypeRepository interface {
	GetByID(ctx context.Context, id int) (*models.MasterAuditLogType, error)
	GetByCode(ctx context.Context, code string) (*models.MasterAuditLogType, error)
	GetAll(ctx context.Context) ([]*models.MasterAuditLogType, error)
} 
