package validator

import "github.com/vnlab/makeshop-payment/src/lib/validator"

type UpdateUserRequest struct {
	FullName   *string `json:"full_name,omitempty" validate:"omitempty,min=1,max=200"`
	RoleID     *int    `json:"role_id,omitempty" validate:"omitempty,min=1"`
	Email      *string `json:"email,omitempty" validate:"omitempty,email"`
	EnabledMFA *bool   `json:"enabled_mfa,omitempty"`
	Password   *string `json:"password" validate:"omitempty,password_policy"`
}

func (r *UpdateUserRequest) Validate() error {
	v := validator.GetValidate()
	return v.Struct(r)
}
