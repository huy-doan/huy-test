package dto

import (
	userDto "github.com/huydq/test/internal/infrastructure/persistence/user/dto"
	persistence "github.com/huydq/test/internal/infrastructure/persistence/util"
)

type AuditLog struct {
	ID            int     `json:"id"`
	UserID        *int    `json:"user_id"`
	AuditLogType  string  `json:"audit_log_type"`
	Description   *string `json:"description"`
	TransactionID *int    `json:"transaction_id"`
	PayoutID      *int    `json:"payout_id"`
	PayinID       *int    `json:"payin_id"`
	UserAgent     *string `json:"user_agent"`
	IPAddress     *string `json:"ip_address"`
	persistence.BaseColumnTimestamp

	// Relationships
	User *userDto.User `json:"user"`
}

func (al *AuditLog) TableName() string {
	return "audit_log"
}
