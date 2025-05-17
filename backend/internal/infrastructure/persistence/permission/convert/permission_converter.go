package convert

import (
	modelPermission "github.com/huydq/test/internal/domain/model/permission"
	modelScreen "github.com/huydq/test/internal/domain/model/screen"
	util "github.com/huydq/test/internal/domain/object/basedatetime"
	objectPermission "github.com/huydq/test/internal/domain/object/permission"
	"github.com/huydq/test/internal/infrastructure/persistence/permission/dto"
	screenConvert "github.com/huydq/test/internal/infrastructure/persistence/screen/convert"
	screenDto "github.com/huydq/test/internal/infrastructure/persistence/screen/dto"
	persistence "github.com/huydq/test/internal/infrastructure/persistence/util"
	"gorm.io/gorm"
)

// ToPermission converts a Permission domain model to a Permission
func ToPermission(permission *modelPermission.Permission) *dto.Permission {
	if permission == nil {
		return nil
	}

	var screenDTO *screenDto.Screen
	if permission.Screen != nil {
		screenDTO = screenConvert.ToScreenDTO(permission.Screen)
	}

	result := &dto.Permission{
		ID:       permission.ID,
		Name:     permission.Name,
		Code:     string(permission.Code),
		ScreenID: permission.ScreenID,
		Screen:   screenDTO,
		BaseColumnTimestamp: persistence.BaseColumnTimestamp{
			CreatedAt: permission.CreatedAt,
			UpdatedAt: permission.UpdatedAt,
		},
	}

	// Handle the conversion from *time.Time to gorm.DeletedAt
	if permission.DeletedAt != nil {
		result.DeletedAt = gorm.DeletedAt{
			Time:  *permission.DeletedAt,
			Valid: true,
		}
	} else {
		result.DeletedAt = gorm.DeletedAt{
			Valid: false,
		}
	}

	return result
}

// ToPermissionModel converts a Permission to a Permission domain model
func ToPermissionModel(dtoObj *dto.Permission) *modelPermission.Permission {
	if dtoObj == nil {
		return nil
	}

	var screenModel *modelScreen.Screen
	if dtoObj.Screen != nil {
		screenModel = screenConvert.ToScreenModel(dtoObj.Screen)
	}

	result := &modelPermission.Permission{
		ID:       dtoObj.ID,
		Name:     dtoObj.Name,
		Code:     objectPermission.PermissionCode(dtoObj.Code),
		ScreenID: dtoObj.ScreenID,
		Screen:   screenModel,
		BaseColumnTimestamp: util.BaseColumnTimestamp{
			CreatedAt: dtoObj.CreatedAt,
			UpdatedAt: dtoObj.UpdatedAt,
		},
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

// ToPermissions converts a list of Permission domain models to a list of Permissions
func ToPermissions(permissions []*modelPermission.Permission) []*dto.Permission {
	if permissions == nil {
		return nil
	}

	result := make([]*dto.Permission, len(permissions))
	for i, permission := range permissions {
		result[i] = ToPermission(permission)
	}
	return result
}

// ToPermissionModels converts a list of Permissions to a list of Permission domain models
func ToPermissionModels(dtos []*dto.Permission) []*modelPermission.Permission {
	if dtos == nil {
		return nil
	}

	result := make([]*modelPermission.Permission, len(dtos))
	for i, dto := range dtos {
		result[i] = ToPermissionModel(dto)
	}
	return result
}
