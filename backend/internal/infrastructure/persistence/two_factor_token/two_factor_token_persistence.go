package persistence

import (
	"context"
	"errors"

	twoFactorModel "github.com/huydq/test/internal/domain/model/two_factor_token"
	twoFactorRepo "github.com/huydq/test/internal/domain/repository/two_factor_token"
	"github.com/huydq/test/internal/infrastructure/persistence/two_factor_token/dto"
	"github.com/huydq/test/internal/pkg/database"
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
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return err
	}
	tokenDTO := dto.ToTwoFactorTokenDTO(token)
	return db.Create(tokenDTO).Error
}

// FindByToken finds a token by its value and user ID
func (r *TwoFactorTokenRepositoryImpl) FindByToken(ctx context.Context, criteria twoFactorModel.TwoFactorToken) (*twoFactorModel.TwoFactorToken, error) {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return nil, err
	}

	var tokenDTO dto.TwoFactorToken
	err = db.Where("user_id = ? AND token = ? AND is_used = ? AND expired_at >= ?",
		criteria.UserID, criteria.Token, criteria.IsUsed, criteria.ExpiredAt).
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
func (r *TwoFactorTokenRepositoryImpl) MarkAsUsed(ctx context.Context, token *twoFactorModel.TwoFactorToken) error {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return err
	}

	tokenDTO := dto.ToTwoFactorTokenDTO(token)
	return db.Model(&dto.TwoFactorToken{}).
		Select("is_used").
		Where("id = ?", token.ID).
		Updates(tokenDTO).Error
}

// InvalidatePreviousTokens soft-deletes any existing tokens for the user with the given MFA type
func (r *TwoFactorTokenRepositoryImpl) InvalidatePreviousTokens(ctx context.Context, criteria twoFactorModel.TwoFactorToken) error {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return err
	}

	return db.WithContext(ctx).
		Where("user_id = ? AND mfa_type = ? AND is_used = ?",
			criteria.UserID, criteria.MFAType, criteria.IsUsed).
		Delete(&dto.TwoFactorToken{}).
		Error
}

// GetLastToken gets the most recently created token for a user and MFA type
func (r *TwoFactorTokenRepositoryImpl) GetLastToken(ctx context.Context, criteria twoFactorModel.TwoFactorToken) (*twoFactorModel.TwoFactorToken, error) {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return nil, err
	}

	var tokenDTO dto.TwoFactorToken
	err = db.Where("user_id = ? AND mfa_type = ?", criteria.UserID, criteria.MFAType).
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
