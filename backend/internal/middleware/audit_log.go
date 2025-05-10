package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	object "github.com/huydq/test/internal/domain/object/audit_log"
	"github.com/huydq/test/internal/domain/service"
	"github.com/labstack/echo/v4"
)

const (
	ContextKey_AuditLogUserID       ContextKey = "user_id"
	ContextKey_AuditLogTargetUserID ContextKey = "target_user_id"
	ContextKey_AuditLogNewRole      ContextKey = "new_role"
	ContextKey_AuditLogPayoutID     ContextKey = "payout_id"
	ContextKey_AuditLogPayinID      ContextKey = "payin_id"
)

type AuditLogOptions struct {
	AuditLogType object.AuditLogType
}

func (m *MiddlewareManager) NewAuditLogger(auditLogService service.AuditLogService) *AuditLogBuilder {
	return &AuditLogBuilder{
		auditLogService: auditLogService,
	}
}

type AuditLogBuilder struct {
	auditLogService service.AuditLogService
	options         AuditLogOptions
}

func (a *AuditLogBuilder) WithType(auditLogType object.AuditLogType) *AuditLogBuilder {
	newBuilder := &AuditLogBuilder{
		auditLogService: a.auditLogService,
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

	ipAddressStr := getIPAddress(c.Request())
	userAgentStr := c.Request().UserAgent()

	ipAddress := object.IPAddress(ipAddressStr)
	userAgent := object.UserAgent(userAgentStr)

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
	case object.AuditLogTypeLogin:
		err = a.auditLogService.LogLoginEvent(ctx, &targetUserID, ipAddressPtr, userAgentPtr)
	case object.AuditLogTypeLogout:
		err = a.auditLogService.LogLogoutEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr)
	case object.AuditLogTypePasswordChange:
		err = a.auditLogService.LogPasswordChangeEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr)
	case object.AuditLogTypePasswordReset:
		err = a.auditLogService.LogPasswordResetEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr, targetUserID)
	case object.AuditLogTypeUserCreate:
		err = a.auditLogService.LogUserCreateEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr, targetUserID)
	case object.AuditLogTypeUserUpdate:
		err = a.auditLogService.LogUserUpdateEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr, targetUserID)
	case object.AuditLogTypeUserDelete:
		err = a.auditLogService.LogUserDeleteEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr, targetUserID)
	case object.AuditLogTypeRoleChange:
		err = a.auditLogService.LogRoleChangeEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr, targetUserID, newRole)
	case object.AuditLogTypePayoutRequest:
		err = a.auditLogService.LogPayoutRequestEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr, payoutID)
	case object.AuditLogTypePayoutApproval:
		err = a.auditLogService.LogPayoutApprovalEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr, payoutID)
	case object.AuditLogTypePayoutReject:
		err = a.auditLogService.LogPayoutRejectEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr, payoutID)
	case object.AuditLogType2FAEnable:
		err = a.auditLogService.Log2FAEnableEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr, targetUserID)
	case object.AuditLogType2FADisable:
		err = a.auditLogService.Log2FADisableEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr, targetUserID)
	case object.AuditLogTypeManualPayinImport:
		err = a.auditLogService.LogManualPayinImportEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr, payinID)
	case object.AuditLogTypePayoutResend:
		err = a.auditLogService.LogPayoutResendEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr, payoutID)
	case object.AuditLogTypePayoutMarkSent:
		err = a.auditLogService.LogPayoutMarkSentEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr, payoutID)
	case object.AuditLogTypeMerchantStatusUpload:
		err = a.auditLogService.LogMerchantStatusUploadEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr)
	case object.AuditLogTypeExternalAPIAccess:
		err = a.auditLogService.LogExternalAPIAccessEvent(ctx, userIDInt, ipAddressPtr, userAgentPtr)
	default:
		err = a.auditLogService.LogEventByType(ctx, userIDInt, ipAddressPtr, userAgentPtr, a.options.AuditLogType, nil)
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

func hasEnoughDataForAuditLog(auditLogType object.AuditLogType, userIDInt *int, targetUserID int, newRole string, payoutID *int, payinID *int) bool {
	switch auditLogType {
	case object.AuditLogTypeLogin:
		return targetUserID != 0
	case object.AuditLogTypeLogout, object.AuditLogTypePasswordChange, object.AuditLogTypePasswordReset:
		return userIDInt != nil
	case object.AuditLogTypeUserCreate, object.AuditLogTypeUserUpdate, object.AuditLogTypeUserDelete, object.AuditLogTypeRoleChange:
		return userIDInt != nil && targetUserID != 0
	case object.AuditLogTypePayoutRequest, object.AuditLogTypePayoutApproval, object.AuditLogTypePayoutReject, object.AuditLogTypePayoutResend, object.AuditLogTypePayoutMarkSent:
		return userIDInt != nil && payoutID != nil
	case object.AuditLogTypeManualPayinImport:
		return userIDInt != nil && payinID != nil
	case object.AuditLogType2FAEnable, object.AuditLogType2FADisable:
		return userIDInt != nil && targetUserID != 0
	default:
		return userIDInt != nil
	}
}
