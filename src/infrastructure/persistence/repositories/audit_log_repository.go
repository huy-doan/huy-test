package repositories

import (
	"context"

	"github.com/vnlab/makeshop-payment/src/domain/models"
	"github.com/vnlab/makeshop-payment/src/domain/repositories"
	"gorm.io/gorm"
)

type AuditLogRepositoryImpl struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) repositories.AuditLogRepository {
	return &AuditLogRepositoryImpl{
		db: db,
	}
}

func (r *AuditLogRepositoryImpl) Create(ctx context.Context, auditLog *models.AuditLog) error {
	return r.db.Create(auditLog).Error
}
