package converter

import (
	"github.com/huydq/test/internal/datastructure/outputdata"
	modelPermission "github.com/huydq/test/internal/domain/model/permission"
	modelScreen "github.com/huydq/test/internal/domain/model/screen"
)

// PermissionModelToOutput converts a Permission domain model to a PermissionOutput
func PermissionModelToOutput(permission *modelPermission.Permission) *outputdata.PermissionOutput {
	if permission == nil {
		return nil
	}

	var screenOutput *outputdata.ScreenOutput
	if permission.Screen != nil {
		screenOutput = ScreenModelToOutput(permission.Screen)
	}

	return &outputdata.PermissionOutput{
		ID:        permission.ID,
		Name:      permission.Name,
		Code:      string(permission.Code),
		ScreenID:  permission.ScreenID,
		Screen:    screenOutput,
		CreatedAt: permission.CreatedAt,
		UpdatedAt: permission.UpdatedAt,
	}
}

// PermissionModelsToOutputs converts a slice of Permission domain models to a slice of PermissionOutputs
func PermissionModelsToOutputs(permissions []*modelPermission.Permission) []*outputdata.PermissionOutput {
	if permissions == nil {
		return nil
	}

	outputs := make([]*outputdata.PermissionOutput, len(permissions))
	for i, permission := range permissions {
		outputs[i] = PermissionModelToOutput(permission)
	}
	return outputs
}

// ScreenModelToOutput converts a Screen domain model to a ScreenOutput
func ScreenModelToOutput(screen *modelScreen.Screen) *outputdata.ScreenOutput {
	if screen == nil {
		return nil
	}

	return &outputdata.ScreenOutput{
		ID:         screen.ID,
		Name:       screen.Name,
		ScreenCode: screen.ScreenCode,
		ScreenPath: screen.ScreenPath,
		CreatedAt:  screen.CreatedAt,
		UpdatedAt:  screen.UpdatedAt,
	}
}

// ScreenModelsToOutputs converts a slice of Screen domain models to a slice of ScreenOutputs
func ScreenModelsToOutputs(screens []*modelScreen.Screen) []*outputdata.ScreenOutput {
	if screens == nil {
		return nil
	}

	outputs := make([]*outputdata.ScreenOutput, len(screens))
	for i, screen := range screens {
		outputs[i] = ScreenModelToOutput(screen)
	}
	return outputs
}
