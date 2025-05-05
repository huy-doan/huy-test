package models

// AuditLog represents an audit log entry in the system
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
	BaseColumnTimestamp

	// Relationships
	User *User `json:"user"`
}

// TableName specifies the table name for AuditLog
func (al *AuditLog) TableName() string {
	return "audit_log"
}

// Audit Log Type Constants
const (
	// User related audit log types
	AuditLogTypeLogin          = "ログイン"
	AuditLogTypeLogout         = "ログアウト"
	AuditLogTypePasswordChange = "パスワード変更"
	AuditLogTypePasswordReset  = "パスワードリセット"
	AuditLogTypeUserCreate     = "ユーザー作成"
	AuditLogTypeUserUpdate     = "ユーザー編集"
	AuditLogTypeUserDelete     = "ユーザー削除"
	AuditLogTypeRoleChange     = "ロール変更"

	// Two-factor authentication related
	AuditLogType2FAEnable  = "２段階認証有効"
	AuditLogType2FADisable = "２段階認証無効"

	// Payout related audit log types
	AuditLogTypePayoutRequest  = "出金申請"
	AuditLogTypePayoutApproval = "出金承認"
	AuditLogTypePayoutReject   = "出金却下"
	AuditLogTypePayoutResend   = "振込再依頼"
	AuditLogTypePayoutMarkSent = "振込を送金済みとする"

	// Payin related audit log types
	AuditLogTypeManualPayinImport = "手動入金取り込み"

	// Report related audit log types
	AuditLogTypePayinReportDownload = "入金レポートをダウンロード"
	AuditLogTypePayinDetailDownload = "入金明細をダウンロード"

	// Merchant related audit log types
	AuditLogTypeMerchantStatusUpload = "加盟店審査状況をアップロード"

	// API related audit log types
	AuditLogTypeExternalAPIAccess = "外部APIアクセス"
)
