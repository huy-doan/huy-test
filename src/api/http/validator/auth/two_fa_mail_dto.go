package validator

import "github.com/huydq/test/src/lib/validator"

type VerifyRequest struct {
	Email string `json:"email" validate:"required,email"`
	Token string `json:"token" validate:"required"`
}

func (r *VerifyRequest) Validate() error {
	v := validator.GetValidate()
	return v.Struct(r)
}
