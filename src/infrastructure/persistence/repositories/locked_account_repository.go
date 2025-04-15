package repositories

import (
	"context"

	"github.com/vnlab/makeshop-payment/src/domain/models"
	"github.com/vnlab/makeshop-payment/src/domain/repositories"
	"gorm.io/gorm"
)

type lockedAccountRepository struct {
	db *gorm.DB
}

// NewLockedAccountRepository creates a new instance of LockedAccountRepository
func NewLockedAccountRepository(db *gorm.DB) repositories.LockedAccountRepository {
	return &lockedAccountRepository{db: db}
}

// Create implements repositories.LockedAccountRepository
func (r *lockedAccountRepository) Create(ctx context.Context, account *models.LockedAccount) error {
	return r.db.WithContext(ctx).Create(account).Error
}

// Update implements repositories.LockedAccountRepository
func (r *lockedAccountRepository) Update(ctx context.Context, account *models.LockedAccount) error {
	return r.db.WithContext(ctx).Save(account).Error
}

// GetByEmail implements repositories.LockedAccountRepository
func (r *lockedAccountRepository) GetByEmail(ctx context.Context, email string) (*models.LockedAccount, error) {
	var account models.LockedAccount
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// GetByUserID implements repositories.LockedAccountRepository
func (r *lockedAccountRepository) GetByUserID(ctx context.Context, userID int) (*models.LockedAccount, error) {
	var account models.LockedAccount
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// GetByID implements repositories.LockedAccountRepository
func (r *lockedAccountRepository) GetByID(ctx context.Context, id int) (*models.LockedAccount, error) {
	var account models.LockedAccount
	err := r.db.WithContext(ctx).First(&account, id).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// List implements repositories.LockedAccountRepository
func (r *lockedAccountRepository) List(ctx context.Context, page, pageSize int) ([]*models.LockedAccount, int, error) {
	var accounts []*models.LockedAccount
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.LockedAccount{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate total pages
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	// Get paginated records
	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(pageSize).
		Order("id DESC").
		Find(&accounts).Error

	if err != nil {
		return nil, 0, err
	}

	return accounts, totalPages, nil
}

// Delete removes a locked account record
func (r *lockedAccountRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.LockedAccount{}, "id = ?", id).Error
}
