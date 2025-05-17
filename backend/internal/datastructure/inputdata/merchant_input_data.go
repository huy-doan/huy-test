package inputdata

import (
	"time"
)

type MerchantListInputData struct {
	Page     int `json:"page" validate:"min=1"`
	PageSize int `json:"page_size" validate:"min=1,max=100"`

	Search         string     `json:"search"`
	ReviewStatus   []int      `json:"review_status"`
	CreatedAtStart *time.Time `json:"created_at_start"`
	CreatedAtEnd   *time.Time `json:"created_at_end"`

	SortField string `json:"sort_field"`
	SortOrder string `json:"sort_order" validate:"omitempty,oneof=asc desc"`
}
