package repositories

import (
	"context"

	models "github.com/huydq/test/src/domain/models"
)

// RoleRepository defines the interface for role data access
type RoleRepository interface {
	// FindByID finds a role by ID
	FindByID(ctx context.Context, id int) (*models.Role, error)

	// FindByCode finds a role by code
	FindByCode(ctx context.Context, code string) (*models.Role, error)
}
