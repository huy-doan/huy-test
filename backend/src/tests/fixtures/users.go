package fixtures

import (
	"context"
	"fmt"
	"time"

	"github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/domain/repositories"
	"golang.org/x/crypto/bcrypt"
)

// GetMockUser returns a mock user with the specified email and password
func GetMockUser(email, password string, role *models.Role) *models.User {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return &models.User{
		ID:           1,
		Email:        email,
		PasswordHash: string(hashedPassword),
		RoleID:       role.ID,
		Role:         role,
	}
}

// CreateUniqueTestUser creates a unique user for testing with a timestamp suffix
// to ensure unique email addresses and lets the database generate the ID
func CreateUniqueTestUser(ctx context.Context, repo repositories.UserRepository, emailPrefix string, password string, role *models.Role) (*models.User, error) {
	// Create a unique email using current timestamp to avoid conflicts
	timestamp := time.Now().UnixNano()
	email := fmt.Sprintf("%s_%d@test.example.com", emailPrefix, timestamp)

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create the user without specifying ID (let the database generate it)
	user := &models.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		RoleID:       role.ID,
		Role:         role,
	}

	// Create the user in the database
	err = repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
