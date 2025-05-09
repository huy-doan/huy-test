package repositories

import (
	"context"

	"github.com/huydq/test/src/domain/models"
)

type PermissionRepository interface {
	// FindByIDs finds multiple permissions by their IDs
	FindByIDs(ctx context.Context, ids []int) ([]*models.Permission, error)

	// List retrieves all permissions
	List(ctx context.Context) ([]*models.Permission, error)
}
