package filter

import (
	"time"
)

// PayoutFilter represents filtering options for payouts
type PayoutFilter struct {
	BaseFilter
	CreatedAt    *time.Time
	SendingDate  *time.Time
	SentDate     *time.Time
	PayoutStatus *int
}

// NewPayoutFilter creates a new PayoutFilter with valid sort fields
func NewPayoutFilter() *PayoutFilter {
	filter := &PayoutFilter{}
	filter.ValidSortFields = map[string]bool{
		"id":            true,
		"created_at":    true,
		"updated_at":    true,
		"sending_date":  true,
		"sent_date":     true,
		"payout_status": true,
	}

	filter.SetPagination(1, 10)

	return filter
}

// ApplyFilters applies all filter conditions based on the filter fields
func (f *PayoutFilter) ApplyFilters() {
	// Apply date filters
	f.AddDateFilter("created_at", Equal, f.CreatedAt)
	f.AddDateFilter("sending_date", Equal, f.SendingDate)
	f.AddDateFilter("sent_date", Equal, f.SentDate)

	// Apply other conditions
	if f.PayoutStatus != nil {
		f.AddCondition("payout_status", Equal, *f.PayoutStatus)
	}
}
