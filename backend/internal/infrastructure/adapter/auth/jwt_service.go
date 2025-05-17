package auth

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/huydq/test/internal/domain/model/user"
	"github.com/huydq/test/internal/pkg/config"
)

// JWTService provides JWT token generation and validation
type JWTService struct {
	secretKey     string
	tokenDuration time.Duration
	blacklist     map[string]time.Time // Map to store blacklisted tokens
	mutex         sync.RWMutex         // Mutex to ensure thread-safety
}

// TokenClaims represents the claims in a JWT token
type TokenClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	RoleID int    `json:"role_id"`
	jwt.RegisteredClaims
}

// NewJWTService creates a new JWTService
func NewJWTService() *JWTService {
	// load secret key from config file
	appConfig := config.GetConfig()
	secret := appConfig.JWTSecret
	if secret == "" {
		secret = "default_jwt_secret_key_change_in_production"
	}

	// Get token duration from environment (default 24 hours)
	var tokenDuration time.Duration
	hours := appConfig.JWTDurationHour
	tokenDuration = time.Duration(hours) * time.Hour

	return &JWTService{
		secretKey:     secret,
		tokenDuration: tokenDuration,
		blacklist:     make(map[string]time.Time),
	}
}

// GenerateToken generates a new JWT token for a user
func (s *JWTService) GenerateToken(user *user.User) (string, error) {
	if user == nil {
		return "", errors.New("user is nil")
	}

	claims := TokenClaims{
		UserID: user.ID,
		Email:  user.Email,
		RoleID: user.RoleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

// ValidateToken validates the provided token string and returns the claims
func (s *JWTService) ValidateToken(tokenString string) (*TokenClaims, error) {
	// Check if the token is blacklisted
	s.mutex.RLock()
	_, blacklisted := s.blacklist[tokenString]
	s.mutex.RUnlock()

	if blacklisted {
		return nil, errors.New("token has been revoked")
	}

	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (any, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// BlacklistToken adds a token to the blacklist
func (s *JWTService) BlacklistToken(tokenString string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.blacklist[tokenString] = time.Now()

	// Auto cleanup expired tokens from the blacklist
	s.cleanupBlacklist()
}

// IsBlacklisted checks if a token is in the blacklist
func (s *JWTService) IsBlacklisted(tokenString string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	_, exists := s.blacklist[tokenString]
	return exists
}

// ExtractUserIDFromToken extracts user ID from a token string
func (s *JWTService) ExtractUserIDFromToken(tokenString string) (int, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// cleanupBlacklist removes expired tokens from the blacklist
func (s *JWTService) cleanupBlacklist() {
	// Implement cleanupBlacklist in a goroutine to avoid blocking the main process
	go func() {
		// Ensure thread-safety when deleting
		s.mutex.Lock()
		defer s.mutex.Unlock()

		// Check and delete tokens that have expired (e.g., over 24 hours)
		threshold := time.Now().Add(-24 * time.Hour)
		for token, blacklistedAt := range s.blacklist {
			if blacklistedAt.Before(threshold) {
				delete(s.blacklist, token)
			}
		}
	}()
}
