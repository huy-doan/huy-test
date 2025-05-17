package mapper

import (
	"github.com/huydq/test/internal/datastructure/inputdata"
	permissionModel "github.com/huydq/test/internal/domain/model/permission"
	model "github.com/huydq/test/internal/domain/model/role"
	generated "github.com/huydq/test/internal/pkg/api/generated"
)

func ToCreateRoleModel(request generated.CreateRoleRequest) (*model.Role, error) {
	permissions := make([]*permissionModel.Permission, len(request.PermissionIds))
	for i, id := range request.PermissionIds {
		permissions[i] = &permissionModel.Permission{ID: id}
	}
	return &model.Role{
		Name:        request.Name,
		Permissions: permissions,
	}, nil
}

func ToUpdateRoleInput(request generated.UpdateRoleRequest) (*inputdata.UpdateRoleInput, error) {
	return &inputdata.UpdateRoleInput{
		Name:          request.Name,
		PermissionIDs: request.PermissionIds,
	}, nil
}

func ToBatchUpdateRolePermissionsInput(request model.BatchUpdateRolePermissionsRequest) (*inputdata.BatchUpdateRolePermissionsInput, error) {
	updates := make([]inputdata.RolePermissionUpdate, len(request))
	for i, item := range request {
		updates[i] = inputdata.RolePermissionUpdate{
			ID:            item.ID,
			PermissionIDs: item.PermissionIDs,
		}
	}
	return &inputdata.BatchUpdateRolePermissionsInput{
		Updates: updates,
	}, nil
}
