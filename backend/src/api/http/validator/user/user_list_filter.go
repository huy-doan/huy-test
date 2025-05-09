package validator

import (
	"github.com/huydq/test/src/lib/validator"
)

type UserListFilter struct {
	Page      int    `json:"page" default:"1" validate:"min=1"`
	PageSize  int    `json:"page_size" default:"10" validate:"min=1"`
	Search    string `json:"search" validate:"omitempty,max=255"`
	RoleID    *int   `json:"role_id" validate:"omitempty,min=1"`
	SortField string `json:"sort_field" validate:"omitempty"`
	SortOrder string `json:"sort_order" validate:"omitempty,oneof=asc desc"`
}

func (f *UserListFilter) Validate() error {
	v := validator.GetValidate()
	return v.Struct(f)
}
