package service

import (
	"context"
	"fmt"

	model "github.com/huydq/test/internal/domain/model/audit_log"
	object "github.com/huydq/test/internal/domain/object/audit_log"
	repository "github.com/huydq/test/internal/domain/repository/audit_log"
)

// Description constants for audit log events
const (
	// User-related descriptions
	DescLogin          = "ログイン"
	DescLogout         = "ログアウト"
	DescPasswordChange = "パスワード変更"
	DescPasswordReset  = "ユーザーID %d のパスワードをリセット"
	DescUserCreate     = "ユーザーID %d を作成"
	DescUserUpdate     = "ユーザーID %d を更新"
	DescUserDelete     = "ユーザーID %d を削除"
	DescRoleChange     = "ユーザーID %d のロールを %s に変更"

	// Two-factor authentication descriptions
	Desc2FAEnable  = "ユーザーID %d の２段階認証を有効化"
	Desc2FADisable = "ユーザーID %d の２段階認証を無効化"

	// Payout-related descriptions
	DescPayoutRequest  = "出金申請"
	DescPayoutApproval = "出金承認"
	DescPayoutReject   = "出金却下"
	DescPayoutResend   = "振込再依頼"
	DescPayoutMarkSent = "振込を送金済みとする"

	// Payin-related descriptions
	DescManualPayinImport = "手動入金取り込み"

	// Other descriptions
	DescMerchantStatusUpload = "加盟店審査状況をアップロード"
	DescExternalAPIAccess    = "外部APIアクセス"
)

type AuditLogService interface {
	// Core operations
	Create(ctx context.Context, auditLog *model.AuditLog) error
	List(ctx context.Context, filter *model.AuditLogFilter) ([]*model.AuditLog, int, int64, error)
	// GetUserIDsWithAuditLogs(ctx context.Context) ([]int, error)

	// Generic event logging
	LogEventByType(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent, auditLogType object.AuditLogType, description *string) error

	// User-related events
	LogLoginEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent) error
	LogLogoutEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent) error
	LogPasswordChangeEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent) error
	LogPasswordResetEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent, targetUserID int) error
	LogUserCreateEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent, targetUserID int) error
	LogUserUpdateEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent, targetUserID int) error
	LogUserDeleteEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent, targetUserID int) error
	LogRoleChangeEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent, targetUserID int, newRole string) error

	// Two-factor authentication events
	Log2FAEnableEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent, targetUserID int) error
	Log2FADisableEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent, targetUserID int) error

	// Payout-related events
	LogPayoutRequestEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent, payoutID *int) error
	LogPayoutApprovalEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent, payoutID *int) error
	LogPayoutRejectEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent, payoutID *int) error
	LogPayoutResendEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent, payoutID *int) error
	LogPayoutMarkSentEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent, payoutID *int) error

	// Payin-related events
	LogManualPayinImportEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent, payinID *int) error

	// Other events
	LogMerchantStatusUploadEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent) error
	LogExternalAPIAccessEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent) error
}

type auditLogServiceImpl struct {
	auditLogRepository repository.AuditLogRepository
}

func NewAuditLogService(auditLogRepository repository.AuditLogRepository) AuditLogService {
	return &auditLogServiceImpl{
		auditLogRepository: auditLogRepository,
	}
}

// Core operations implementations

func (s *auditLogServiceImpl) Create(ctx context.Context, auditLog *model.AuditLog) error {
	return s.auditLogRepository.Create(ctx, auditLog)
}

func (s *auditLogServiceImpl) List(ctx context.Context, filter *model.AuditLogFilter) ([]*model.AuditLog, int, int64, error) {
	return s.auditLogRepository.List(ctx, filter)
}

// func (s *auditLogServiceImpl) GetUserIDsWithAuditLogs(ctx context.Context) ([]int, error) {
// 	return s.auditLogRepository.GetUserIDsWithAuditLogs(ctx)
// }

// Generic event logging implementation

func (s *auditLogServiceImpl) LogEventByType(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent, auditLogType object.AuditLogType, description *string) error {
	logDescription := string(auditLogType)
	if description != nil {
		logDescription = *description
	}

	auditLog := model.NewAuditLog(model.NewAuditLogParams{
		UserID:       userID,
		Description:  &logDescription,
		AuditLogType: auditLogType,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
	})

	return s.auditLogRepository.Create(ctx, auditLog)
}

// User-related events implementations

func (s *auditLogServiceImpl) LogLoginEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent) error {
	description := DescLogin
	return s.LogEventByType(ctx, userID, ipAddress, userAgent, object.AuditLogTypeLogin, &description)
}

func (s *auditLogServiceImpl) LogLogoutEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent) error {
	description := DescLogout
	return s.LogEventByType(ctx, userID, ipAddress, userAgent, object.AuditLogTypeLogout, &description)
}

func (s *auditLogServiceImpl) LogPasswordChangeEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent) error {
	description := DescPasswordChange
	return s.LogEventByType(ctx, userID, ipAddress, userAgent, object.AuditLogTypePasswordChange, &description)
}

func (s *auditLogServiceImpl) LogPasswordResetEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent, targetUserID int) error {
	description := fmt.Sprintf(DescPasswordReset, targetUserID)
	return s.LogEventByType(ctx, userID, ipAddress, userAgent, object.AuditLogTypePasswordReset, &description)
}

func (s *auditLogServiceImpl) LogUserCreateEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent, targetUserID int) error {
	description := fmt.Sprintf(DescUserCreate, targetUserID)
	return s.LogEventByType(ctx, userID, ipAddress, userAgent, object.AuditLogTypeUserCreate, &description)
}

func (s *auditLogServiceImpl) LogUserUpdateEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent, targetUserID int) error {
	description := fmt.Sprintf(DescUserUpdate, targetUserID)
	return s.LogEventByType(ctx, userID, ipAddress, userAgent, object.AuditLogTypeUserUpdate, &description)
}

func (s *auditLogServiceImpl) LogUserDeleteEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent, targetUserID int) error {
	description := fmt.Sprintf(DescUserDelete, targetUserID)
	return s.LogEventByType(ctx, userID, ipAddress, userAgent, object.AuditLogTypeUserDelete, &description)
}

func (s *auditLogServiceImpl) LogRoleChangeEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent, targetUserID int, newRole string) error {
	description := fmt.Sprintf(DescRoleChange, targetUserID, newRole)
	return s.LogEventByType(ctx, userID, ipAddress, userAgent, object.AuditLogTypeRoleChange, &description)
}

// Two-factor authentication events implementations

func (s *auditLogServiceImpl) Log2FAEnableEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent, targetUserID int) error {
	description := fmt.Sprintf(Desc2FAEnable, targetUserID)
	return s.LogEventByType(ctx, userID, ipAddress, userAgent, object.AuditLogType2FAEnable, &description)
}

func (s *auditLogServiceImpl) Log2FADisableEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent, targetUserID int) error {
	description := fmt.Sprintf(Desc2FADisable, targetUserID)
	return s.LogEventByType(ctx, userID, ipAddress, userAgent, object.AuditLogType2FADisable, &description)
}

// Payout-related events implementations

func (s *auditLogServiceImpl) LogPayoutRequestEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent, payoutID *int) error {
	description := DescPayoutRequest
	return s.LogEventByTypeWithPayoutID(ctx, userID, ipAddress, userAgent, object.AuditLogTypePayoutRequest, &description, payoutID)
}

func (s *auditLogServiceImpl) LogPayoutApprovalEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent, payoutID *int) error {
	description := DescPayoutApproval
	return s.LogEventByTypeWithPayoutID(ctx, userID, ipAddress, userAgent, object.AuditLogTypePayoutApproval, &description, payoutID)
}

func (s *auditLogServiceImpl) LogPayoutRejectEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent, payoutID *int) error {
	description := DescPayoutReject
	return s.LogEventByTypeWithPayoutID(ctx, userID, ipAddress, userAgent, object.AuditLogTypePayoutReject, &description, payoutID)
}

func (s *auditLogServiceImpl) LogPayoutResendEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent, payoutID *int) error {
	description := DescPayoutResend
	return s.LogEventByTypeWithPayoutID(ctx, userID, ipAddress, userAgent, object.AuditLogTypePayoutResend, &description, payoutID)
}

func (s *auditLogServiceImpl) LogPayoutMarkSentEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent, payoutID *int) error {
	description := DescPayoutMarkSent
	return s.LogEventByTypeWithPayoutID(ctx, userID, ipAddress, userAgent, object.AuditLogTypePayoutMarkSent, &description, payoutID)
}

// Helper method for payout-related events
func (s *auditLogServiceImpl) LogEventByTypeWithPayoutID(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent, auditLogType object.AuditLogType, description *string, payoutID *int) error {
	logDescription := string(auditLogType)
	if description != nil {
		logDescription = *description
	}

	auditLog := model.NewAuditLog(model.NewAuditLogParams{
		UserID:       userID,
		Description:  &logDescription,
		AuditLogType: auditLogType,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		PayoutID:     payoutID,
	})

	return s.auditLogRepository.Create(ctx, auditLog)
}

// Payin-related events implementations

func (s *auditLogServiceImpl) LogManualPayinImportEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent, payinID *int) error {
	description := DescManualPayinImport
	return s.LogEventByTypeWithPayinID(ctx, userID, ipAddress, userAgent, object.AuditLogTypeManualPayinImport, &description, payinID)
}

// Helper method for payin-related events
func (s *auditLogServiceImpl) LogEventByTypeWithPayinID(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent, auditLogType object.AuditLogType, description *string, payinID *int) error {
	logDescription := string(auditLogType)
	if description != nil {
		logDescription = *description
	}

	auditLog := model.NewAuditLog(model.NewAuditLogParams{
		UserID:       userID,
		Description:  &logDescription,
		AuditLogType: auditLogType,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		PayinID:      payinID,
	})

	return s.auditLogRepository.Create(ctx, auditLog)
}

// Other events implementations

func (s *auditLogServiceImpl) LogMerchantStatusUploadEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent) error {
	description := DescMerchantStatusUpload
	return s.LogEventByType(ctx, userID, ipAddress, userAgent, object.AuditLogTypeMerchantStatusUpload, &description)
}

func (s *auditLogServiceImpl) LogExternalAPIAccessEvent(ctx context.Context, userID *int, ipAddress *object.IPAddress, userAgent *object.UserAgent) error {
	description := DescExternalAPIAccess
	return s.LogEventByType(ctx, userID, ipAddress, userAgent, object.AuditLogTypeExternalAPIAccess, &description)
}
