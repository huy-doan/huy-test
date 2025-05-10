package two_factor_token

import (
	"context"

	"github.com/huydq/test/internal/domain/model/two_factor_token"
)

// TwoFactorTokenRepository defines the interface for two factor token data access
type TwoFactorTokenRepository interface {
	// Create creates a new two-factor token
	Create(ctx context.Context, token *two_factor_token.TwoFactorToken) error

	// FindByToken finds a token by its value and user ID
	FindByToken(ctx context.Context, userID int, token string) (*two_factor_token.TwoFactorToken, error)

	// MarkAsUsed marks a token as used
	MarkAsUsed(ctx context.Context, tokenID int) error

	// InvalidatePreviousTokens soft-deletes any existing tokens for the user with the given MFA type
	InvalidatePreviousTokens(ctx context.Context, userID int, mfaType int) error

	// GetLastToken gets the most recently created token for a user and MFA type
	GetLastToken(ctx context.Context, userID int, mfaType int) (*two_factor_token.TwoFactorToken, error)
}
