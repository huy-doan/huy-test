package mapper

import (
	"time"

	controllerBase "github.com/huydq/test/internal/controller/base"
	"github.com/huydq/test/internal/domain/model/payout"
	"github.com/huydq/test/internal/domain/model/util"
	objectPayout "github.com/huydq/test/internal/domain/object/payout"
)

// PayoutListRequest represents the request parameters for listing payouts
type PayoutListRequest struct {
	controllerBase.PaginationRequest
	CreatedAt    string `query:"created_at" validate:"omitempty"`
	SendingDate  string `query:"sending_date" validate:"omitempty"`
	SentDate     string `query:"sent_date" validate:"omitempty"`
	PayoutStatus int    `query:"payout_status" validate:"omitempty"`
}

// ToPayoutFilter converts the request to a payout filter
func (r *PayoutListRequest) ToPayoutFilter() *payout.PayoutFilter {
	filter := payout.NewPayoutFilter()

	// Set pagination
	if r.Page > 0 {
		filter.SetPagination(r.Page, r.PageSize)
	}

	if r.SortField != "" {
		sortDirection := util.Ascending
		if r.SortOrder == "desc" {
			sortDirection = util.Descending
		}
		filter.SetSort(r.SortField, sortDirection)
	}

	if r.CreatedAt != "" {
		createdAt, err := time.Parse(time.RFC3339, r.CreatedAt)
		if err == nil {
			filter.CreatedAt = &createdAt
		}
	}

	if r.SendingDate != "" {
		sendingDate, err := time.Parse(time.RFC3339, r.SendingDate)
		if err == nil {
			filter.SendingDate = &sendingDate
		}
	}

	if r.SentDate != "" {
		sentDate, err := time.Parse(time.RFC3339, r.SentDate)
		if err == nil {
			filter.SentDate = &sentDate
		}
	}

	if r.PayoutStatus != 0 {
		status, success := objectPayout.GetPayoutStatusFromInt(r.PayoutStatus)
		if success {
			filter.PayoutStatus = &status
		}
	}

	return filter
}
