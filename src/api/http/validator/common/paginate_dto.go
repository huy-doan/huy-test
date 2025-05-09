package validator

import "github.com/huydq/test/src/lib/validator"

type PaginationRequest struct {
	Page     int `form:"page" validate:"omitempty,min=1"`
	PageSize int `form:"page_size" validate:"omitempty,min=1,max=1000"`
}

func (r *PaginationRequest) Validate() error {
	v := validator.GetValidate()
	return v.Struct(r)
}
