package persistence

import (
	"context"

	model "github.com/huydq/test/internal/domain/model/permission"
	repository "github.com/huydq/test/internal/domain/repository/permission"
	"github.com/huydq/test/internal/infrastructure/persistence/permission/convert"
	"github.com/huydq/test/internal/infrastructure/persistence/permission/dto"
	"gorm.io/gorm"
)

type PermissionRepositoryImpl struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) repository.PermissionRepository {
	return &PermissionRepositoryImpl{
		db: db,
	}
}

func (r *PermissionRepositoryImpl) FindByIDs(ctx context.Context, ids []int) ([]*model.Permission, error) {
	var permissionDTOs []*dto.PermissionDTO

	err := r.db.WithContext(ctx).
		Preload("Screen").
		Where("id IN ?", ids).
		Find(&permissionDTOs).Error

	if err != nil {
		return nil, err
	}

	permissions := convert.ToPermissionModels(permissionDTOs)
	return permissions, nil
}

func (r *PermissionRepositoryImpl) List(ctx context.Context) ([]*model.Permission, error) {
	var permissionDTOs []*dto.PermissionDTO

	err := r.db.WithContext(ctx).
		Preload("Screen").
		Find(&permissionDTOs).Error

	if err != nil {
		return nil, err
	}

	permissions := convert.ToPermissionModels(permissionDTOs)
	return permissions, nil
}

func (r *PermissionRepositoryImpl) GetPermissionCodesByRoleID(ctx context.Context, roleID int) ([]string, error) {
	var permissionCodes []string
	err := r.db.WithContext(ctx).
		Model(&dto.PermissionDTO{}).
		Select("permission.code").
		Joins("JOIN role_permission ON role_permission.permission_id = permission.id").
		Where("role_permission.role_id = ? AND role_permission.deleted_at IS NULL", roleID).
		Where("permission.deleted_at IS NULL").
		Pluck("permission.code", &permissionCodes).
		Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return permissionCodes, nil
}
