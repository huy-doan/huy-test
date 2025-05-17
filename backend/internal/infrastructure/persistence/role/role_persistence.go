package persistence

import (
	"context"

	model "github.com/huydq/test/internal/domain/model/role"
	repository "github.com/huydq/test/internal/domain/repository/role"
	"github.com/huydq/test/internal/infrastructure/persistence/role/convert"
	"github.com/huydq/test/internal/infrastructure/persistence/role/dto"
	"github.com/huydq/test/internal/pkg/database"
	"gorm.io/gorm"
)

type RoleRepositoryImpl struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) repository.RoleRepository {
	return &RoleRepositoryImpl{
		db: db,
	}
}

func (r *RoleRepositoryImpl) FindByID(ctx context.Context, id int) (*model.Role, error) {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return nil, err
	}
	var roleDTO dto.Role
	err = db.
		Preload("Permissions.Screen").
		First(&roleDTO, id).Error

	if err != nil {
		return nil, err
	}

	role := convert.ToRoleModel(&roleDTO)
	return role, nil
}

func (r *RoleRepositoryImpl) FindByCode(ctx context.Context, code string) (*model.Role, error) {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return nil, err
	}

	var roleDTO dto.Role
	err = db.
		Preload("Permissions.Screen").
		Where("code = ?", code).
		First(&roleDTO).Error

	if err != nil {
		return nil, err
	}

	role := convert.ToRoleModel(&roleDTO)
	return role, nil
}

func (r *RoleRepositoryImpl) FindByName(ctx context.Context, name string) (*model.Role, error) {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return nil, err
	}

	var roleDTO dto.Role
	err = db.
		Preload("Permissions.Screen").
		Where("name = ?", name).
		First(&roleDTO).Error

	if err != nil {
		return nil, err
	}

	role := convert.ToRoleModel(&roleDTO)
	return role, nil
}

func (r *RoleRepositoryImpl) Create(ctx context.Context, role *model.Role) error {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return err
	}
	roleDTO := convert.ToRoleDTO(role)
	result := db.Create(roleDTO)

	if result.Error == nil {
		role.ID = roleDTO.ID
	}

	return result.Error
}

func (r *RoleRepositoryImpl) Update(ctx context.Context, role *model.Role) error {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return err
	}

	roleDTO := convert.ToRoleDTO(role)
	if result := db.Save(roleDTO); result.Error != nil {
		return result.Error
	}
	err = db.Model(&roleDTO).Association("Permissions").Replace(roleDTO.Permissions)
	if err != nil {
		return err
	}
	return nil
}

func (r *RoleRepositoryImpl) Delete(ctx context.Context, id int) error {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return err
	}

	return db.Delete(&dto.Role{}, id).Error
}

func (r *RoleRepositoryImpl) List(ctx context.Context) ([]*model.Role, error) {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return nil, err
	}

	var roleDTOs []*dto.Role
	err = db.
		Preload("Permissions.Screen").
		Find(&roleDTOs).Error

	if err != nil {
		return nil, err
	}

	roles := convert.ToRoleModels(roleDTOs)
	return roles, nil
}
