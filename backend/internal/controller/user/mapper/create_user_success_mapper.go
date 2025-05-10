package mapper

import (
	"github.com/huydq/test/internal/domain/model/user"
	"github.com/labstack/echo/v4"
)

type CreateUserSuccessResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    *DetailedUserData `json:"data"`
}

type CreateUserSuccessMapper struct {
	ctx echo.Context
}

func NewCreateUserSuccessMapper(ctx echo.Context) *CreateUserSuccessMapper {
	return &CreateUserSuccessMapper{
		ctx: ctx,
	}
}

func (m *CreateUserSuccessMapper) ToCreateUserSuccessResponse(user *user.User) *CreateUserSuccessResponse {
	var mfaTypeData *MFATypeData
	if user.EnabledMFA {
		mfaType := "Email"
		mfaTypeData = &MFATypeData{
			ID:       user.MFAType,
			Title:    mfaType,
			IsActive: true,
		}
	}

	roleData := &RoleData{
		ID:   user.RoleID,
		Name: "システム管理者",
		Code: "SYSTEM_ADMIN",
	}

	if user.Role != nil {
		roleData.ID = user.Role.ID
		roleData.Name = user.Role.Name
		roleData.Code = string(user.Role.Code)
	}

	return &CreateUserSuccessResponse{
		Success: true,
		Message: "ユーザーが正常に作成されました",
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
