package repositories

import (
	"context"

	"github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/domain/repositories/filter"
)

// AuditLogRepository defines the interface for audit log operations
type AuditLogRepository interface {
	Create(ctx context.Context, auditLog *models.AuditLog) error
	List(ctx context.Context, filter *filter.AuditLogFilter) ([]*models.AuditLog, int, int64, error)
}
