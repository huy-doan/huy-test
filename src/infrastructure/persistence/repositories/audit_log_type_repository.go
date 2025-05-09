package repositories

import (
	"context"

	"github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/domain/repositories"
	"gorm.io/gorm"
)

// AuditLogTypeRepositoryImpl implements the AuditLogTypeRepository interface
type AuditLogTypeRepositoryImpl struct {
	db *gorm.DB
}

// NewAuditLogTypeRepository creates a new AuditLogTypeRepository
func NewAuditLogTypeRepository(db *gorm.DB) repositories.AuditLogTypeRepository {
	return &AuditLogTypeRepositoryImpl{
		db: db,
	}
}

// GetByID retrieves a master audit log type by its ID
func (r *AuditLogTypeRepositoryImpl) GetByID(ctx context.Context, id int) (*models.AuditLogType, error) {
	var auditLogType models.AuditLogType
	err := r.db.WithContext(ctx).First(&auditLogType, id).Error
	if err != nil {
		return nil, err
	}
	return &auditLogType, nil
}

// GetByCode retrieves a master audit log type by its code
func (r *AuditLogTypeRepositoryImpl) GetByCode(ctx context.Context, code string) (*models.AuditLogType, error) {
	var auditLogType models.AuditLogType
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&auditLogType).Error
	if err != nil {
		return nil, err
	}
	return &auditLogType, nil
}

// GetAll retrieves all master audit log types
func (r *AuditLogTypeRepositoryImpl) GetAll(ctx context.Context) ([]*models.AuditLogType, error) {
	var auditLogTypes []*models.AuditLogType
	err := r.db.WithContext(ctx).Find(&auditLogTypes).Error
	if err != nil {
		return nil, err
	}
	return auditLogTypes, nil
}
