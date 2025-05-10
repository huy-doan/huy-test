package persistence

import (
	"context"

	model "github.com/huydq/test/internal/domain/model/role"
	repository "github.com/huydq/test/internal/domain/repository/role"
	permissionDto "github.com/huydq/test/internal/infrastructure/persistence/permission/dto"
	"github.com/huydq/test/internal/infrastructure/persistence/role/convert"
	"github.com/huydq/test/internal/infrastructure/persistence/role/dto"
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
	var roleDTO dto.RoleDTO

	err := r.db.WithContext(ctx).
		Preload("Permissions.Screen").
		First(&roleDTO, id).Error

	if err != nil {
		return nil, err
	}

	role := convert.ToRoleModel(&roleDTO)
	return role, nil
}

func (r *RoleRepositoryImpl) FindByCode(ctx context.Context, code string) (*model.Role, error) {
	var roleDTO dto.RoleDTO

	err := r.db.WithContext(ctx).
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
	var roleDTO dto.RoleDTO

	err := r.db.WithContext(ctx).
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
	roleDTO := convert.ToRoleDTO(role)
	result := r.db.WithContext(ctx).Create(roleDTO)

	// Update the ID in the original model
	if result.Error == nil {
		role.ID = roleDTO.ID
	}

	return result.Error
}

func (r *RoleRepositoryImpl) Update(ctx context.Context, role *model.Role) error {
	roleDTO := convert.ToRoleDTO(role)

	// Start transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Update basic role information
	if err := tx.Model(&dto.RoleDTO{}).Where("id = ?", role.ID).Updates(map[string]interface{}{
		"name": role.Name,
		"code": string(role.Code),
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update permissions if present
	if role.Permissions != nil {
		// Clear existing associations
		if err := tx.Model(&roleDTO).Association("Permissions").Clear(); err != nil {
			tx.Rollback()
			return err
		}

		// Add new associations
		if len(role.Permissions) > 0 {
			var permissionIDs []int
			for _, perm := range role.Permissions {
				permissionIDs = append(permissionIDs, perm.ID)
			}

			if err := tx.Model(&roleDTO).Association("Permissions").Append(
				tx.Where("id IN ?", permissionIDs).Find(&[]*permissionDto.PermissionDTO{}),
			); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

func (r *RoleRepositoryImpl) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&dto.RoleDTO{}, id).Error
}

func (r *RoleRepositoryImpl) List(ctx context.Context) ([]*model.Role, error) {
	var roleDTOs []*dto.RoleDTO

	err := r.db.WithContext(ctx).
		Preload("Permissions.Screen").
		Find(&roleDTOs).Error

	if err != nil {
		return nil, err
	}

	roles := convert.ToRoleModels(roleDTOs)
	return roles, nil
}
