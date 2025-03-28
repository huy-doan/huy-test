package validator

import "github.com/huydq/ddd-project/src/lib/validator"

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (r *LoginRequest) Validate() error {
	v := validator.GetValidate()
	return v.Struct(r)
}
