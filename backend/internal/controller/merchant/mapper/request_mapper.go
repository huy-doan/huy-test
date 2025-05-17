package mapper

import (
	"time"

	"github.com/huydq/test/internal/datastructure/inputdata"
	generated "github.com/huydq/test/internal/pkg/api/generated"
)

func ToMerchantListInputData(request *generated.MerchantListRequest) *inputdata.MerchantListInputData {
	var startTime, endTime *time.Time

	if request.CreatedAtStart != "" {
		if t, err := time.Parse(time.RFC3339, request.CreatedAtStart); err == nil {
			startTime = &t
		}
	}

	if request.CreatedAtEnd != "" {
		if t, err := time.Parse(time.RFC3339, request.CreatedAtEnd); err == nil {
			endTime = &t
		}
	}

	return &inputdata.MerchantListInputData{
		Page:           request.Page,
		PageSize:       request.PageSize,
		Search:         request.Search,
		ReviewStatus:   request.ReviewStatus,
		CreatedAtStart: startTime,
		CreatedAtEnd:   endTime,
		SortField:      request.SortField,
		SortOrder:      request.SortOrder,
	}
}
