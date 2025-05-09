package validator

import "github.com/huydq/test/src/lib/validator"

type UpdateUserRequest struct {
	LastName      *string `json:"last_name,omitempty" validate:"omitempty,min=1,max=100"`
	FirstName     *string `json:"first_name,omitempty" validate:"omitempty,min=1,max=100"`
	LastNameKana  *string `json:"last_name_kana,omitempty" validate:"omitempty,min=1,max=100,kana"`
	FirstNameKana *string `json:"first_name_kana,omitempty" validate:"omitempty,min=1,max=100,kana"`
	RoleID        *int    `json:"role_id,omitempty" validate:"omitempty,min=1"`
	Email         *string `json:"email,omitempty" validate:"omitempty,email"`
	EnabledMFA    *bool   `json:"enabled_mfa,omitempty"`
}

func (r *UpdateUserRequest) Validate() error {
	v := validator.GetValidate()
	return v.Struct(r)
}
