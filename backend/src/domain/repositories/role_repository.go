package repositories

import (
	"context"

	models "github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/domain/repositories/filter"
)

// RoleRepository defines the interface for role data access
type RoleRepository interface {
	// FindByID finds a role by ID
	FindByID(ctx context.Context, id int) (*models.Role, error)

	// FindByCode finds a role by code
	FindByCode(ctx context.Context, code string) (*models.Role, error)

	// FindByName finds a role by name
	FindByName(ctx context.Context, name string) (*models.Role, error)

	// Create creates a new role
	Create(ctx context.Context, role *models.Role) error

	// Update updates an existing role
	Update(ctx context.Context, role *models.Role) error

	// Delete soft-deletes a role by ID
	Delete(ctx context.Context, id int) error

	// List lists all roles with pagination
	List(ctx context.Context, filter *filter.RoleFilter) ([]*models.Role, int, int64, error)
}
