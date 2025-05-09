package repositories

import (
	"context"
	"errors"
	"math"

	models "github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/domain/repositories"
	domainFilter "github.com/huydq/test/src/domain/repositories/filter"
	infraFilter "github.com/huydq/test/src/infrastructure/persistence/repositories/filter"
	"gorm.io/gorm"
)

// RoleRepositoryImpl implements the RoleRepository interface
type RoleRepositoryImpl struct {
	db            *gorm.DB
	filterBuilder *infraFilter.GormFilterBuilder
}

// NewRoleRepository creates a new RoleRepository
func NewRoleRepository(db *gorm.DB) repositories.RoleRepository {
	return &RoleRepositoryImpl{
		db:            db,
		filterBuilder: infraFilter.NewGormFilterBuilder(),
	}
}

// FindByID finds a role by ID
func (r *RoleRepositoryImpl) FindByID(ctx context.Context, id int) (*models.Role, error) {
	var role models.Role
	result := r.db.Preload("Permissions").First(&role, id)
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
	result := r.db.Where("code = ?", code).Preload("Permissions").First(&role)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if role not found
		}
		return nil, result.Error
	}
	return &role, nil
}

// FindByName finds a role by name
func (r *RoleRepositoryImpl) FindByName(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role
	result := r.db.Where("name = ?", name).First(&role)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if role not found
		}
		return nil, result.Error
	}
	return &role, nil
}

func (r *RoleRepositoryImpl) Create(ctx context.Context, role *models.Role) error {
	result := r.db.Create(role)
	return result.Error
}

func (r *RoleRepositoryImpl) Update(ctx context.Context, role *models.Role) error {
	if result := r.db.Save(role); result.Error != nil {
		return result.Error
	}
	err := r.db.Model(&role).Association("Permissions").Replace(role.Permissions)
	if err != nil {
		return err
	}
	return nil
}

func (r *RoleRepositoryImpl) Delete(ctx context.Context, id int) error {
	result := r.db.Delete(&models.Role{}, id)
	return result.Error
}

func (r *RoleRepositoryImpl) List(ctx context.Context, filter *domainFilter.RoleFilter) ([]*models.Role, int, int64, error) {
	var roles []*models.Role
	var count int64

	// Count total records
	if err := r.db.Model(&models.Role{}).Count(&count).Error; err != nil {
		return nil, 0, 0, err
	}

	if filter != nil {
		filter.ApplyFilters()
	} else {
		filter = domainFilter.NewRoleFilter()
	}

	query := r.db.WithContext(ctx).Model(&models.Role{})
	query = r.filterBuilder.ApplyBaseFilter(query, &filter.BaseFilter)

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, 0, err
	}

	query = r.filterBuilder.ApplyPagination(query, filter.Pagination)

	query = query.Preload("Permissions")

	if err := query.Find(&roles).Error; err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(count) / float64(filter.Pagination.PageSize)))

	return roles, totalPages, int64(count), nil
}
