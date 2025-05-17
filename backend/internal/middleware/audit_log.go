package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	model "github.com/huydq/test/internal/domain/model/audit_log"
	object "github.com/huydq/test/internal/domain/object/audit_log"
	"github.com/huydq/test/internal/domain/service"
	"github.com/labstack/echo/v4"
)

const (
	ContextKey_AuditLogTargetUserID ContextKey = "targetUserId"
	ContextKey_AuditLogNewRole      ContextKey = "newRole"
	ContextKey_AuditLogPayoutID     ContextKey = "payoutId"
	ContextKey_AuditLogPayinID      ContextKey = "payinId"
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
	userIDVal := c.Get(string(ContextKey_AuthUserIDKey))
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

	ctx := c.Request().Context()

	auditLogModel := model.NewAuditLogWithData(
		userIDInt,
		a.options.AuditLogType,
		ipAddressPtr,
		userAgentPtr,
		targetUserID,
		newRole,
		payoutID,
		payinID,
	)

	err := a.auditLogService.CreateAuditLog(ctx, auditLogModel)
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
