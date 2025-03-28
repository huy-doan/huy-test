package validator

import "github.com/huydq/demo/src/lib/validator"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (r *LoginRequest) Validate() error {
	return validator.ValidateStruct(r)
}
