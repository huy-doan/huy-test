package validator

import "github.com/vnlab/makeshop-payment/src/lib/validator"

type VerifyRequest struct {
	Email string `json:"email" validate:"required,email"`
	Token string `json:"token" validate:"required"`
}

func (r *VerifyRequest) Validate() error {
	v := validator.GetValidate()
	return v.Struct(r)
}

// ResendCodeRequest represents the request to resend a 2FA code
type ResendCodeRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func (r *ResendCodeRequest) Validate() error {
	v := validator.GetValidate()
	return v.Struct(r)
}
