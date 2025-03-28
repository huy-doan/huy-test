package validator

import "github.com/huydq/demo/src/lib/validator"

type PaginationRequest struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=1000"`
}

// Validate kiểm tra dữ liệu của PaginationRequest
func (r *PaginationRequest) Validate() error {
	v := validator.GetValidate()
	return v.Struct(r)
}
