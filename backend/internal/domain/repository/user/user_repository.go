package user

import (
	"context"

	"github.com/huydq/test/internal/datastructure/inputdata"
	"github.com/huydq/test/internal/domain/model/user"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// FindByID finds a user by ID
	FindByID(ctx context.Context, id int) (*user.User, error)

	// FindByEmail finds a user by email
	FindByEmail(ctx context.Context, email string) (*user.User, error)

	// Create creates a new user
	Create(ctx context.Context, user *user.User) error

	// Update updates an existing user
	Update(ctx context.Context, user *user.User) error

	// Delete soft-deletes a user by ID
	Delete(ctx context.Context, id int) error

	// List lists users with filtering and pagination
	List(ctx context.Context, params *inputdata.UserListInputData) ([]*user.User, int, int, error)

	// GetUsersWithAuditLogs retrieves users who have audit log entries
	GetUsersWithAuditLogs(ctx context.Context) ([]*user.User, error)
}
