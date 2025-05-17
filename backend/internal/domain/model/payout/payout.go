package model

import (
	"time"

	userModel "github.com/huydq/test/internal/domain/model/user"
	util "github.com/huydq/test/internal/domain/object/basedatetime"
	object "github.com/huydq/test/internal/domain/object/payout"
)

type Payout struct {
	ID int
	util.BaseColumnTimestamp

	PayoutStatus          object.PayoutStatus
	Total                 float64
	TotalCount            int
	SendingDate           time.Time
	SentDate              time.Time
	AozoraTransferApplyNo string
	ApprovalID            *int
	UserID                int
	User                  *userModel.User

	PayoutRecordCount     int
	PayoutRecordSumAmount float64
}

type NewPayoutParams struct {
	ID int
	util.BaseColumnTimestamp
	PayoutStatus          object.PayoutStatus
	Total                 float64
	TotalCount            int
	SendingDate           time.Time
	SentDate              time.Time
	AozoraTransferApplyNo string
	ApprovalID            *int
	UserID                int
	User                  *userModel.User
	PayoutRecordCount     int
	PayoutRecordSumAmount float64
}

func NewPayout(params NewPayoutParams) *Payout {
	return &Payout{
		ID:                    params.ID,
		PayoutStatus:          params.PayoutStatus,
		Total:                 params.Total,
		TotalCount:            params.TotalCount,
		SendingDate:           params.SendingDate,
		SentDate:              params.SentDate,
		AozoraTransferApplyNo: params.AozoraTransferApplyNo,
		ApprovalID:            params.ApprovalID,
		UserID:                params.UserID,
		User:                  params.User,
		PayoutRecordCount:     params.PayoutRecordCount,
		PayoutRecordSumAmount: params.PayoutRecordSumAmount,
		BaseColumnTimestamp:   params.BaseColumnTimestamp,
	}
}

func (p *Payout) IsProcessed() bool {
	return p.PayoutStatus == object.PayoutStatusProcessed
}

func (p *Payout) CanBeProcessed() bool {
	return p.PayoutStatus == object.PayoutStatusCreated
}

func (p *Payout) MarkAsProcessed(sentDate time.Time) {
	p.PayoutStatus = object.PayoutStatusProcessed
	p.SentDate = sentDate
}
