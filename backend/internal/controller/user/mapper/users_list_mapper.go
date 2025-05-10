package mapper

import (
	"github.com/huydq/test/internal/datastructure/inputdata"
	"github.com/huydq/test/internal/datastructure/outputdata"
	"github.com/huydq/test/internal/domain/model/user"
	"github.com/huydq/test/internal/pkg/validator"
)

type UserListRequest struct {
	Page      int    `json:"page" query:"page" default:"1" validate:"min=1"`
	PageSize  int    `json:"page_size" query:"page_size" default:"10" validate:"min=1"`
	Search    string `json:"search" query:"search" validate:"omitempty,max=255"`
	RoleID    *int   `json:"role_id" query:"role_id" validate:"omitempty,min=1"`
	SortField string `json:"sort_field" query:"sort_field" validate:"omitempty"`
	SortOrder string `json:"sort_order" query:"sort_order" validate:"omitempty,oneof=asc desc"`
}

func (r *UserListRequest) Validate() error {
	v := validator.NewValidator()
	return v.Validate(r)
}

func (r *UserListRequest) ToUserListInputData() *inputdata.UserListInputData {
	return &inputdata.UserListInputData{
		Page:      r.Page,
		PageSize:  r.PageSize,
		Search:    r.Search,
		RoleID:    r.RoleID,
		SortField: r.SortField,
		SortOrder: r.SortOrder,
	}
}

type UserListResponse struct {
	Users      []*user.User `json:"users"`
	TotalPages int          `json:"total_pages"`
	Total      int          `json:"total"`
	Page       int          `json:"page"`
}

func (r *UserListResponse) ToResponseMap() map[string]interface{} {
	return map[string]interface{}{
		"success": true,
		"message": "Users retrieved successfully",
		"data": map[string]interface{}{
			"users":       r.Users,
			"total_pages": r.TotalPages,
			"total":       r.Total,
			"page":        r.Page,
		},
	}
}

func CreateUserListResponse(users []*user.User, totalPages int, total int, page int) *UserListResponse {
	return &UserListResponse{
		Users:      users,
		TotalPages: totalPages,
		Total:      total,
		Page:       page,
	}
}

func FromUserListOutputData(data *outputdata.UserListOutputData) *UserListResponse {
	return &UserListResponse{
		Users:      data.Users,
		TotalPages: data.TotalPages,
		Total:      data.TotalCount,
		Page:       data.CurrentPage,
	}
}
