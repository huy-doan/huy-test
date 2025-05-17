package object

type AuditLogType string

const (
	// User related audit log types
	AuditLogTypeLogin          AuditLogType = "ログイン"
	AuditLogTypeLogout         AuditLogType = "ログアウト"
	AuditLogTypePasswordChange AuditLogType = "パスワード変更"
	AuditLogTypePasswordReset  AuditLogType = "パスワードリセット"
	AuditLogTypeUserCreate     AuditLogType = "ユーザー作成"
	AuditLogTypeUserUpdate     AuditLogType = "ユーザー編集"
	AuditLogTypeUserDelete     AuditLogType = "ユーザー削除"
	AuditLogTypeRoleChange     AuditLogType = "ロール変更"

	// Two-factor authentication related
	AuditLogType2FAEnable  AuditLogType = "２段階認証有効"
	AuditLogType2FADisable AuditLogType = "２段階認証無効"

	// Payout related audit log types
	AuditLogTypePayoutRequest  AuditLogType = "出金申請"
	AuditLogTypePayoutApproval AuditLogType = "出金承認"
	AuditLogTypePayoutReject   AuditLogType = "出金却下"
	AuditLogTypePayoutResend   AuditLogType = "振込再依頼"
	AuditLogTypePayoutMarkSent AuditLogType = "振込を送金済みとする"

	// Payin related audit log types
	AuditLogTypeManualPayinImport AuditLogType = "手動入金取り込み"

	// Report related audit log types
	AuditLogTypePayinReportDownload AuditLogType = "入金レポートをダウンロード"
	AuditLogTypePayinDetailDownload AuditLogType = "入金明細をダウンロード"

	// Merchant related audit log types
	AuditLogTypeMerchantStatusUpload AuditLogType = "加盟店審査状況をアップロード"

	// API related audit log types
	AuditLogTypeExternalAPIAccess AuditLogType = "外部APIアクセス"
)
