package fixtures

import (
	"github.com/huydq/test/src/domain/models"
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
