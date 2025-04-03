package repositories

import (
	"context"

	"github.com/vnlab/makeshop-payment/src/domain/models"
	"github.com/vnlab/makeshop-payment/src/domain/repositories"
	"gorm.io/gorm"
)

// MasterAuditLogTypeRepositoryImpl implements the MasterAuditLogTypeRepository interface
type MasterAuditLogTypeRepositoryImpl struct {
	db *gorm.DB
}

// NewMasterAuditLogTypeRepository creates a new MasterAuditLogTypeRepository
func NewMasterAuditLogTypeRepository(db *gorm.DB) repositories.MasterAuditLogTypeRepository {
	return &MasterAuditLogTypeRepositoryImpl{
		db: db,
	}
}

// GetByID retrieves a master audit log type by its ID
func (r *MasterAuditLogTypeRepositoryImpl) GetByID(ctx context.Context, id int) (*models.MasterAuditLogType, error) {
	var auditLogType models.MasterAuditLogType
	err := r.db.WithContext(ctx).First(&auditLogType, id).Error
	if err != nil {
		return nil, err
	}
	return &auditLogType, nil
}

// GetByCode retrieves a master audit log type by its code
func (r *MasterAuditLogTypeRepositoryImpl) GetByCode(ctx context.Context, code string) (*models.MasterAuditLogType, error) {
	var auditLogType models.MasterAuditLogType
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&auditLogType).Error
	if err != nil {
		return nil, err
	}
	return &auditLogType, nil
}

// GetAll retrieves all master audit log types
func (r *MasterAuditLogTypeRepositoryImpl) GetAll(ctx context.Context) ([]*models.MasterAuditLogType, error) {
	var auditLogTypes []*models.MasterAuditLogType
	err := r.db.WithContext(ctx).Find(&auditLogTypes).Error
	if err != nil {
		return nil, err
	}
	return auditLogTypes, nil
} 
