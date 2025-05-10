package mapper

import (
	"github.com/huydq/test/internal/domain/model/user"
	"github.com/labstack/echo/v4"
)

// UpdateUserSuccessResponse represents the standardized response for user update
type UpdateUserSuccessResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    *DetailedUserData `json:"data"`
}

// UpdateUserSuccessMapper handles mapping for user update responses
type UpdateUserSuccessMapper struct {
	ctx echo.Context
}

// NewUpdateUserSuccessMapper creates a new mapper for user update responses
func NewUpdateUserSuccessMapper(ctx echo.Context) *UpdateUserSuccessMapper {
	return &UpdateUserSuccessMapper{
		ctx: ctx,
	}
}

// ToUpdateUserSuccessResponse creates a user update success response
func (m *UpdateUserSuccessMapper) ToUpdateUserSuccessResponse(user *user.User) *UpdateUserSuccessResponse {
	var mfaTypeData *MFATypeData
	if user.EnabledMFA {
		mfaType := "Email"
		mfaTypeData = &MFATypeData{
			ID:       user.MFAType,
			Title:    mfaType,
			IsActive: true,
		}
	}

	var roleData *RoleData
	if user.Role != nil {
		roleData = &RoleData{
			ID:   user.Role.ID,
			Name: user.Role.Name,
			Code: string(user.Role.Code),
		}
	}

	return &UpdateUserSuccessResponse{
		Success: true,
		Message: "ユーザー情報が正常に更新されました",
		Data: &DetailedUserData{
			ID:         user.ID,
			Email:      user.Email,
			FullName:   user.FullName,
			EnabledMFA: user.EnabledMFA,
			MFAType:    mfaTypeData,
			Role:       roleData,
		},
	}
}
