package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/vnlab/makeshop-payment/src/domain/models"
	"github.com/vnlab/makeshop-payment/src/domain/repositories"
	"github.com/vnlab/makeshop-payment/src/domain/repositories/filter"
)

type RoleUsecase struct {
	roleRepo       repositories.RoleRepository
	permissionRepo repositories.PermissionRepository
}

// NewRoleUsecase creates a new role usecase with necessary repositories
func NewRoleUsecase(
	roleRepo repositories.RoleRepository,
	permissionRepo repositories.PermissionRepository,
) *RoleUsecase {
	return &RoleUsecase{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

type UpdateRoleRequest struct {
	Name          string `json:"name" binding:"required"`
	PermissionIDs []int  `json:"permission_ids"`
}

func (r *RoleUsecase) ListRoles(ctx context.Context, filter *filter.RoleFilter) ([]*models.Role, int, int64, error) {
	return r.roleRepo.List(ctx, filter)
}

func (r *RoleUsecase) GetRoleByID(ctx context.Context, id int) (*models.Role, error) {
	return r.roleRepo.FindByID(ctx, id)
}

// ValidatePermissions checks if all permission IDs exist in the database
func (r *RoleUsecase) ValidatePermissions(ctx context.Context, permissionIDs []int) error {
	if len(permissionIDs) == 0 {
		return nil
	}

	// Get all valid permissions from the database
	validPermissions, err := r.permissionRepo.FindByIDs(ctx, permissionIDs)
	if err != nil {
		return err
	}

	// If the count of found permissions doesn't match the requested permissions count,
	// some permissions don't exist
	if len(validPermissions) != len(permissionIDs) {
		return fmt.Errorf("one or more permission IDs do not exist")
	}

	return nil
}

func (r *RoleUsecase) CreateRole(ctx context.Context, role *models.Role) error {
	// Check if role name or code already exists
	existingByName, err := r.roleRepo.FindByName(ctx, role.Name)
	if err != nil {
		return err
	}
	if existingByName != nil {
		return fmt.Errorf("role name already exists")
	}

	existingByCode, err := r.roleRepo.FindByCode(ctx, role.Code)
	if err != nil {
		return err
	}
	if existingByCode != nil {
		return fmt.Errorf("role code already exists")
	}

	// Validate permissions if they exist
	if len(role.Permissions) > 0 {
		permissionIDs := make([]int, 0, len(role.Permissions))
		for _, perm := range role.Permissions {
			permissionIDs = append(permissionIDs, perm.ID)
		}

		if err := r.ValidatePermissions(ctx, permissionIDs); err != nil {
			return err
		}
	}

	return r.roleRepo.Create(ctx, role)
}

func (r *RoleUsecase) UpdateRole(ctx context.Context, id int, updateData *UpdateRoleRequest) error {
	role, err := r.roleRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}

	// Check if updated name already exists for another role
	if role.Name != updateData.Name {
		existingByName, err := r.roleRepo.FindByName(ctx, updateData.Name)
		if err != nil {
			return err
		}
		if existingByName != nil && existingByName.ID != id {
			return fmt.Errorf("role name already exists")
		}
	}

	// Validate permissions if they exist
	if len(updateData.PermissionIDs) > 0 {
		if err := r.ValidatePermissions(ctx, updateData.PermissionIDs); err != nil {
			return err
		}
	}

	// Update role fields
	role.Name = updateData.Name

	// Update permissions (would need a method to sync permissions with the role)
	permissions, err := r.permissionRepo.FindByIDs(ctx, updateData.PermissionIDs)
	if err != nil {
		return err
	}
	role.Permissions = permissions

	err = r.roleRepo.Update(ctx, role)
	if err != nil {
		return err
	}
	return nil
}

func (r *RoleUsecase) DeleteRole(ctx context.Context, id int) error {
	return r.roleRepo.Delete(ctx, id)
}

// BatchUpdateRolePermissions updates permissions for multiple roles at once
func (r *RoleUsecase) BatchUpdateRolePermissions(ctx context.Context, updates []struct {
	ID            int
	PermissionIDs []int
}) ([]int, error) {
	// Collect all permission IDs for validation
	allPermissionIDs := make(map[int]struct{})
	for _, update := range updates {
		for _, permID := range update.PermissionIDs {
			allPermissionIDs[permID] = struct{}{}
		}
	}

	// Convert to slice for validation
	permIDsToValidate := make([]int, 0, len(allPermissionIDs))
	for permID := range allPermissionIDs {
		permIDsToValidate = append(permIDsToValidate, permID)
	}

	// Validate all permissions in one go
	if err := r.ValidatePermissions(ctx, permIDsToValidate); err != nil {
		return nil, err
	}

	// Process each role update
	successfulUpdates := make([]int, 0, len(updates))

	for _, update := range updates {
		// Get the role
		role, err := r.roleRepo.FindByID(ctx, update.ID)
		if err != nil {
			continue // Skip this role if there's an error
		}
		if role == nil {
			continue // Skip non-existent roles
		}

		// Get permissions for this role
		permissions, err := r.permissionRepo.FindByIDs(ctx, update.PermissionIDs)
		if err != nil {
			continue // Skip this role if there's an error getting permissions
		}

		// Update the role's permissions
		role.Permissions = permissions

		// Save the changes
		err = r.roleRepo.Update(ctx, role)
		if err != nil {
			continue // Skip this role if update fails
		}

		// Record successful update
		successfulUpdates = append(successfulUpdates, role.ID)
	}

	return successfulUpdates, nil
}
