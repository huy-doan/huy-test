package repositories

import (
	"context"

	"github.com/huydq/test/src/domain/models"
)

// AuditLogRepository defines the interface for audit log operations
type AuditLogRepository interface {
	Create(ctx context.Context, auditLog *models.AuditLog) error
}
