package two_fa

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/huydq/test/internal/datastructure/inputdata"
	"github.com/huydq/test/internal/datastructure/outputdata"
	twoFactorModel "github.com/huydq/test/internal/domain/model/two_factor_token"
	userModel "github.com/huydq/test/internal/domain/model/user"
	twoFactorRepo "github.com/huydq/test/internal/domain/repository/two_factor_token"
	userRepo "github.com/huydq/test/internal/domain/repository/user"
	"github.com/huydq/test/internal/infrastructure/adapter/auth"
	"github.com/huydq/test/internal/infrastructure/adapter/email"

	"github.com/huydq/test/internal/pkg/config"
	"github.com/huydq/test/internal/pkg/database"
)

// TwoFAUsecase interface defines the contract for two-factor authentication use cases
type TwoFAUsecase interface {
	Generate2FAToken(ctx context.Context, input *inputdata.GenerateTwoFAInputData) (*outputdata.GenerateTwoFAOutputData, error)
	Verify2FAToken(ctx context.Context, input *inputdata.VerifyTwoFAInputData) (*outputdata.VerifyTwoFAOutputData, error)
	CanResendCode(ctx context.Context, input *inputdata.CanResendCodeInputData) (*outputdata.CanResendCodeOutputData, error)
}

// TwoFAUsecaseImpl implements the TwoFAUsecase interface
type TwoFAUsecaseImpl struct {
	userRepo       userRepo.UserRepository
	twoFactorRepo  twoFactorRepo.TwoFactorTokenRepository
	jwtService     *auth.JWTService
	mailService    *email.MailService
	tokenExpiryMin int
}

// NewTwoFAUsecase creates a new TwoFAUsecase implementation
func NewTwoFAUsecase(
	userRepo userRepo.UserRepository,
	twoFactorRepo twoFactorRepo.TwoFactorTokenRepository,
	jwtService *auth.JWTService,
	mailService *email.MailService,
) TwoFAUsecase {
	return &TwoFAUsecaseImpl{
		userRepo:       userRepo,
		twoFactorRepo:  twoFactorRepo,
		jwtService:     jwtService,
		mailService:    mailService,
		tokenExpiryMin: config.GetConfig().MFATokenExpiryMinutes,
	}
}

// Generate2FAToken generates a new 2FA token for a user
func (uc *TwoFAUsecaseImpl) Generate2FAToken(ctx context.Context, input *inputdata.GenerateTwoFAInputData) (*outputdata.GenerateTwoFAOutputData, error) {
	user, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user.not_found")
	}

	if err := uc.twoFactorRepo.InvalidatePreviousTokens(ctx, input.UserID, input.MFAType); err != nil {
		return nil, err
	}

	token := fmt.Sprintf("%06d", rand.Intn(1000000))
	expiredAt := time.Now().Add(time.Duration(uc.tokenExpiryMin) * time.Minute)

	twoFactorToken := &twoFactorModel.TwoFactorToken{
		UserID:    input.UserID,
		Token:     token,
		MFAType:   input.MFAType,
		IsUsed:    false,
		ExpiredAt: expiredAt,
	}

	if err := uc.twoFactorRepo.Create(ctx, twoFactorToken); err != nil {
		return nil, err
	}

	uc.sendVerificationEmail(user, token)

	return &outputdata.GenerateTwoFAOutputData{
		MFAType:   input.MFAType,
		ExpiresIn: int64(uc.tokenExpiryMin * 60),
	}, nil
	tx, err := database.NewTx[*outputdata.GenerateTwoFAOutputData](ctx)
	if err != nil {
		return nil, err
	}

	result, err := tx.Transact(ctx, func(ctx context.Context) (*outputdata.GenerateTwoFAOutputData, error) {
		if err := uc.twoFactorRepo.InvalidatePreviousTokens(ctx, input.UserID, input.MFAType); err != nil {
			return nil, err
		}

		token := fmt.Sprintf("%06d", rand.Intn(1000000))
		expiredAt := time.Now().Add(time.Duration(uc.tokenExpiryMin) * time.Minute)

		twoFactorToken := &twoFactorModel.TwoFactorToken{
			UserID:    input.UserID,
			Token:     token,
			MFAType:   input.MFAType,
			IsUsed:    false,
			ExpiredAt: expiredAt,
		}

		if err := uc.twoFactorRepo.Create(ctx, twoFactorToken); err != nil {
			return nil, err
		}

		uc.sendVerificationEmail(user, token)

		return &outputdata.GenerateTwoFAOutputData{
			MFAType:   input.MFAType,
			ExpiresIn: int64(uc.tokenExpiryMin * 60),
		}, nil
	})

	return result, err
}

// Verify2FAToken verifies a 2FA token
func (uc *TwoFAUsecaseImpl) Verify2FAToken(ctx context.Context, input *inputdata.VerifyTwoFAInputData) (*outputdata.VerifyTwoFAOutputData, error) {
	user, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user.not_found")
	}

	token, err := uc.twoFactorRepo.FindByToken(ctx, user.ID, input.Token)
	if err != nil {
		return nil, err
	}

	if token == nil {
		return nil, errors.New("mfa.invalid_token")
	}

	if token.ExpiredAt.Before(time.Now()) {
		return nil, errors.New("mfa.expired_token")
	}

	if err := uc.twoFactorRepo.MarkAsUsed(ctx, token.ID); err != nil {
		return nil, err
	}

	jwtToken, err := uc.jwtService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &outputdata.VerifyTwoFAOutputData{
		Token: jwtToken,
		User:  user,
	}, nil
}

// CanResendCode checks if a user can resend a verification code
func (uc *TwoFAUsecaseImpl) CanResendCode(ctx context.Context, input *inputdata.CanResendCodeInputData) (*outputdata.CanResendCodeOutputData, error) {
	const resendInterval = 1

	lastToken, err := uc.twoFactorRepo.GetLastToken(ctx, input.UserID, input.MFAType)
	if err != nil {
		return nil, err
	}

	if lastToken == nil {
		return &outputdata.CanResendCodeOutputData{
			CanResend:     true,
			RemainingTime: 0,
		}, nil
	}

	earliestNextResendTime := time.Now().Add(-resendInterval * time.Minute)
	remainingTime := time.Until(lastToken.CreatedAt.Add(time.Duration(resendInterval) * time.Minute))
	return &outputdata.CanResendCodeOutputData{
		CanResend:     lastToken.CreatedAt.Before(earliestNextResendTime),
		RemainingTime: int(remainingTime.Seconds()),
	}, nil
}

// Helper method to send verification email
func (uc *TwoFAUsecaseImpl) sendVerificationEmail(user *userModel.User, token string) {
	uc.mailService.SendMailByTemplateID(email.TemplateID2FACode, email.TwoFACodeEmailData{
		Email:          user.Email,
		ToName:         user.FullName,
		Token:          token,
		TokenExpiryMin: uc.tokenExpiryMin,
	})
}
