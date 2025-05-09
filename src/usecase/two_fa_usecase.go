package usecase

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	validator "github.com/huydq/test/src/api/http/validator/auth"
	models "github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/domain/repositories"
	"github.com/huydq/test/src/infrastructure/auth"
	"github.com/huydq/test/src/infrastructure/config"
)

// TwoFAUsecase handles two-factor authentication business logic
type TwoFAUsecase struct {
	userRepo       repositories.UserRepository
	twoFactorRepo  repositories.TwoFactorTokenRepository
	jwtService     *auth.JWTService
	tokenExpiryMin int
}

// NewTwoFAUsecase creates a new TwoFAUsecase
func NewTwoFAUsecase(
	userRepo repositories.UserRepository,
	twoFactorRepo repositories.TwoFactorTokenRepository,
	jwtService *auth.JWTService,
) *TwoFAUsecase {
	return &TwoFAUsecase{
		userRepo:       userRepo,
		twoFactorRepo:  twoFactorRepo,
		jwtService:     jwtService,
		tokenExpiryMin: config.GetConfig().MFATokenExpiryMinutes,
	}
}

// GenerateVerificationResponse represents a response for generating a verification token
type GenerateVerificationResponse struct {
	MFAType   int   `json:"mfa_type"`
	ExpiresIn int64 `json:"expires_in"`
}

// VerifyResponse represents a 2FA verification response
type VerifyResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

// Generate2FAToken generates a new 2FA token for a user
func (uc *TwoFAUsecase) Generate2FAToken(ctx context.Context, userID int, mfaType int) (*GenerateVerificationResponse, error) {
	// Get user
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user.not_found")
	}

	// Invalidate any previous tokens for this user with the same MFA type
	if err := uc.twoFactorRepo.InvalidatePreviousTokens(ctx, userID, mfaType); err != nil {
		return nil, err
	}

	token := fmt.Sprintf("%06d", rand.Intn(1000000))
	expiredAt := time.Now().Add(time.Duration(uc.tokenExpiryMin) * time.Minute)

	twoFactorToken := &models.TwoFactorToken{
		UserID:    userID,
		Token:     token,
		MFAType:   mfaType,
		IsUsed:    false,
		ExpiredAt: expiredAt,
	}

	if err := uc.twoFactorRepo.Create(ctx, twoFactorToken); err != nil {
		return nil, err
	}

	// TODO: Send code into user's email
	fmt.Printf("2FA Token for user %d: %s (expires at %v)\n", userID, token, expiredAt)

	return &GenerateVerificationResponse{
		MFAType:   mfaType,
		ExpiresIn: int64(uc.tokenExpiryMin * 60),
	}, nil
}

// VerifyToken verifies a 2FA token
func (uc *TwoFAUsecase) Verify2FAToken(ctx context.Context, req validator.VerifyRequest) (*VerifyResponse, error) {
	user, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user.not_found")
	}

	token, err := uc.twoFactorRepo.FindByToken(ctx, user.ID, req.Token)
	if err != nil {
		return nil, err
	}

	if token == nil {
		return nil, errors.New("mfa.invalid_token")
	}

	if !token.IsValid() {
		return nil, errors.New("mfa.expired_token")
	}

	if err := uc.twoFactorRepo.MarkAsUsed(ctx, token.ID); err != nil {
		return nil, err
	}

	jwtToken, err := uc.jwtService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &VerifyResponse{
		Token: jwtToken,
		User:  user,
	}, nil
}
