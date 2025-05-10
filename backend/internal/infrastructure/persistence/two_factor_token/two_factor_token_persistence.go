package two_factor_token

import (
	"context"
	"errors"
	"time"

	twoFactorModel "github.com/huydq/test/internal/domain/model/two_factor_token"
	twoFactorRepo "github.com/huydq/test/internal/domain/repository/two_factor_token"
	"github.com/huydq/test/internal/infrastructure/persistence/two_factor_token/dto"
	"gorm.io/gorm"
)

// TwoFactorTokenRepositoryImpl implements the TwoFactorTokenRepository interface
type TwoFactorTokenRepositoryImpl struct {
	db *gorm.DB
}

// NewTwoFactorTokenRepository creates a new TwoFactorTokenRepository
func NewTwoFactorTokenRepository(db *gorm.DB) twoFactorRepo.TwoFactorTokenRepository {
	return &TwoFactorTokenRepositoryImpl{db: db}
}

// Create creates a new two-factor token
func (r *TwoFactorTokenRepositoryImpl) Create(ctx context.Context, token *twoFactorModel.TwoFactorToken) error {
	tokenDTO := dto.ToTwoFactorTokenDTO(token)
	return r.db.WithContext(ctx).Create(tokenDTO).Error
}

// FindByToken finds a token by its value and user ID
func (r *TwoFactorTokenRepositoryImpl) FindByToken(ctx context.Context, userID int, token string) (*twoFactorModel.TwoFactorToken, error) {
	var tokenDTO dto.TwoFactorTokenDTO
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND token = ? AND is_used = ? AND expired_at >= ?", userID, token, false, time.Now()).
		First(&tokenDTO).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return tokenDTO.ToTwoFactorTokenModel(), nil
}

// MarkAsUsed marks a token as used
func (r *TwoFactorTokenRepositoryImpl) MarkAsUsed(ctx context.Context, tokenID int) error {
	tokenDTO := &dto.TwoFactorTokenDTO{ID: tokenID}
	return r.db.WithContext(ctx).Model(tokenDTO).Updates(map[string]any{
		"is_used": true,
	}).Error
}

// InvalidatePreviousTokens soft-deletes any existing tokens for the user with the given MFA type
func (r *TwoFactorTokenRepositoryImpl) InvalidatePreviousTokens(ctx context.Context, userID int, mfaType int) error {
	tokenDTO := &dto.TwoFactorTokenDTO{UserID: userID, MFAType: mfaType, IsUsed: false}
	return r.db.WithContext(ctx).
		Model(tokenDTO).
		Where("deleted_at IS NULL").
		Update("deleted_at", time.Now()).
		Error
}

// GetLastToken gets the most recently created token for a user and MFA type
func (r *TwoFactorTokenRepositoryImpl) GetLastToken(ctx context.Context, userID int, mfaType int) (*twoFactorModel.TwoFactorToken, error) {
	var tokenDTO dto.TwoFactorTokenDTO
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND mfa_type = ?", userID, mfaType).
		Order("id DESC").
		First(&tokenDTO).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return tokenDTO.ToTwoFactorTokenModel(), nil
}
