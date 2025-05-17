package outputdata

import (
	"time"
)

// AuditLogOutput represents the output data for a single audit log
type AuditLogOutput struct {
	ID            int       `json:"id"`
	UserID        *int      `json:"user_id"`
	AuditLogType  string    `json:"audit_log_type"`
	Description   *string   `json:"description"`
	TransactionID *int      `json:"transaction_id"`
	PayoutID      *int      `json:"payout_id"`
	PayinID       *int      `json:"payin_id"`
	UserAgent     *string   `json:"user_agent"`
	IPAddress     *string   `json:"ip_address"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ListAuditLogOutput represents the output data for a list of audit logs
type ListAuditLogOutput struct {
	AuditLogs  []*AuditLogOutput `json:"audit_logs"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
	Total      int64             `json:"total"`
}
