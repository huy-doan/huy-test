package models

// AuditLogType represents the type of audit log
type AuditLogType struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
	BaseColumnTimestamp
}

// TableName specifies the table name for AuditLogType
func (AuditLogType) TableName() string {
	return "audit_log_type"
}

// Audit log type constants
const (
	AuditTypeUserLogin                  int = 1
	AuditTypeUserLogout                 int = 2
	AuditTypePasswordChange             int = 3
	AuditTypeUserCreate                 int = 4
	AuditTypeUserUpdate                 int = 5
	AuditTypeUserDelete                 int = 6
	AuditTypeRoleAssign                 int = 7
	AuditTypeRoleUpdate                 int = 8
	AuditTypeTransactionBusinessApprove int = 9
	AuditTypeTransactionBusinessReject  int = 10
	AuditTypeTransactionAccountApprove  int = 11
	AuditTypeTransactionAccountReject   int = 12
	AuditTypeTransferBusinessApprove    int = 13
	AuditTypeTransferBusinessReject     int = 14
	AuditTypeTransferAccountApprove     int = 15
	AuditTypeTransferAccountReject      int = 16
	AuditTypeMFAEnabled                 int = 17
	AuditTypeMFADisabled                int = 18
	AuditTypeAccountLocked              int = 19
	AuditTypeAccountUnlocked            int = 20
	AuditTypeAPIAccess                  int = 21
)
