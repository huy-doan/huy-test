package mapper

import (
	"github.com/huydq/test/internal/domain/model/user"
	"github.com/labstack/echo/v4"
)

// UserListSuccessResponse represents the standard response for user listing
type UserListSuccessResponse struct {
	Success bool                 `json:"success"`
	Message string               `json:"message"`
	Data    *UserListSuccessData `json:"data"`
}

// UserListSuccessData represents the data in the user list response
type UserListSuccessData struct {
	Users      []*DetailedUserData `json:"users"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
	Total      int                 `json:"total"`
	TotalPages int                 `json:"total_pages"`
}

// UserListSuccessMapper handles mapping for user list responses
type UserListSuccessMapper struct {
	ctx echo.Context
}

// NewUserListSuccessMapper creates a new mapper for user list responses
func NewUserListSuccessMapper(ctx echo.Context) *UserListSuccessMapper {
	return &UserListSuccessMapper{
		ctx: ctx,
	}
}

// ToUserListSuccessResponse creates a user list success response
func (m *UserListSuccessMapper) ToUserListSuccessResponse(
	users []*user.User,
	totalPages int,
	total int,
	page int,
	pageSize int,
) *UserListSuccessResponse {
	detailedUsers := make([]*DetailedUserData, 0, len(users))

	for _, u := range users {
		var mfaTypeData *MFATypeData
		if u.EnabledMFA {
			mfaType := "Email"
			mfaTypeData = &MFATypeData{
				ID:       u.MFAType,
				Title:    mfaType,
				IsActive: true,
			}
		}

		// Always create role data to ensure it appears in response
		roleData := &RoleData{
			ID:   u.RoleID,
			Name: "システム管理者",
			Code: "SYSTEM_ADMIN",
		}

		// Override with actual role data if available
		if u.Role != nil {
			roleData.ID = u.Role.ID
			roleData.Name = u.Role.Name
			roleData.Code = string(u.Role.Code)
		}

		detailedUsers = append(detailedUsers, &DetailedUserData{
			ID:         u.ID,
			Email:      u.Email,
			FullName:   u.FullName,
			EnabledMFA: u.EnabledMFA,
			MFAType:    mfaTypeData,
			Role:       roleData,
		})
	}

	return &UserListSuccessResponse{
		Success: true,
		Message: "ユーザー一覧を正常に取得しました",
		Data: &UserListSuccessData{
			Users:      detailedUsers,
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}
}
