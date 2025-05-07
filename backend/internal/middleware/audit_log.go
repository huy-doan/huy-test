package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/vnlab/makeshop-payment/src/domain/models"
	"github.com/vnlab/makeshop-payment/src/usecase"
)

const (
	ContextKey_AuditLogUserID       ContextKey = "user_id"
	ContextKey_AuditLogTargetUserID ContextKey = "target_user_id"
	ContextKey_AuditLogNewRole      ContextKey = "new_role"
	ContextKey_AuditLogPayoutID     ContextKey = "payout_id"
	ContextKey_AuditLogPayinID      ContextKey = "payin_id"
)

type AuditLogOptions struct {
	AuditLogType string
}

func (m *MiddlewareManager) NewAuditLogger(auditLogUsecase *usecase.AuditLogUsecase) *AuditLogBuilder {
	return &AuditLogBuilder{
		auditLogUsecase: auditLogUsecase,
	}
}

type AuditLogBuilder struct {
	auditLogUsecase *usecase.AuditLogUsecase
	options         AuditLogOptions
}

func (a *AuditLogBuilder) WithType(auditLogType string) *AuditLogBuilder {
	newBuilder := &AuditLogBuilder{
		auditLogUsecase: a.auditLogUsecase,
		options: AuditLogOptions{
			AuditLogType: auditLogType,
		},
	}
	return newBuilder
}

// AsMiddleware returns an Echo middleware function
func (a *AuditLogBuilder) AsMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)

			// Check if the response was successful
			if c.Response().Status >= 200 && c.Response().Status < 300 {
				a.logAuditEvent(c)
			}

			return err
		}
	}
}

// AsResponseMiddleware returns a middleware that only logs the audit event
func (a *AuditLogBuilder) AsResponseMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			a.logAuditEvent(c)
			return err
		}
	}
}

func (a *AuditLogBuilder) logAuditEvent(c echo.Context) {
	userIDVal := c.Get(string(ContextKey_AuditLogUserID))

	ipAddress := getIPAddress(c.Request())
	userAgent := c.Request().UserAgent()

	var userIDInt *int
	switch v := userIDVal.(type) {
	case int:
		id := v
		userIDInt = &id
	case string:
		id := parseInt(v)
		userIDInt = &id
	default:
		log.Printf("Unexpected user ID type: %T", userIDVal)
	}

	var targetUserID int
	var newRole string
	var payoutID *int
	var payinID *int

	if targetIDVal := c.Get(string(ContextKey_AuditLogTargetUserID)); targetIDVal != nil {
		switch v := targetIDVal.(type) {
		case int:
			targetUserID = v
		case string:
			targetUserID = parseInt(v)
		}
	}

	if roleVal := c.Get(string(ContextKey_AuditLogNewRole)); roleVal != nil {
		if role, ok := roleVal.(string); ok {
			newRole = role
		}
	}

	if poutID := c.Get(string(ContextKey_AuditLogPayoutID)); poutID != nil {
		if id, ok := poutID.(*int); ok {
			payoutID = id
		}
	}

	if pinID := c.Get(string(ContextKey_AuditLogPayinID)); pinID != nil {
		if id, ok := pinID.(*int); ok {
			payinID = id
		}
	}

	ipAddressPtr := &ipAddress
	userAgentPtr := &userAgent

	if !hasEnoughDataForAuditLog(a.options.AuditLogType, userIDInt, targetUserID, newRole, payoutID, payinID) {
		log.Printf("Skipping audit logging due to insufficient data for event type: %s", a.options.AuditLogType)
		return
	}

	var err error
	ctx := c.Request().Context()

	switch a.options.AuditLogType {
	case models.AuditLogTypeLogin:
		err = a.auditLogUsecase.LogLoginEvent(ctx, &targetUserID, ipAddressPtr, userAgentPtr)
	case models.AuditLogTypeLogout:
		err = a.auditLogUsecase.LogLogoutEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr)
	case models.AuditLogTypePasswordChange:
		err = a.auditLogUsecase.LogPasswordChangeEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr)
	case models.AuditLogTypePasswordReset:
		err = a.auditLogUsecase.LogPasswordResetEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr, targetUserID)
	case models.AuditLogTypeUserCreate:
		err = a.auditLogUsecase.LogUserCreateEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr, targetUserID)
	case models.AuditLogTypeUserUpdate:
		err = a.auditLogUsecase.LogUserUpdateEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr, targetUserID)
	case models.AuditLogTypeUserDelete:
		err = a.auditLogUsecase.LogUserDeleteEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr, targetUserID)
	case models.AuditLogTypeRoleChange:
		err = a.auditLogUsecase.LogRoleChangeEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr, targetUserID, newRole)
	case models.AuditLogTypePayoutRequest:
		err = a.auditLogUsecase.LogPayoutRequestEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr, payoutID)
	case models.AuditLogTypePayoutApproval:
		err = a.auditLogUsecase.LogPayoutApprovalEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr, payoutID)
	case models.AuditLogTypePayoutReject:
		err = a.auditLogUsecase.LogPayoutRejectEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr, payoutID)
	case models.AuditLogType2FAEnable:
		err = a.auditLogUsecase.Log2FAEnableEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr, targetUserID)
	case models.AuditLogType2FADisable:
		err = a.auditLogUsecase.Log2FADisableEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr, targetUserID)
	case models.AuditLogTypeManualPayinImport:
		err = a.auditLogUsecase.LogManualPayinImportEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr, payinID)
	case models.AuditLogTypePayoutResend:
		err = a.auditLogUsecase.LogPayoutResendEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr, payoutID)
	case models.AuditLogTypePayoutMarkSent:
		err = a.auditLogUsecase.LogPayoutMarkSentEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr, payoutID)
	case models.AuditLogTypeMerchantStatusUpload:
		err = a.auditLogUsecase.LogMerchantStatusUploadEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr)
	case models.AuditLogTypeExternalAPIAccess:
		err = a.auditLogUsecase.LogExternalAPIAccessEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr)
	default:
		err = a.auditLogUsecase.LogEventByType(ctx, userIDInt, ipAddressPtr, userAgentPtr, a.options.AuditLogType, nil)
	}

	if err != nil {
		fmt.Printf("Failed to log audit event: %v\n", err)
	}
}

func getIPAddress(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return strings.Split(ip, ",")[0]
	}

	return r.RemoteAddr
}

func parseInt(s string) int {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	if err != nil {
		return 0
	}
	return i
}

func hasEnoughDataForAuditLog(auditLogType string, userIDInt *int, targetUserID int, newRole string, payoutID *int, payinID *int) bool {
	switch auditLogType {
	case models.AuditLogTypeLogin, models.AuditLogTypeLogout, models.AuditLogTypePasswordChange, models.AuditLogTypePasswordReset:
		return userIDInt != nil
	case models.AuditLogTypeUserCreate, models.AuditLogTypeUserUpdate, models.AuditLogTypeUserDelete, models.AuditLogTypeRoleChange:
		return userIDInt != nil && targetUserID != 0
	case models.AuditLogTypePayoutRequest, models.AuditLogTypePayoutApproval, models.AuditLogTypePayoutReject, models.AuditLogTypePayoutResend, models.AuditLogTypePayoutMarkSent:
		return userIDInt != nil && payoutID != nil
	case models.AuditLogTypeManualPayinImport:
		return userIDInt != nil && payinID != nil
	case models.AuditLogType2FAEnable, models.AuditLogType2FADisable:
		return userIDInt != nil && targetUserID != 0
	default:
		return userIDInt != nil
	}
}
