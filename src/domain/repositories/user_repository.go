package repositories

import (
	"context"

	models "github.com/huydq/test/src/domain/models"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// FindByID finds a user by ID
	FindByID(ctx context.Context, id int) (*models.User, error)

	// FindByEmail finds a user by email
	FindByEmail(ctx context.Context, email string) (*models.User, error)

	// Create creates a new user
	Create(ctx context.Context, user *models.User) error

	// Update updates an existing user
	Update(ctx context.Context, user *models.User) error

	// Delete soft-deletes a user by ID
	Delete(ctx context.Context, id int) error

	// List lists all users with pagination
	List(ctx context.Context, page, pageSize int) ([]*models.User, int, error)
}
