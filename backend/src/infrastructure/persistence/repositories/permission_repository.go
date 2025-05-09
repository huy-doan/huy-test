package repositories

import (
	"context"

	"github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/domain/repositories"
	"gorm.io/gorm"
)

type PermissionRepositoryImpl struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) repositories.PermissionRepository {
	return &PermissionRepositoryImpl{
		db: db,
	}
}

func (r *PermissionRepositoryImpl) FindByIDs(ctx context.Context, ids []int) ([]*models.Permission, error) {
	var permissions []*models.Permission
	result := r.db.Where("id IN ?", ids).
		Preload("Screen").
		Find(&permissions)
	if result.Error != nil {
		return nil, result.Error
	}
	return permissions, nil
}

func (r *PermissionRepositoryImpl) List(ctx context.Context) ([]*models.Permission, error) {
	var permissions []*models.Permission
	result := r.db.
		Preload("Screen").
		Find(&permissions)
	if result.Error != nil {
		return nil, result.Error
	}
	return permissions, nil
}
