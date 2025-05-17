package persistence

import (
	"context"

	model "github.com/huydq/test/internal/domain/model/permission"
	repository "github.com/huydq/test/internal/domain/repository/permission"
	"github.com/huydq/test/internal/infrastructure/persistence/permission/convert"
	"github.com/huydq/test/internal/infrastructure/persistence/permission/dto"
	"github.com/huydq/test/internal/pkg/database"
	"gorm.io/gorm"
)

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) repository.PermissionRepository {
	return &PermissionRepository{
		db: db,
	}
}

func (r *PermissionRepository) FindByIDs(ctx context.Context, ids []int) ([]*model.Permission, error) {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return nil, err
	}
	var permissionDTOs []*dto.Permission

	err = db.
		Preload("Screen").
		Where("id IN ?", ids).
		Find(&permissionDTOs).Error

	if err != nil {
		return nil, err
	}

	permissions := convert.ToPermissionModels(permissionDTOs)
	return permissions, nil
}

func (r *PermissionRepository) List(ctx context.Context) ([]*model.Permission, error) {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return nil, err
	}
	var permissionDTOs []*dto.Permission

	err = db.
		Preload("Screen").
		Find(&permissionDTOs).Error

	if err != nil {
		return nil, err
	}

	permissions := convert.ToPermissionModels(permissionDTOs)
	return permissions, nil
}
