package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	tokenModel "github.com/huydq/test/internal/domain/model/token"
	tokenRepo "github.com/huydq/test/internal/domain/repository/token"
	config "github.com/huydq/test/internal/pkg/config"
)

type AccessTokenService interface {
	AddToken(ctx context.Context, tokenString string) error
	GetVerifyToken(ctx context.Context, tokenString string) (*tokenModel.Token, error)
	UpdateToken(ctx context.Context, token *tokenModel.Token) error
	IsBlacklisted(ctx context.Context, tokenString string) (bool, error)
}

type accessTokenService struct {
	tokenRepo tokenRepo.TokenRepository
}

func NewAccessTokenService(
	tokenRepo tokenRepo.TokenRepository,
) AccessTokenService {
	return &accessTokenService{
		tokenRepo: tokenRepo,
	}
}

// AddToken creates a new access token in the database
func (s *accessTokenService) AddToken(ctx context.Context, tokenString string) error {
	tokenExpiredMinutes := config.GetConfig().MFATokenExpiryMinutes
	expiredAt := time.Now().Add(time.Duration(tokenExpiredMinutes) * time.Minute)
	hash := sha256.Sum256([]byte(tokenString))
	hashedToken := hex.EncodeToString(hash[:])

	tokenModel := &tokenModel.Token{
		Token:     hashedToken,
		IsActive:  true,
		ExpiredAt: expiredAt,
	}

	if err := s.tokenRepo.Create(ctx, tokenModel); err != nil {
		return err
	}
	return nil
}

// GetVerifyToken verifies a token and returns its details
func (s *accessTokenService) GetVerifyToken(ctx context.Context, tokenString string) (*tokenModel.Token, error) {
	hash := sha256.Sum256([]byte(tokenString))
	hashedToken := hex.EncodeToString(hash[:])
	token, err := s.tokenRepo.FindByToken(ctx, hashedToken)

	if token == nil {
		return nil, errors.New("トークンが見つかりません")
	}

	return token, err
}

// UpdateToken updates a token's status in the database
func (s *accessTokenService) UpdateToken(ctx context.Context, token *tokenModel.Token) error {
	token.Invalidate()
	return s.tokenRepo.Update(ctx, token)
}

// IsBlacklisted checks if a token is blacklisted
func (s *accessTokenService) IsBlacklisted(ctx context.Context, tokenString string) (bool, error) {
	verifyToken, err := s.GetVerifyToken(ctx, tokenString)
	if err != nil {
		return false, err
	}

	return !verifyToken.IsActive, nil
}
