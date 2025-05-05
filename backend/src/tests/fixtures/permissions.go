package fixtures

import (
	"context"

	"github.com/vnlab/makeshop-payment/src/domain/models"
	"github.com/vnlab/makeshop-payment/src/domain/repositories"
)

// GetMockPermission creates a mock permission with the specified attributes
func GetMockPermission(id int, name, code string) *models.Permission {
	return &models.Permission{
		ID:   id,
		Name: name,
		Code: code,
	}
}

// GetSeedPermissions returns all permissions from the database
func GetSeedPermissions(ctx context.Context, repo repositories.PermissionRepository) []*models.Permission {
	permissions, err := repo.List(ctx)
	if err != nil || len(permissions) == 0 {
		return []*models.Permission{} // Return empty slice on error
	}
	return permissions
}

// GetAdminPermission returns a permission that represents admin capabilities
// In the new seed data, there's no single "ADMIN" permission, but USER_MANAGE is a key admin permission
func GetAdminPermission(ctx context.Context, repo repositories.PermissionRepository) *models.Permission {
	// USER_MANAGE is one of the core admin permissions in the new seed data
	adminPermission, err := repo.FindByIDs(ctx, []int{1}) // USER_MANAGE has ID 1
	if err != nil || len(adminPermission) == 0 {
		// If the permission isn't found, return a mock for USER_MANAGE
		return GetMockPermission(1, "ユーザー管理", "USER_MANAGE")
	}
	return adminPermission[0]
}

// FindPermissionByID finds a permission by ID
func FindPermissionByID(ctx context.Context, repo repositories.PermissionRepository, id int) *models.Permission {
	permissions, err := repo.FindByIDs(ctx, []int{id})
	if err != nil || len(permissions) == 0 {
		return nil
	}

	return permissions[0]
}

// FindPermissionsByIDs finds permissions by IDs
func FindPermissionsByIDs(ctx context.Context, repo repositories.PermissionRepository, ids []int) []*models.Permission {
	permissions, err := repo.FindByIDs(ctx, ids)
	if err != nil {
		return []*models.Permission{}
	}

	return permissions
}
