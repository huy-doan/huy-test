package models

// AuditLog represents an audit log entry in the system
type AuditLog struct {
	ID            int     `json:"id"`
	UserID        *int    `json:"user_id"`
	AuditTypeID   uint8   `json:"audit_type_id"`
	Description   *string `json:"description"`
	TransactionID *int    `json:"transaction_id"`
	OutcomingID   *int    `json:"outcoming_id"`
	UserAgent     *string `json:"user_agent"`
	IPAddress     *string `json:"ip_address"`
	BaseColumnTimestamp

	// Relationships
	User               *User               `json:"user"`
	MasterAuditLogType *MasterAuditLogType `json:"master_audit_log_type"`
}

// TableName specifies the table name for AuditLog
func (al *AuditLog) TableName() string {
	return "audit_logs"
}
