package user

import (
	"context"
	"errors"
	"math"

	"github.com/huydq/test/internal/datastructure/inputdata"
	userModel "github.com/huydq/test/internal/domain/model/user"
	userRepo "github.com/huydq/test/internal/domain/repository/user"
	"github.com/huydq/test/internal/infrastructure/persistence/user/dto"

	"github.com/huydq/test/internal/pkg/database"
	"gorm.io/gorm"
)

// UserRepositoryImpl implements the UserRepository interface
type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *gorm.DB) userRepo.UserRepository {
	return &UserRepositoryImpl{db: db}
}

// FindByID finds a user by ID
func (r *UserRepositoryImpl) FindByID(ctx context.Context, id int) (*userModel.User, error) {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return nil, err
	}

	var userDTO dto.User
	if err := db.Preload("Role").First(&userDTO, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return userDTO.ToUserModel(), nil
}

// FindByEmail finds a user by email
func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*userModel.User, error) {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return nil, err
	}

	var userDTO dto.User
	if err := db.Where("email = ?", email).Preload("Role").First(&userDTO).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return userDTO.ToUserModel(), nil
}

// Create creates a new user
func (r *UserRepositoryImpl) Create(ctx context.Context, user *userModel.User) error {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return err
	}

	userDTO := dto.ToUserDTO(user)
	if err := db.Create(userDTO).Error; err != nil {
		return err
	}

	user.ID = userDTO.ID
	return nil
}

// Update updates an existing user
func (r *UserRepositoryImpl) Update(ctx context.Context, user *userModel.User) error {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return err
	}

	return db.Save(dto.ToUserDTO(user)).Error
}

// Delete soft-deletes a user by ID
func (r *UserRepositoryImpl) Delete(ctx context.Context, id int) error {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return err
	}

	return db.Delete(&dto.User{}, id).Error
}

// List lists users with filtering and pagination
func (r *UserRepositoryImpl) List(ctx context.Context, params *inputdata.UserListInputData) ([]*userModel.User, int, int, error) {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return nil, 0, 0, err
	}

	query := db.Model(&dto.User{})
	query = r.applyFilters(query, params)

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, 0, err
	}

	query = r.db.Model(&dto.User{})
	query = r.applyFilters(query, params)
	query = r.applyPagination(query, params)
	query = r.applySorting(query, params)
	query = query.Preload("Role")

	var userDTOs []dto.User
	if err := query.Find(&userDTOs).Error; err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(count) / float64(params.PageSize)))

	users := make([]*userModel.User, len(userDTOs))
	for i, dto := range userDTOs {
		users[i] = dto.ToUserModel()
	}

	return users, totalPages, int(count), nil
}

// applyFilters applies search and role filters to the query
func (r *UserRepositoryImpl) applyFilters(query *gorm.DB, params *inputdata.UserListInputData) *gorm.DB {
	if params.Search != "" {
		query = query.Where("user.full_name LIKE ? OR user.email LIKE ?",
			"%"+params.Search+"%", "%"+params.Search+"%")
	}

	if params.RoleID != nil && *params.RoleID != 0 {
		query = query.Where("user.role_id = ?", params.RoleID)
	}

	return query
}

// applyPagination applies pagination to the query
func (r *UserRepositoryImpl) applyPagination(query *gorm.DB, params *inputdata.UserListInputData) *gorm.DB {
	offset := (params.Page - 1) * params.PageSize
	return query.Offset(offset).Limit(params.PageSize)
}

// applySorting applies sorting to the query
func (r *UserRepositoryImpl) applySorting(query *gorm.DB, params *inputdata.UserListInputData) *gorm.DB {
	if params.SortField != "" {
		sortFieldMappings := map[string]string{
			"full_name": "user.full_name",
			"email":     "user.email",
		}

		if dbField, ok := sortFieldMappings[params.SortField]; ok {
			query = query.Order(dbField + " " + params.SortOrder)
		}
	}

	return query
}

// GetUsersWithAuditLogs retrieves users who have audit log entries
func (r *UserRepositoryImpl) GetUsersWithAuditLogs(ctx context.Context) ([]*userModel.User, error) {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return nil, err
	}

	// First get the user IDs that have audit logs
	var userIDs []int
	if err := db.WithContext(ctx).
		Table("audit_log").
		Select("DISTINCT audit_log.user_id").
		Where("audit_log.user_id IS NOT NULL").
		Pluck("user_id", &userIDs).Error; err != nil {
		return nil, err
	}

	if len(userIDs) == 0 {
		return []*userModel.User{}, nil
	}

	// Then retrieve those users with Role preloaded
	var userDTOs []dto.User
	if err := r.db.WithContext(ctx).
		Preload("Role").
		Where("id IN ?", userIDs).
		Order("full_name").
		Find(&userDTOs).Error; err != nil {
		return nil, err
	}

	users := make([]*userModel.User, len(userDTOs))
	for i, dto := range userDTOs {
		users[i] = dto.ToUserModel()
	}

	return users, nil
}
