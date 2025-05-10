package converter

import (
	"github.com/huydq/test/internal/datastructure/inputdata"
	"github.com/huydq/test/internal/datastructure/outputdata"
	modelRole "github.com/huydq/test/internal/domain/model/role"
	objectRole "github.com/huydq/test/internal/domain/object/role"
)

// RoleCreateInputToModel converts a CreateRoleInput to a Role domain model
func RoleCreateInputToModel(input *inputdata.CreateRoleInput) *modelRole.Role {
	if input == nil {
		return nil
	}

	return &modelRole.Role{
		Name: input.Name,
		Code: objectRole.RoleCode(input.Code),
	}
}

// RoleUpdateInputToModel applies UpdateRoleInput changes to an existing Role model
func RoleUpdateInputToModel(role *modelRole.Role, input *inputdata.UpdateRoleInput) {
	if role == nil || input == nil {
		return
	}

	role.Name = input.Name
	role.Code = objectRole.RoleCode(input.Code)
	// Permissions will be handled separately
}

// RoleModelToOutput converts a Role domain model to a RoleOutput
func RoleModelToOutput(role *modelRole.Role) *outputdata.RoleOutput {
	if role == nil {
		return nil
	}

	var permissionsOutput []*outputdata.PermissionOutput
	if role.Permissions != nil {
		permissionsOutput = PermissionModelsToOutputs(role.Permissions)
	}

	return &outputdata.RoleOutput{
		ID:          role.ID,
		Name:        role.Name,
		Code:        string(role.Code),
		Permissions: permissionsOutput,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}

// RoleModelsToOutputs converts a slice of Role domain models to a slice of RoleOutputs
func RoleModelsToOutputs(roles []*modelRole.Role) []*outputdata.RoleOutput {
	if roles == nil {
		return nil
	}

	outputs := make([]*outputdata.RoleOutput, len(roles))
	for i, role := range roles {
		outputs[i] = RoleModelToOutput(role)
	}
	return outputs
}

// BatchRolePermissionUpdateInputToModel converts BatchUpdateRolePermissionsInput to domain structure
func BatchRolePermissionUpdateInputToModel(input *inputdata.BatchUpdateRolePermissionsInput) []struct {
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

// BatchUpdateResultToOutput converts successful update IDs to BatchUpdateRolePermissionsOutput
func BatchUpdateResultToOutput(successfulIDs []int) *outputdata.BatchUpdateRolePermissionsOutput {
	return &outputdata.BatchUpdateRolePermissionsOutput{
		SuccessfulUpdates: successfulIDs,
	}
}
