package models

import "time"

const (
	PayoutStatusDraft     int = 1 // ドラフト
	PayoutStatusCreated   int = 2 // 振込データ作成済み
	PayoutStatusProcessed int = 3 // 送金手続き済み
)

type Payout struct {
	ID int `json:"id"`
	BaseColumnTimestamp

	PayoutStatus          int       `json:"payout_status"`
	Total                 float64   `json:"total"`
	TotalCount            int       `json:"total_count"`
	SendingDate           time.Time `json:"sending_date"`
	SentDate              time.Time `json:"sent_date"`
	AozoraTransferApplyNo string    `json:"aozora_transfer_apply_no"`
	ApprovalID            *int      `json:"approval_id"`
	UserID                int       `json:"user_id"`

	User     *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Approval *Approval `json:"approval,omitempty" gorm:"foreignKey:ApprovalID"`

	PayoutRecordCount     int     `json:"payout_record_count" gorm:"-"`
	PayoutRecordSumAmount float64 `json:"payout_record_sum_amount" gorm:"-"`
}

func (Payout) TableName() string {
	return "payout"
}
