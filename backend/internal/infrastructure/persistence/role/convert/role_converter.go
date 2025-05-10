package convert

import (
	modelPermission "github.com/huydq/test/internal/domain/model/permission"
	modelRole "github.com/huydq/test/internal/domain/model/role"
	util "github.com/huydq/test/internal/domain/object/basedatetime"
	objectRole "github.com/huydq/test/internal/domain/object/role"
	permissionConvert "github.com/huydq/test/internal/infrastructure/persistence/permission/convert"
	permissionDto "github.com/huydq/test/internal/infrastructure/persistence/permission/dto"
	"github.com/huydq/test/internal/infrastructure/persistence/role/dto"
	persistence "github.com/huydq/test/internal/infrastructure/persistence/util"
	"gorm.io/gorm"
)

// ToRoleDTO converts a Role domain model to a RoleDTO
func ToRoleDTO(role *modelRole.Role) *dto.RoleDTO {
	if role == nil {
		return nil
	}

	var permissionDTOs []*permissionDto.PermissionDTO
	if role.Permissions != nil {
		permissionDTOs = permissionConvert.ToPermissionDTOs(role.Permissions)
	}

	result := &dto.RoleDTO{
		ID:   role.ID,
		Name: role.Name,
		Code: string(role.Code),
		BaseColumnTimestamp: persistence.BaseColumnTimestamp{
			CreatedAt: role.CreatedAt,
			UpdatedAt: role.UpdatedAt,
		},
		Permissions: permissionDTOs,
	}

	// Handle the conversion from *time.Time to gorm.DeletedAt
	if role.DeletedAt != nil {
		result.DeletedAt = gorm.DeletedAt{
			Time:  *role.DeletedAt,
			Valid: true,
		}
	} else {
		result.DeletedAt = gorm.DeletedAt{
			Valid: false,
		}
	}

	return result
}

// ToRoleModel converts a RoleDTO to a Role domain model
func ToRoleModel(dtoObj *dto.RoleDTO) *modelRole.Role {
	if dtoObj == nil {
		return nil
	}

	var permissions []*modelPermission.Permission
	if dtoObj.Permissions != nil {
		permissions = permissionConvert.ToPermissionModels(dtoObj.Permissions)
	}

	result := &modelRole.Role{
		ID:   dtoObj.ID,
		Name: dtoObj.Name,
		Code: objectRole.RoleCode(dtoObj.Code),
		BaseColumnTimestamp: util.BaseColumnTimestamp{
			CreatedAt: dtoObj.CreatedAt,
			UpdatedAt: dtoObj.UpdatedAt,
		},
		Permissions: permissions,
	}

	// Handle the conversion from gorm.DeletedAt to *time.Time
	if dtoObj.DeletedAt.Valid {
		deletedAt := dtoObj.DeletedAt.Time
		result.DeletedAt = &deletedAt
	} else {
		result.DeletedAt = nil
	}

	return result
}

// ToRoleDTOs converts a list of Role domain models to a list of RoleDTOs
func ToRoleDTOs(roles []*modelRole.Role) []*dto.RoleDTO {
	if roles == nil {
		return nil
	}

	result := make([]*dto.RoleDTO, len(roles))
	for i, role := range roles {
		result[i] = ToRoleDTO(role)
	}
	return result
}

// ToRoleModels converts a list of RoleDTOs to a list of Role domain models
func ToRoleModels(dtos []*dto.RoleDTO) []*modelRole.Role {
	if dtos == nil {
		return nil
	}

	result := make([]*modelRole.Role, len(dtos))
	for i, dto := range dtos {
		result[i] = ToRoleModel(dto)
	}
	return result
}
