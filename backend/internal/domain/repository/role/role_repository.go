package repository

import (
	"context"

	model "github.com/huydq/test/internal/domain/model/role"
)

// RoleRepository defines the interface for role data access
type RoleRepository interface {
	// FindByID finds a role by ID
	FindByID(ctx context.Context, id int) (*model.Role, error)

	// FindByCode finds a role by code
	FindByCode(ctx context.Context, code string) (*model.Role, error)

	// FindByName finds a role by name
	FindByName(ctx context.Context, name string) (*model.Role, error)

	// Create creates a new role
	Create(ctx context.Context, role *model.Role) error

	// Update updates an existing role
	Update(ctx context.Context, role *model.Role) error

	// Delete soft-deletes a role by ID
	Delete(ctx context.Context, id int) error

	// List lists all roles with pagination
	List(ctx context.Context) ([]*model.Role, error)
}
