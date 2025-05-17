package model

import (
	"time"

	"github.com/huydq/test/internal/domain/model/util"
	objectPayout "github.com/huydq/test/internal/domain/object/payout"
)

type PayoutFilter struct {
	util.BaseFilter
	CreatedAt    *time.Time
	SendingDate  *time.Time
	SentDate     *time.Time
	PayoutStatus *objectPayout.PayoutStatus
}

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

func (f *PayoutFilter) ApplyFilters() {
	f.AddDateFilter("created_at", util.Equal, f.CreatedAt)
	f.AddDateFilter("sending_date", util.Equal, f.SendingDate)
	f.AddDateFilter("sent_date", util.Equal, f.SentDate)

	if f.PayoutStatus != nil {
		f.AddCondition("payout_status", util.Equal, int(*f.PayoutStatus))
	}
}
