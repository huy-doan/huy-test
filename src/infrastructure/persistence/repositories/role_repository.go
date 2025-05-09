package repositories

import (
	"context"
	"errors"

	models "github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/domain/repositories"
	"gorm.io/gorm"
)

// RoleRepositoryImpl implements the RoleRepository interface
type RoleRepositoryImpl struct {
	db *gorm.DB
}

// NewRoleRepository creates a new RoleRepository
func NewRoleRepository(db *gorm.DB) repositories.RoleRepository {
	return &RoleRepositoryImpl{
		db: db,
	}
}

// FindByID finds a role by ID
func (r *RoleRepositoryImpl) FindByID(ctx context.Context, id int) (*models.Role, error) {
	var role models.Role
	result := r.db.First(&role, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if role not found
		}
		return nil, result.Error
	}
	return &role, nil
}

// FindByCode finds a role by code
func (r *RoleRepositoryImpl) FindByCode(ctx context.Context, code string) (*models.Role, error) {
	var role models.Role
	result := r.db.Where("code = ?", code).First(&role)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if role not found
		}
		return nil, result.Error
	}
	return &role, nil
}
