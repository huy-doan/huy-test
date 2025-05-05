package validator

import "github.com/vnlab/makeshop-payment/src/lib/validator"

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password_policy"`
	FullName string `json:"full_name" validate:"required,min=1,max=200"`
}

func (r *RegisterRequest) Validate() error {
	v := validator.GetValidate()
	return v.Struct(r)
}
