package model

import (
	"fmt"
	"time"

	userModel "github.com/huydq/test/internal/domain/model/user"
	object "github.com/huydq/test/internal/domain/object/audit_log"
	util "github.com/huydq/test/internal/domain/object/basedatetime"
)

// Description constants for audit log events
const (
	// User-related descriptions
	DescLogin          = "ログインしました。"
	DescLogout         = "ログアウトしました。"
	DescPasswordChange = "パスワード変更しました。"
	DescPasswordReset  = "ユーザー（%d）のパスワードパスワードリセット。"
	DescUserCreate     = "ユーザー（%d）を作成しました。"
	DescUserUpdate     = "ユーザー（%d）編集しました。"
	DescUserDelete     = "ユーザー（%d）無効にしました。"
	DescRoleChange     = "ユーザー（%d）を「%s」ロールに変更しました。"

	// Two-factor authentication descriptions
	Desc2FAEnable  = "ユーザー（%d）の２段階認証を有効しました。"
	Desc2FADisable = "ユーザー（%d）の２段階認証を無効しました。"

	// Payout-related descriptions
	DescPayoutRequest  = "出金申請しました。"
	DescPayoutApproval = "出金承認しました。"
	DescPayoutReject   = "出金却下しました。"
	DescPayoutResend   = "振込再依頼を行いました。"
	DescPayoutMarkSent = "振込データを送金済みとしました。"

	// Payin-related descriptions
	DescManualPayinImport = "手動入金取り込みを行いました。"

	// Other descriptions
	DescMerchantStatusUpload = "加盟店審査状況をアップロードしました。"
	DescExternalAPIAccess    = "振込APIを実行しました。"
)

type AuditLog struct {
	ID            int
	UserID        *int
	User          *userModel.User
	AuditLogType  object.AuditLogType
	Description   *string
	TransactionID *int
	PayoutID      *int
	PayinID       *int
	TargetUserID  *int
	NewRole       *string
	UserAgent     *object.UserAgent
	IPAddress     *object.IPAddress
	util.BaseColumnTimestamp
}

type AuditLogGenerator struct {
	UserID        *int
	AuditLogType  object.AuditLogType
	Description   *string
	TransactionID *int
	PayoutID      *int
	PayinID       *int
	UserAgent     *object.UserAgent
	IPAddress     *object.IPAddress
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
	TargetUserID  *int
	NewRole       *string
}

// NewAuditLogGenerator creates a new generator with required base fields
func NewAuditLogGenerator(userID *int, auditLogType object.AuditLogType, ipAddress *object.IPAddress, userAgent *object.UserAgent) *AuditLogGenerator {
	return &AuditLogGenerator{
		UserID:       userID,
		AuditLogType: auditLogType,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
	}
}

// Map of audit log types to their corresponding description templates
var descriptionTemplates = map[object.AuditLogType]string{
	object.AuditLogTypeLogin:                DescLogin,
	object.AuditLogTypeLogout:               DescLogout,
	object.AuditLogTypePasswordChange:       DescPasswordChange,
	object.AuditLogTypePasswordReset:        DescPasswordReset,
	object.AuditLogTypeUserCreate:           DescUserCreate,
	object.AuditLogTypeUserUpdate:           DescUserUpdate,
	object.AuditLogTypeUserDelete:           DescUserDelete,
	object.AuditLogTypeRoleChange:           DescRoleChange,
	object.AuditLogType2FAEnable:            Desc2FAEnable,
	object.AuditLogType2FADisable:           Desc2FADisable,
	object.AuditLogTypePayoutRequest:        DescPayoutRequest,
	object.AuditLogTypePayoutApproval:       DescPayoutApproval,
	object.AuditLogTypePayoutReject:         DescPayoutReject,
	object.AuditLogTypePayoutResend:         DescPayoutResend,
	object.AuditLogTypePayoutMarkSent:       DescPayoutMarkSent,
	object.AuditLogTypeManualPayinImport:    DescManualPayinImport,
	object.AuditLogTypeMerchantStatusUpload: DescMerchantStatusUpload,
	object.AuditLogTypeExternalAPIAccess:    DescExternalAPIAccess,
}

// getDescription returns the appropriate description based on the audit log type
func (g *AuditLogGenerator) getDescription() string {
	template, exists := descriptionTemplates[g.AuditLogType]
	if !exists {
		return string(g.AuditLogType)
	}

	switch g.AuditLogType {
	case object.AuditLogTypePasswordReset,
		object.AuditLogTypeUserCreate,
		object.AuditLogTypeUserUpdate,
		object.AuditLogTypeUserDelete,
		object.AuditLogType2FAEnable,
		object.AuditLogType2FADisable:
		if g.TargetUserID != nil {
			return fmt.Sprintf(template, *g.TargetUserID)
		}
	case object.AuditLogTypeRoleChange:
		if g.TargetUserID != nil && g.NewRole != nil {
			return fmt.Sprintf(template, *g.TargetUserID, *g.NewRole)
		}
	default:
		return template
	}

	return string(g.AuditLogType)
}

// Generate builds and returns a new AuditLog instance
func (g *AuditLogGenerator) Generate() *AuditLog {
	if g.Description == nil {
		defaultDesc := g.getDescription()
		g.Description = &defaultDesc
	}

	return &AuditLog{
		UserID:        g.UserID,
		AuditLogType:  g.AuditLogType,
		Description:   g.Description,
		TransactionID: g.TransactionID,
		PayoutID:      g.PayoutID,
		PayinID:       g.PayinID,
		UserAgent:     g.UserAgent,
		IPAddress:     g.IPAddress,
		BaseColumnTimestamp: util.BaseColumnTimestamp{
			CreatedAt: g.CreatedAt,
			UpdatedAt: g.UpdatedAt,
			DeletedAt: g.DeletedAt,
		},
	}
}

func NewAuditLog() *AuditLog {
	return &AuditLog{
		BaseColumnTimestamp: util.BaseColumnTimestamp{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

type NewAuditLogParams struct {
	UserID        *int                `json:"user_id"`
	User          *userModel.User     `json:"user"`
	AuditLogType  object.AuditLogType `json:"audit_log_type"`
	Description   *string             `json:"description"`
	TransactionID *int                `json:"transaction_id"`
	PayoutID      *int                `json:"payout_id"`
	PayinID       *int                `json:"payin_id"`
	UserAgent     *object.UserAgent   `json:"user_agent"`
	IPAddress     *object.IPAddress   `json:"ip_address"`
	util.BaseColumnTimestamp
}

func NewAuditLogWithData(
	userID *int,
	auditLogType object.AuditLogType,
	ipAddress *object.IPAddress,
	userAgent *object.UserAgent,
	targetUserID int,
	newRole string,
	payoutID *int,
	payinID *int,
) *AuditLog {
	auditLog := NewAuditLog()
	auditLog.UserID = userID
	auditLog.AuditLogType = auditLogType
	auditLog.IPAddress = ipAddress
	auditLog.UserAgent = userAgent

	// Special case for login where the target user ID becomes the user ID
	if auditLogType == object.AuditLogTypeLogin && targetUserID != 0 {
		auditLog.UserID = &targetUserID
	} else if targetUserID != 0 {
		auditLog.TargetUserID = &targetUserID
	}

	if newRole != "" {
		auditLog.NewRole = &newRole
	}

	if payoutID != nil {
		auditLog.PayoutID = payoutID
	}

	if payinID != nil {
		auditLog.PayinID = payinID
	}

	return auditLog
}
