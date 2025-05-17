package mapper

import (
	"github.com/huydq/test/internal/datastructure/inputdata"
	"github.com/huydq/test/internal/domain/model/user"
	generated "github.com/huydq/test/internal/pkg/api/generated"
	utils "github.com/huydq/test/internal/pkg/utils"
)

func ToCreateUserInputData(r *generated.CreateUserRequest) *inputdata.CreateUserInputData {
	email := string(r.Email)
	return &inputdata.CreateUserInputData{
		Email:      email,
		Password:   r.Password,
		FullName:   r.FullName,
		RoleID:     r.RoleId,
		EnabledMFA: r.EnabledMfa,
	}
}

func ToUpdateUserInputData(r *generated.UpdateUserRequest) *inputdata.UpdateUserInputData {
	var email *string
	if r.Email != nil {
		str := string(*r.Email)
		email = &str
	}
	return &inputdata.UpdateUserInputData{
		Email:      email,
		Password:   r.Password,
		FullName:   r.FullName,
		RoleID:     r.RoleId,
		EnabledMFA: r.EnabledMfa,
	}
}

func ToUserListInputData(r *generated.UserListRequest) *inputdata.UserListInputData {
	return &inputdata.UserListInputData{
		Page:      r.Page,
		PageSize:  r.PageSize,
		Search:    r.Search,
		RoleID:    r.RoleId,
		SortField: r.SortField,
		SortOrder: string(r.SortOrder),
	}
}

func ToUserListSuccessData(users []*user.User,
	totalPages int,
	total int,
	page int,
	pageSize int) *generated.UserListResponse {

	detailedUsers := make([]generated.User, 0, len(users))
	for _, u := range users {
		user := generated.User{
			Id:         utils.ToPtr(u.ID),
			Email:      utils.ToPtr(u.Email),
			FullName:   utils.ToPtr(u.FullName),
			EnabledMfa: utils.ToPtr(u.EnabledMFA),
			Role: &generated.Role{
				Id:   utils.ToPtr(u.Role.ID),
				Name: utils.ToPtr(u.Role.Name),
			},
			MfaType: &generated.MfaType{
				Id:       utils.ToPtr(u.MFAType),
				IsActive: utils.ToPtr(true),
				Title:    utils.ToPtr("Email"),
			},
			CreatedAt: utils.ToPtr(u.CreatedAt),
			UpdatedAt: utils.ToPtr(u.UpdatedAt),
		}

		detailedUsers = append(detailedUsers, user)
	}

	return &generated.UserListResponse{
		Users:      detailedUsers,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}

func ToDetailedUserData(user *user.User) *generated.User {
	return &generated.User{
		Id:         utils.ToPtr(user.ID),
		Email:      utils.ToPtr(user.Email),
		FullName:   utils.ToPtr(user.FullName),
		EnabledMfa: utils.ToPtr(user.EnabledMFA),
		Role: &generated.Role{
			Id:   utils.ToPtr(user.Role.ID),
			Name: utils.ToPtr(user.Role.Name),
		},
		MfaType: &generated.MfaType{
			Id:       utils.ToPtr(user.MFAType),
			IsActive: utils.ToPtr(true),
			Title:    utils.ToPtr("Email"),
		},
	}
}
