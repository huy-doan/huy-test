package repositories

import (
	"context"

	models "github.com/huydq/test/src/domain/models"
)

type TwoFactorTokenRepository interface {
	Create(ctx context.Context, token *models.TwoFactorToken) error

	FindByToken(ctx context.Context, userID int, token string) (*models.TwoFactorToken, error)

	MarkAsUsed(ctx context.Context, tokenID int) error

	InvalidatePreviousTokens(ctx context.Context, userID int, mfaType int) error
}
