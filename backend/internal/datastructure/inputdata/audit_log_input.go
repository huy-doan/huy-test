package inputdata

import (
	"time"
)

// CreateAuditLogInput represents the input data for creating a new audit log
type CreateAuditLogInput struct {
	UserID        *int    `json:"user_id"`
	AuditLogType  string  `json:"audit_log_type"`
	Description   *string `json:"description"`
	TransactionID *int    `json:"transaction_id"`
	PayoutID      *int    `json:"payout_id"`
	PayinID       *int    `json:"payin_id"`
	UserAgent     *string `json:"user_agent"`
	IPAddress     *string `json:"ip_address"`
}

// ListAuditLogInput represents the input data for listing audit logs
type ListAuditLogInput struct {
	UserID       *int       `json:"user_id"`
	AuditLogType *string    `json:"audit_log_type"`
	CreatedAt    *time.Time `json:"created_at"`
	Description  *string    `json:"description"`
	Page         int        `json:"page"`
	PageSize     int        `json:"page_size"`
	SortField    string     `json:"sort_field"`
	SortOrder    string     `json:"sort_order"`
}
