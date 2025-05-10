package model

import (
	object "github.com/huydq/test/internal/domain/object/audit_log"
	util "github.com/huydq/test/internal/domain/object/basedatetime"
)

type AuditLog struct {
	ID            int                 `json:"id"`
	UserID        *int                `json:"user_id"`
	AuditLogType  object.AuditLogType `json:"audit_log_type"`
	Description   *string             `json:"description"`
	TransactionID *int                `json:"transaction_id"`
	PayoutID      *int                `json:"payout_id"`
	PayinID       *int                `json:"payin_id"`
	UserAgent     *object.UserAgent   `json:"user_agent"`
	IPAddress     *object.IPAddress   `json:"ip_address"`
	util.BaseColumnTimestamp
}

func NewAuditLog(param NewAuditLogParams) *AuditLog {
	return &AuditLog{
		UserID:              param.UserID,
		AuditLogType:        param.AuditLogType,
		Description:         param.Description,
		TransactionID:       param.TransactionID,
		PayoutID:            param.PayoutID,
		PayinID:             param.PayinID,
		UserAgent:           param.UserAgent,
		IPAddress:           param.IPAddress,
		BaseColumnTimestamp: param.BaseColumnTimestamp,
	}
}

type NewAuditLogParams struct {
	UserID        *int                `json:"user_id"`
	AuditLogType  object.AuditLogType `json:"audit_log_type"`
	Description   *string             `json:"description"`
	TransactionID *int                `json:"transaction_id"`
	PayoutID      *int                `json:"payout_id"`
	PayinID       *int                `json:"payin_id"`
	UserAgent     *object.UserAgent   `json:"user_agent"`
	IPAddress     *object.IPAddress   `json:"ip_address"`
	util.BaseColumnTimestamp
}
