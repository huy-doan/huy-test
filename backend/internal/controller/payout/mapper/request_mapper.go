package mapper

import (
	"time"

	payoutModel "github.com/huydq/test/internal/domain/model/payout"
	"github.com/huydq/test/internal/domain/model/util"
	objectPayout "github.com/huydq/test/internal/domain/object/payout"
	generated "github.com/huydq/test/internal/pkg/api/generated"
)

// ToPayoutFilter converts the request to a payout filter
func ToPayoutFilter(request *generated.PayoutListRequest) *payoutModel.PayoutFilter {
	filter := payoutModel.NewPayoutFilter()

	// Set pagination
	if request.Page > 0 {
		filter.SetPagination(request.Page, request.PageSize)
	}

	if request.SortField != "" {
		sortDirection := util.Ascending
		if request.SortOrder == "desc" {
			sortDirection = util.Descending
		}
		filter.SetSort(request.SortField, sortDirection)
	}

	if request.CreatedAt != "" {
		createdAt, err := time.Parse(time.RFC3339, request.CreatedAt)
		if err == nil {
			filter.CreatedAt = &createdAt
		}
	}

	if request.SendingDate != "" {
		sendingDate, err := time.Parse(time.RFC3339, request.SendingDate)
		if err == nil {
			filter.SendingDate = &sendingDate
		}
	}

	if request.SentDate != "" {
		sentDate, err := time.Parse(time.RFC3339, request.SentDate)
		if err == nil {
			filter.SentDate = &sentDate
		}
	}

	if request.PayoutStatus != 0 {
		status, success := objectPayout.GetPayoutStatusFromInt(request.PayoutStatus)
		if success {
			filter.PayoutStatus = &status
		}
	}

	return filter
}
