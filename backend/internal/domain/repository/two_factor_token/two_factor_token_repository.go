package model

import (
	"context"

	model "github.com/huydq/test/internal/domain/model/two_factor_token"
)

// TwoFactorTokenRepository defines the interface for two factor token data access
type TwoFactorTokenRepository interface {
	// Create creates a new two-factor token
	Create(ctx context.Context, token *model.TwoFactorToken) error

	// FindByToken finds a token by its value and user ID
	FindByToken(ctx context.Context, criteria model.TwoFactorToken) (*model.TwoFactorToken, error)

	// MarkAsUsed marks a token as used
	MarkAsUsed(ctx context.Context, token *model.TwoFactorToken) error

	// InvalidatePreviousTokens soft-deletes any existing tokens for the user with the given MFA type
	InvalidatePreviousTokens(ctx context.Context, criteria model.TwoFactorToken) error

	// GetLastToken gets the most recently created token for a user and MFA type
	GetLastToken(ctx context.Context, criteria model.TwoFactorToken) (*model.TwoFactorToken, error)
}
