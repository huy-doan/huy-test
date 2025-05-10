package dto

import (
	"time"

	util "github.com/huydq/test/internal/domain/object/basedatetime"
	approvalDto "github.com/huydq/test/internal/infrastructure/persistence/approval/dto"
	userDto "github.com/huydq/test/internal/infrastructure/persistence/user/dto"
)

type PayoutDTO struct {
	ID int
	util.BaseColumnTimestamp

	PayoutStatus          int
	Total                 float64
	TotalCount            int
	SendingDate           time.Time
	SentDate              time.Time
	AozoraTransferApplyNo string
	ApprovalID            *int
	UserID                int

	User     *userDto.UserDTO         `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Approval *approvalDto.ApprovalDTO `gorm:"foreignKey:ApprovalID"`
}

func (PayoutDTO) TableName() string {
	return "payout"
}
