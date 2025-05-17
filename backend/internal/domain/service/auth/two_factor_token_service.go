package auth

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	twoFactorModel "github.com/huydq/test/internal/domain/model/two_factor_token"
	twoFactorRepo "github.com/huydq/test/internal/domain/repository/two_factor_token"
	userRepo "github.com/huydq/test/internal/domain/repository/user"
	"github.com/huydq/test/internal/infrastructure/adapter/email"
	"github.com/huydq/test/internal/pkg/config"
)

type TwoFactorTokenService interface {
	Create2FAToken(ctx context.Context, userID int, mfaType int, userEmail, fullName string) (*twoFactorModel.TwoFactorToken, error)
	CanResendToken(ctx context.Context, userID int, mfaType int) (bool, int, error)
}

// TwoFactorTokenServiceImpl implements the TwoFactorTokenService interface
type TwoFactorTokenServiceImpl struct {
	userRepo       userRepo.UserRepository
	twoFactorRepo  twoFactorRepo.TwoFactorTokenRepository
	mailService    *email.MailService
	tokenExpiryMin int
}

// NewTwoFactorTokenService creates a new TwoFactorTokenService implementation
func NewTwoFactorTokenService(
	userRepo userRepo.UserRepository,
	twoFactorRepo twoFactorRepo.TwoFactorTokenRepository,
	mailService *email.MailService,
) TwoFactorTokenService {
	return &TwoFactorTokenServiceImpl{
		userRepo:       userRepo,
		twoFactorRepo:  twoFactorRepo,
		mailService:    mailService,
		tokenExpiryMin: config.GetConfig().MFATokenExpiryMinutes,
	}
}

// Create2FAToken creates a new two-factor token
func (s *TwoFactorTokenServiceImpl) Create2FAToken(ctx context.Context, userID int, mfaType int, userEmail, fullName string) (*twoFactorModel.TwoFactorToken, error) {
	token := fmt.Sprintf("%06d", rand.Intn(1000000))
	expiredAt := time.Now().Add(time.Duration(s.tokenExpiryMin) * time.Minute)

	twoFactorToken := &twoFactorModel.TwoFactorToken{
		UserID:    userID,
		Token:     token,
		MFAType:   mfaType,
		IsUsed:    false,
		ExpiredAt: expiredAt,
	}

	if err := s.twoFactorRepo.Create(ctx, twoFactorToken); err != nil {
		return nil, err
	}

	s.mailService.SendMailByTemplateID("send_2fa_code", email.TwoFACodeEmailData{
		Email:          userEmail,
		ToName:         fullName,
		Token:          token,
		TokenExpiryMin: s.tokenExpiryMin,
	})

	return twoFactorToken, nil
}

// CanResendToken checks if a user can resend a 2FA token
func (s *TwoFactorTokenServiceImpl) CanResendToken(ctx context.Context, userID int, mfaType int) (bool, int, error) {
	resendInterval := config.GetConfig().MFATokenResendInterval

	criteria := twoFactorModel.TwoFactorToken{
		UserID:  userID,
		MFAType: mfaType,
	}

	lastToken, err := s.twoFactorRepo.GetLastToken(ctx, criteria)
	if err != nil {
		return false, 0, err
	}

	if lastToken == nil {
		return true, 0, nil
	}

	earliestNextResendTime := time.Now().Add(-time.Duration(resendInterval) * time.Minute)
	remainingTime := time.Until(lastToken.CreatedAt.Add(time.Duration(resendInterval) * time.Minute))

	return lastToken.CreatedAt.Before(earliestNextResendTime), int(remainingTime.Seconds()), nil
}
