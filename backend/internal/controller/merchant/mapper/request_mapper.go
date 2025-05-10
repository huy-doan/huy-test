package mapper

import (
	"time"

	"github.com/huydq/test/internal/datastructure/inputdata"
)

type MerchantListRequest struct {
	Page     int `query:"page" json:"page" validate:"omitempty,min=1"`
	PageSize int `query:"page_size" json:"page_size" validate:"omitempty,min=1"`

	Search         string `query:"search" json:"search"`
	ReviewStatus   []int  `query:"review_status" json:"review_status"`
	CreatedAtStart string `query:"created_at_start" json:"created_at_start"`
	CreatedAtEnd   string `query:"created_at_end" json:"created_at_end"`

	SortField string `query:"sort_field" json:"sort_field"`
	SortOrder string `query:"sort_order" json:"sort_order" validate:"omitempty,oneof=asc desc"`
}

func (r *MerchantListRequest) ToMerchantListInputData() *inputdata.MerchantListInputData {
	var startTime, endTime *time.Time

	if r.CreatedAtStart != "" {
		if t, err := time.Parse(time.RFC3339, r.CreatedAtStart); err == nil {
			startTime = &t
		}
	}

	if r.CreatedAtEnd != "" {
		if t, err := time.Parse(time.RFC3339, r.CreatedAtEnd); err == nil {
			endTime = &t
		}
	}

	return &inputdata.MerchantListInputData{
		Page:           r.Page,
		PageSize:       r.PageSize,
		Search:         r.Search,
		ReviewStatus:   r.ReviewStatus,
		CreatedAtStart: startTime,
		CreatedAtEnd:   endTime,
		SortField:      r.SortField,
		SortOrder:      r.SortOrder,
	}
}
