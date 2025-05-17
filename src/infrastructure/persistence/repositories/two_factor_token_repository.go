package repositories

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	models "github.com/huydq/test/src/domain/models"
	repositories "github.com/huydq/test/src/domain/repositories"
)

// TwoFactorTokenRepositoryImpl implements TwoFactorTokenRepository
type TwoFactorTokenRepositoryImpl struct {
	db *gorm.DB
}

// NewTwoFactorTokenRepository creates a new TwoFactorTokenRepository
func NewTwoFactorTokenRepository(db *gorm.DB) repositories.TwoFactorTokenRepository {
	return &TwoFactorTokenRepositoryImpl{db: db}
}

// Create creates a new two-factor token
func (r *TwoFactorTokenRepositoryImpl) Create(ctx context.Context, token *models.TwoFactorToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

// FindByToken finds a token by its value and user ID
func (r *TwoFactorTokenRepositoryImpl) FindByToken(ctx context.Context, userID int, token string) (*models.TwoFactorToken, error) {
	var result models.TwoFactorToken
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND token = ? AND is_used = ? AND expired_at >= ?", userID, token, false, time.Now()).
		First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

// MarkAsUsed marks a token as used
func (r *TwoFactorTokenRepositoryImpl) MarkAsUsed(ctx context.Context, tokenID int) error {
	t := &models.TwoFactorToken{ID: tokenID}
	return r.db.WithContext(ctx).Model(t).Updates(map[string]interface{}{
		"is_used": true,
	}).Error
}

// InvalidatePreviousTokens soft-deletes any existing tokens for the user with the given MFA type
func (r *TwoFactorTokenRepositoryImpl) InvalidatePreviousTokens(ctx context.Context, userID int, mfaType int) error {
	t := &models.TwoFactorToken{UserID: userID, MFAType: mfaType, IsUsed: false}
	return r.db.WithContext(ctx).
		Model(t).
		Where("deleted_at IS NULL").
		Update("deleted_at", time.Now()).
		Error
}
