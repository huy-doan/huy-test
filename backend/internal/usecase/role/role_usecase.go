package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/huydq/test/internal/datastructure/inputdata"
	"github.com/huydq/test/internal/datastructure/outputdata"
	model "github.com/huydq/test/internal/domain/model/role"
	"github.com/huydq/test/internal/domain/service"

	"gorm.io/gorm"
)

type RoleUsecase interface {
	CreateRole(ctx context.Context, role *model.Role) (*model.Role, error)
	UpdateRole(ctx context.Context, id int, input *inputdata.UpdateRoleInput) (*model.Role, error)
	DeleteRole(ctx context.Context, id int) error
	GetRoleByID(ctx context.Context, id int) (*model.Role, error)
	ListRoles(ctx context.Context) (*outputdata.RoleListOutput, error)
	BatchUpdateRolePermissions(ctx context.Context, input *inputdata.BatchUpdateRolePermissionsInput) (*outputdata.BatchUpdateRolePermissionsOutput, error)
}

type roleUsecaseImpl struct {
	roleService service.RoleService
}

func NewRoleUsecase(
	roleService service.RoleService,
) RoleUsecase {
	return &roleUsecaseImpl{
		roleService: roleService,
	}
}

func (u *roleUsecaseImpl) CreateRole(ctx context.Context, role *model.Role) (*model.Role, error) {
	existingByName, err := u.roleService.GetRoleByName(ctx, role.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingByName != nil {
		return nil, fmt.Errorf("role with name '%s' already exists", role.Name)
	}

	if err := u.roleService.CreateRole(ctx, role); err != nil {
		return nil, err
	}

	createdRole, err := u.roleService.GetRoleByID(ctx, role.ID)
	if err != nil {
		return nil, err
	}

	return createdRole, nil
}

func (u *roleUsecaseImpl) UpdateRole(ctx context.Context, id int, input *inputdata.UpdateRoleInput) (*model.Role, error) {
	role, err := u.roleService.GetRoleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, fmt.Errorf("role with ID %d not found", id)
	}

	if role.Name != input.Name {
		existingByName, err := u.roleService.GetRoleByName(ctx, input.Name)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existingByName != nil && existingByName.ID != id {
			return nil, fmt.Errorf("role with name '%s' already exists", input.Name)
		}
	}

	if input == nil {
		return nil, fmt.Errorf("role or input cannot be nil")
	}

	role.Name = input.Name

	if input.PermissionIDs != nil {
		permissions, err := u.roleService.GetPermissionsByIDs(ctx, input.PermissionIDs)
		if err != nil {
			return nil, err
		}

		if len(permissions) != len(input.PermissionIDs) {
			return nil, fmt.Errorf("one or more permission IDs do not exist")
		}

		role.Permissions = permissions
	}

	if err := u.roleService.UpdateRole(ctx, role); err != nil {
		return nil, err
	}

	updatedRole, err := u.roleService.GetRoleByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return updatedRole, nil
}

func (u *roleUsecaseImpl) DeleteRole(ctx context.Context, id int) error {
	role, err := u.roleService.GetRoleByID(ctx, id)
	if err != nil {
		return err
	}
	if role == nil {
		return fmt.Errorf("role with ID %d not found", id)
	}

	return u.roleService.DeleteRole(ctx, id)
}

func (u *roleUsecaseImpl) GetRoleByID(ctx context.Context, id int) (*model.Role, error) {
	role, err := u.roleService.GetRoleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, fmt.Errorf("role with ID %d not found", id)
	}

	return role, nil
}

func (u *roleUsecaseImpl) ListRoles(ctx context.Context) (*outputdata.RoleListOutput, error) {
	roles, err := u.roleService.ListRoles(ctx)
	if err != nil {
		return nil, err
	}

	return &outputdata.RoleListOutput{
		Roles: roles,
		Total: int64(len(roles)),
	}, nil
}

func (u *roleUsecaseImpl) BatchUpdateRolePermissions(ctx context.Context, input *inputdata.BatchUpdateRolePermissionsInput) (*outputdata.BatchUpdateRolePermissionsOutput, error) {
	updates := batchRolePermissionUpdateInputToModel(input)

	allPermissionIDs := make(map[int]struct{})
	for _, update := range updates {
		for _, permID := range update.PermissionIDs {
			allPermissionIDs[permID] = struct{}{}
		}
	}

	permIDsToValidate := make([]int, 0, len(allPermissionIDs))
	for permID := range allPermissionIDs {
		permIDsToValidate = append(permIDsToValidate, permID)
	}

	if len(permIDsToValidate) > 0 {
		permissions, err := u.roleService.GetPermissionsByIDs(ctx, permIDsToValidate)
		if err != nil {
			return nil, err
		}

		if len(permissions) != len(permIDsToValidate) {
			return nil, fmt.Errorf("one or more permission IDs do not exist")
		}
	}

	successfulUpdates := make([]int, 0, len(updates))

	for _, update := range updates {
		err := u.roleService.UpdateRolePermissions(ctx, update.ID, update.PermissionIDs)
		if err != nil {
			continue
		}

		successfulUpdates = append(successfulUpdates, update.ID)
	}

	return batchUpdateResultToOutput(successfulUpdates), nil
}

func batchUpdateResultToOutput(successfulIDs []int) *outputdata.BatchUpdateRolePermissionsOutput {
	return &outputdata.BatchUpdateRolePermissionsOutput{
		SuccessfulUpdates: successfulIDs,
	}
}

func batchRolePermissionUpdateInputToModel(input *inputdata.BatchUpdateRolePermissionsInput) []struct {
	ID            int
	PermissionIDs []int
} {
	if input == nil {
		return nil
	}

	updates := make([]struct {
		ID            int
		PermissionIDs []int
	}, len(input.Updates))

	for i, update := range input.Updates {
		updates[i] = struct {
			ID            int
			PermissionIDs []int
		}{
			ID:            update.ID,
			PermissionIDs: update.PermissionIDs,
		}
	}

	return updates
}
