package auth

import (
	"context"
	"errors"
	"time"

	"github.com/huydq/test/internal/datastructure/inputdata"
	"github.com/huydq/test/internal/datastructure/outputdata"
	twoFactorModel "github.com/huydq/test/internal/domain/model/two_factor_token"
	object "github.com/huydq/test/internal/domain/object/mfa_type"
	twoFactorRepo "github.com/huydq/test/internal/domain/repository/two_factor_token"

	userRepo "github.com/huydq/test/internal/domain/repository/user"
	authService "github.com/huydq/test/internal/domain/service/auth"
	"github.com/huydq/test/internal/infrastructure/adapter/auth"
	config "github.com/huydq/test/internal/pkg/config"
	"github.com/huydq/test/internal/pkg/database"
)

type AuthUsecase struct {
	userRepo              userRepo.UserRepository
	twoFactorRepo         twoFactorRepo.TwoFactorTokenRepository
	jwtService            *auth.JWTService
	twoFactorTokenService authService.TwoFactorTokenService
	accessTokenService    authService.AccessTokenService
}

func NewAuthUsecase(
	userRepo userRepo.UserRepository,
	twoFactorRepo twoFactorRepo.TwoFactorTokenRepository,
	jwtService *auth.JWTService,
	twoFactorTokenService authService.TwoFactorTokenService,
	accessTokenService authService.AccessTokenService,
) *AuthUsecase {
	return &AuthUsecase{
		userRepo:              userRepo,
		twoFactorRepo:         twoFactorRepo,
		jwtService:            jwtService,
		twoFactorTokenService: twoFactorTokenService,
		accessTokenService:    accessTokenService,
	}
}

// Login handles the login request
func (uc *AuthUsecase) Login(ctx context.Context, input *inputdata.LoginInputData) (*outputdata.LoginOutputData, error) {
	user, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if user == nil || !user.VerifyPassword(input.Password) {
		return nil, errors.New("login.failed")
	}

	if user.EnabledMFA {
		tx, err := database.NewTx[*outputdata.GenerateTwoFAOutputData](ctx)
		if err != nil {
			return nil, err
		}

		verificationResp, err := tx.Transact(ctx, func(ctx context.Context) (*outputdata.GenerateTwoFAOutputData, error) {
			criteria := twoFactorModel.TwoFactorToken{
				UserID:  user.ID,
				MFAType: user.MFAType,
			}
			if err := uc.twoFactorRepo.InvalidatePreviousTokens(ctx, criteria); err != nil {
				return nil, err
			}

			_, err = uc.twoFactorTokenService.Create2FAToken(ctx, user.ID, user.MFAType, user.Email, user.FullName)
			if err != nil {
				return nil, err
			}

			return &outputdata.GenerateTwoFAOutputData{
				MFAType:   user.MFAType,
				ExpiresIn: config.GetConfig().MFATokenExpiryMinutes * 60,
			}, nil
		})

		if err != nil {
			return nil, errors.New("mfa.generate_failed")
		}

		return &outputdata.LoginOutputData{
			User:        user,
			RequiresMFA: true,
			MFAInfo: &outputdata.MFAInfo{
				Type:      object.MFAType(user.MFAType).String(),
				ExpiresIn: int(verificationResp.ExpiresIn),
			},
		}, nil
	}

	token, err := uc.jwtService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	if err := uc.accessTokenService.AddToken(ctx, token); err != nil {
		return nil, err
	}

	return &outputdata.LoginOutputData{
		Token:       token,
		User:        user,
		RequiresMFA: false,
	}, nil
}

func (uc *AuthUsecase) Logout(ctx context.Context, token string) error {
	verifyToken, err := uc.accessTokenService.GetVerifyToken(ctx, token)
	if err != nil {
		return errors.New("トークンブラックリストを作成できません")

	}

	if verifyToken != nil {
		if verifyToken.IsActive {
			err := uc.accessTokenService.UpdateToken(ctx, verifyToken)
			if err != nil {
				return errors.New("トークンブラックリストを作成できません")

			}
		}
	}

	return err
}

// GetMe retrieves a user by their ID
func (uc *AuthUsecase) GetMe(ctx context.Context, userID int) (*outputdata.UserProfileOutputData, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user.not_found")
	}

	return &outputdata.UserProfileOutputData{
		User: user,
	}, nil
}

// Verify2FAToken verifies a 2FA token
func (uc *AuthUsecase) Verify2FAToken(ctx context.Context, input *inputdata.VerifyTwoFAInputData) (*outputdata.VerifyTwoFAOutputData, error) {
	user, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user.not_found")
	}

	token, err := uc.twoFactorRepo.FindByToken(ctx, twoFactorModel.TwoFactorToken{
		UserID:    user.ID,
		Token:     input.Token,
		IsUsed:    false,
		ExpiredAt: time.Now(),
	})
	if err != nil {
		return nil, err
	}

	if token == nil {
		return nil, errors.New("mfa.invalid_token")
	}

	if token.ExpiredAt.Before(time.Now()) {
		return nil, errors.New("mfa.expired_token")
	}

	// update token status to used
	token.MarkAsUsed()
	if err := uc.twoFactorRepo.MarkAsUsed(ctx, token); err != nil {
		return nil, err
	}

	jwtToken, err := uc.jwtService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	if err := uc.accessTokenService.AddToken(ctx, jwtToken); err != nil {
		return nil, err
	}

	return &outputdata.VerifyTwoFAOutputData{
		Token: jwtToken,
		User:  user,
	}, nil
}

// ResendCode handles resending a 2FA code to the user
func (uc *AuthUsecase) ResendCode(ctx context.Context, input *inputdata.ResendCodeInputData) (*outputdata.ResendCodeOutputData, error) {
	user, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user.not_found")
	}

	canResend, remainingTime, err := uc.twoFactorTokenService.CanResendToken(ctx, user.ID, input.MfaType)
	if err != nil {
		return nil, err
	}

	canResendOutput := &outputdata.CanResendCodeOutputData{
		CanResend:     canResend,
		RemainingTime: remainingTime,
	}

	if !canResendOutput.CanResend {
		return &outputdata.ResendCodeOutputData{
			CanResend: false,
			ExpiresIn: canResendOutput.RemainingTime,
		}, nil
	}

	tx, err := database.NewTx[*outputdata.GenerateTwoFAOutputData](ctx)
	if err != nil {
		return nil, err
	}

	generateOutput, err := tx.Transact(ctx, func(ctx context.Context) (*outputdata.GenerateTwoFAOutputData, error) {
		// Create invalidation criteria
		criteria := twoFactorModel.TwoFactorToken{
			UserID:  user.ID,
			MFAType: input.MfaType,
		}

		if err := uc.twoFactorRepo.InvalidatePreviousTokens(ctx, criteria); err != nil {
			return nil, err
		}

		_, err = uc.twoFactorTokenService.Create2FAToken(ctx, user.ID, input.MfaType, user.Email, user.FullName)
		if err != nil {
			return nil, err
		}

		return &outputdata.GenerateTwoFAOutputData{
			MFAType:   input.MfaType,
			ExpiresIn: config.GetConfig().MFATokenExpiryMinutes * 60,
		}, nil
	})

	if err != nil {
		return nil, errors.New("mfa.generate_failed")
	}

	return &outputdata.ResendCodeOutputData{
		CanResend: true,
		ExpiresIn: generateOutput.ExpiresIn,
	}, nil
}
