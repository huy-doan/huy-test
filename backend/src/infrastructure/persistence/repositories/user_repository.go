package repositories

import (
	"context"
	"errors"
	"math"

	validator "github.com/huydq/test/src/api/http/validator/user"
	models "github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/domain/repositories"
	"gorm.io/gorm"
)

// UserRepositoryImpl implements the UserRepository interface
type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

// FindByID finds a user by ID
func (r *UserRepositoryImpl) FindByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	result := r.db.Preload("Role").
		Preload("Role.Permissions").
		First(&user, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if user not found
		}
		return nil, result.Error
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	result := r.db.Preload("Role").
		Where("email = ?", email).
		First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if user not found
		}
		return nil, result.Error
	}
	return &user, nil
}

// Create creates a new user
func (r *UserRepositoryImpl) Create(ctx context.Context, user *models.User) error {
	return r.db.Create(user).Error
}

// Update updates an existing user
func (r *UserRepositoryImpl) Update(ctx context.Context, user *models.User) error {
	return r.db.Save(user).Error
}

// Delete soft-deletes a user by ID
func (r *UserRepositoryImpl) Delete(ctx context.Context, id int) error {
	return r.db.Delete(&models.User{}, id).Error
}

// List lists users with filtering and pagination
func (r *UserRepositoryImpl) List(ctx context.Context, params validator.UserListFilter) ([]*models.User, int, int, error) {
	var users []*models.User
	var count int64
	query := r.db.Model(&models.User{}).Joins("Role")

	if params.Search != "" {
		query = query.Where("user.full_name LIKE ? OR user.email LIKE ? OR Role.name LIKE ?",
			"%"+params.Search+"%",
			"%"+params.Search+"%",
			"%"+params.Search+"%")
	}

	if params.RoleID != nil && *params.RoleID != 0 {
		query = query.Where("user.role_id = ?", params.RoleID)
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	if params.SortField != "" {
		sortFieldMappings := map[string]string{
			"full_name": "user.full_name",
			"email":     "user.email",
		}

		if dbField, ok := sortFieldMappings[params.SortField]; ok {
			query = query.Order(dbField + " " + params.SortOrder)
		}
	}

	if err := query.Offset(offset).
		Preload("Role").
		Limit(params.PageSize).
		Find(&users).Error; err != nil {

		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(count) / float64(params.PageSize)))

	return users, totalPages, int(count), nil
}

// GetUsersWithAuditLogs retrieves users who have audit log entries
func (r *UserRepositoryImpl) GetUsersWithAuditLogs(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	query := r.db.WithContext(ctx).
		Distinct("user.*").
		Table("user").
		Select("user.*").
		Joins("JOIN audit_log ON user.id = audit_log.user_id").
		Where("audit_log.user_id IS NOT NULL").
		Where("user.deleted_at IS NULL").
		Order("user.full_name")

	// Execute the query
	err := query.Preload("Role").Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}
