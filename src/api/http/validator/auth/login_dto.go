package validator

import "github.com/huydq/test/src/lib/validator"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (r *LoginRequest) Validate() error {
	v := validator.GetValidate()
	return v.Struct(r)
}
