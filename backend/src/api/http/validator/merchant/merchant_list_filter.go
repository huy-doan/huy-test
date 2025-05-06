package validator

import (
	"time"

	"github.com/vnlab/makeshop-payment/src/lib/validator"
)

type MerchantListFilter struct {
	Page           int        `json:"page" default:"1" validate:"min=1"`
	PageSize       int        `json:"page_size" default:"10" validate:"min=1"`
	Search         string     `json:"search" validate:"omitempty,max=255"`
	ReviewStatus   []int      `json:"review_status" validate:"omitempty,dive,min=1,max=3"`
	CreatedAtStart *time.Time `json:"created_at_start" validate:"omitempty"`
	CreatedAtEnd   *time.Time `json:"created_at_end" validate:"omitempty"`
	SortField      string     `json:"sort_field" validate:"omitempty"`
	SortOrder      string     `json:"sort_order" validate:"omitempty,oneof=asc desc"`
}

func (f *MerchantListFilter) Validate() error {
	v := validator.GetValidate()
	return v.Struct(f)
}
