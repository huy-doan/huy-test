package validator

import "github.com/vnlab/makeshop-payment/src/lib/validator"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (r *LoginRequest) Validate() error {
	return validator.ValidateStruct(r)
}
