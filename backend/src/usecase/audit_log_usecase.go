package usecase

import (
	"context"

	"github.com/vnlab/makeshop-payment/src/domain/models"
	"github.com/vnlab/makeshop-payment/src/domain/repositories"
	"github.com/vnlab/makeshop-payment/src/domain/repositories/filter"
	"github.com/vnlab/makeshop-payment/src/lib/i18n"
)

// AuditLogUsecase handles audit logging business logic
type AuditLogUsecase struct {
	auditLogRepo     repositories.AuditLogRepository
	auditLogTypeRepo repositories.AuditLogTypeRepository
	userRepo         repositories.UserRepository
}

// NewAuditLogUsecase creates a new AuditLogUsecase
func NewAuditLogUsecase(
	auditLogRepo repositories.AuditLogRepository,
	auditLogTypeRepo repositories.AuditLogTypeRepository,
	userRepo repositories.UserRepository,
) *AuditLogUsecase {
	return &AuditLogUsecase{
		auditLogRepo:     auditLogRepo,
		auditLogTypeRepo: auditLogTypeRepo,
		userRepo:         userRepo,
	}
}

// LogEventByType logs an event with the specified audit log type and custom description
func (uc *AuditLogUsecase) LogEventByType(ctx context.Context, userID *int, ipAddress *string, userAgent *string, auditLogType string, description *string) error {
	logDescription := auditLogType
	if description != nil {
		logDescription = *description
	}

	auditLog := &models.AuditLog{
		UserID:       userID,
		Description:  &logDescription,
		AuditLogType: auditLogType,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
	}

	return uc.auditLogRepo.Create(ctx, auditLog)
}

func (uc *AuditLogUsecase) LogLoginEvent(ctx context.Context, userID *int, ipAddress *string, userAgent *string) error {
	description := i18n.T(ctx, "audit_log.login")
	return uc.LogEventByType(ctx, userID, ipAddress, userAgent, models.AuditLogTypeLogin, &description)
}

func (uc *AuditLogUsecase) LogLogoutEvent(ctx context.Context, userID *int, ipAddress *string, userAgent *string) error {
	description := i18n.T(ctx, "audit_log.logout")
	return uc.LogEventByType(ctx, userID, ipAddress, userAgent, models.AuditLogTypeLogout, &description)
}

func (uc *AuditLogUsecase) LogPasswordChangeEvent(ctx context.Context, userID *int, ipAddress *string, userAgent *string) error {
	description := i18n.T(ctx, "audit_log.password_change")
	return uc.LogEventByType(ctx, userID, ipAddress, userAgent, models.AuditLogTypePasswordChange, &description)
}

func (uc *AuditLogUsecase) LogPasswordResetEvent(ctx context.Context, userID *int, ipAddress *string, userAgent *string, targetUserID int) error {
	description := i18n.T(ctx, "audit_log.password_reset", targetUserID)
	return uc.LogEventByType(ctx, userID, ipAddress, userAgent, models.AuditLogTypePasswordReset, &description)
}

func (uc *AuditLogUsecase) LogUserCreateEvent(ctx context.Context, userID *int, ipAddress *string, userAgent *string, targetUserID int) error {
	description := i18n.T(ctx, "audit_log.user_create", targetUserID)
	return uc.LogEventByType(ctx, userID, ipAddress, userAgent, models.AuditLogTypeUserCreate, &description)
}

func (uc *AuditLogUsecase) LogUserUpdateEvent(ctx context.Context, userID *int, ipAddress *string, userAgent *string, targetUserID int) error {
	description := i18n.T(ctx, "audit_log.user_update", targetUserID)
	return uc.LogEventByType(ctx, userID, ipAddress, userAgent, models.AuditLogTypeUserUpdate, &description)
}

// LogUserDeleteEvent logs a user deletion event
func (uc *AuditLogUsecase) LogUserDeleteEvent(ctx context.Context, userID *int, ipAddress *string, userAgent *string, targetUserID int) error {
	description := i18n.T(ctx, "audit_log.user_delete", targetUserID)
	return uc.LogEventByType(ctx, userID, ipAddress, userAgent, models.AuditLogTypeUserDelete, &description)
}

// LogRoleChangeEvent logs a role change event
func (uc *AuditLogUsecase) LogRoleChangeEvent(ctx context.Context, userID *int, ipAddress *string, userAgent *string, targetUserID int, newRole string) error {
	description := i18n.T(ctx, "audit_log.role_change", targetUserID, newRole)
	return uc.LogEventByType(ctx, userID, ipAddress, userAgent, models.AuditLogTypeRoleChange, &description)
}

func (uc *AuditLogUsecase) LogPayoutRequestEvent(ctx context.Context, userID *int, ipAddress *string, userAgent *string, payoutID *int) error {
	description := i18n.T(ctx, "audit_log.payout_request")

	auditLog := &models.AuditLog{
		UserID:       userID,
		Description:  &description,
		AuditLogType: models.AuditLogTypePayoutRequest,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		PayoutID:     payoutID,
	}

	return uc.auditLogRepo.Create(ctx, auditLog)
}

func (uc *AuditLogUsecase) LogPayoutApprovalEvent(ctx context.Context, userID *int, ipAddress *string, userAgent *string, payoutID *int) error {
	description := i18n.T(ctx, "audit_log.payout_approval")

	auditLog := &models.AuditLog{
		UserID:       userID,
		Description:  &description,
		AuditLogType: models.AuditLogTypePayoutApproval,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		PayoutID:     payoutID,
	}

	return uc.auditLogRepo.Create(ctx, auditLog)
}

func (uc *AuditLogUsecase) LogPayoutRejectEvent(ctx context.Context, userID *int, ipAddress *string, userAgent *string, payoutID *int) error {
	description := i18n.T(ctx, "audit_log.payout_reject")

	auditLog := &models.AuditLog{
		UserID:       userID,
		Description:  &description,
		AuditLogType: models.AuditLogTypePayoutReject,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		PayoutID:     payoutID,
	}

	return uc.auditLogRepo.Create(ctx, auditLog)
}

func (uc *AuditLogUsecase) Log2FAEnableEvent(ctx context.Context, userID *int, ipAddress *string, userAgent *string, targetUserID int) error {
	description := i18n.T(ctx, "audit_log.2fa_enable", targetUserID)
	return uc.LogEventByType(ctx, userID, ipAddress, userAgent, models.AuditLogType2FAEnable, &description)
}

func (uc *AuditLogUsecase) Log2FADisableEvent(ctx context.Context, userID *int, ipAddress *string, userAgent *string, targetUserID int) error {
	description := i18n.T(ctx, "audit_log.2fa_disable", targetUserID)
	return uc.LogEventByType(ctx, userID, ipAddress, userAgent, models.AuditLogType2FADisable, &description)
}

func (uc *AuditLogUsecase) LogManualPayinImportEvent(ctx context.Context, userID *int, ipAddress *string, userAgent *string, payinID *int) error {
	description := i18n.T(ctx, "audit_log.manual_payin_import")

	auditLog := &models.AuditLog{
		UserID:       userID,
		Description:  &description,
		AuditLogType: models.AuditLogTypeManualPayinImport,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		PayinID:      payinID,
	}

	return uc.auditLogRepo.Create(ctx, auditLog)
}

func (uc *AuditLogUsecase) LogPayoutResendEvent(ctx context.Context, userID *int, ipAddress *string, userAgent *string, payoutID *int) error {
	description := i18n.T(ctx, "audit_log.payout_resend")

	auditLog := &models.AuditLog{
		UserID:       userID,
		Description:  &description,
		AuditLogType: models.AuditLogTypePayoutResend,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		PayoutID:     payoutID,
	}

	return uc.auditLogRepo.Create(ctx, auditLog)
}

func (uc *AuditLogUsecase) LogPayoutMarkSentEvent(ctx context.Context, userID *int, ipAddress *string, userAgent *string, payoutID *int) error {
	description := i18n.T(ctx, "audit_log.payout_mark_sent")

	auditLog := &models.AuditLog{
		UserID:       userID,
		Description:  &description,
		AuditLogType: models.AuditLogTypePayoutMarkSent,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		PayoutID:     payoutID,
	}

	return uc.auditLogRepo.Create(ctx, auditLog)
}

func (uc *AuditLogUsecase) LogMerchantStatusUploadEvent(ctx context.Context, userID *int, ipAddress *string, userAgent *string) error {
	description := i18n.T(ctx, "audit_log.merchant_status_upload")
	return uc.LogEventByType(ctx, userID, ipAddress, userAgent, models.AuditLogTypeMerchantStatusUpload, &description)
}

func (uc *AuditLogUsecase) LogExternalAPIAccessEvent(ctx context.Context, userID *int, ipAddress *string, userAgent *string) error {
	description := i18n.T(ctx, "audit_log.external_api_access")
	return uc.LogEventByType(ctx, userID, ipAddress, userAgent, models.AuditLogTypeExternalAPIAccess, &description)
}

func (uc *AuditLogUsecase) ListAuditLogs(ctx context.Context, filter *filter.AuditLogFilter) ([]*models.AuditLog, int, int64, error) {
	return uc.auditLogRepo.List(ctx, filter)
}

func (uc *AuditLogUsecase) GetAuditLogUsers(ctx context.Context) ([]*models.User, error) {
	return uc.userRepo.GetUsersWithAuditLogs(ctx)
}
