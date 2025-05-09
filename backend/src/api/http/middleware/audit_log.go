package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/usecase"
)

const (
	TargetUserIDKey contextKey = "target_user_id"
	NewRoleKey      contextKey = "new_role"
	PayoutIDKey     contextKey = "payout_id"
	PayinIDKey      contextKey = "payin_id"
)

type AuditLogOptions struct {
	AuditLogType string
}

func NewAuditLogger(auditLogUsecase *usecase.AuditLogUsecase) *AuditLogBuilder {
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

func (a *AuditLogBuilder) AsMiddleware() MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return a.createMiddleware(next)
	}
}

func (a *AuditLogBuilder) AsResponseMiddleware() MiddlewareFunc {
	return CreateResponseMiddleware(func(w http.ResponseWriter, r *http.Request) {
		a.logAuditEvent(w, r)
	})
}

func (a *AuditLogBuilder) createMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alw := &auditLogResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next(alw, r)

		if alw.statusCode >= 200 && alw.statusCode < 300 {
			a.logAuditEvent(alw, r)
		}
	}
}

func (a *AuditLogBuilder) logAuditEvent(_ http.ResponseWriter, r *http.Request) {
	userIDVal := r.Context().Value(UserIDKey)

	ipAddress := getIPAddress(r)
	userAgent := r.UserAgent()

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

	if targetIDVal := r.Context().Value(TargetUserIDKey); targetIDVal != nil {
		switch v := targetIDVal.(type) {
		case int:
			targetUserID = v
		case string:
			targetUserID = parseInt(v)
		}
	}

	if roleVal := r.Context().Value(NewRoleKey); roleVal != nil {
		if role, ok := roleVal.(string); ok {
			newRole = role
		}
	}

	if poutID := r.Context().Value(PayoutIDKey); poutID != nil {
		if id, ok := poutID.(*int); ok {
			payoutID = id
		}
	}

	if pinID := r.Context().Value(PayinIDKey); pinID != nil {
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

	switch a.options.AuditLogType {
	case models.AuditLogTypeLogin:
		err = a.auditLogUsecase.LogLoginEvent(r.Context(), &targetUserID, ipAddressPtr, userAgentPtr)
	case models.AuditLogTypeLogout:
		err = a.auditLogUsecase.LogLogoutEvent(r.Context(), userIDInt, ipAddressPtr, userAgentPtr)
	case models.AuditLogTypePasswordChange:
		err = a.auditLogUsecase.LogPasswordChangeEvent(r.Context(), userIDInt, ipAddressPtr, userAgentPtr)
	case models.AuditLogTypePasswordReset:
		err = a.auditLogUsecase.LogPasswordResetEvent(r.Context(), userIDInt, ipAddressPtr, userAgentPtr, targetUserID)
	case models.AuditLogTypeUserCreate:
		err = a.auditLogUsecase.LogUserCreateEvent(r.Context(), userIDInt, ipAddressPtr, userAgentPtr, targetUserID)
	case models.AuditLogTypeUserUpdate:
		err = a.auditLogUsecase.LogUserUpdateEvent(r.Context(), userIDInt, ipAddressPtr, userAgentPtr, targetUserID)
	case models.AuditLogTypeUserDelete:
		err = a.auditLogUsecase.LogUserDeleteEvent(r.Context(), userIDInt, ipAddressPtr, userAgentPtr, targetUserID)
	case models.AuditLogTypeRoleChange:
		err = a.auditLogUsecase.LogRoleChangeEvent(r.Context(), userIDInt, ipAddressPtr, userAgentPtr, targetUserID, newRole)
	case models.AuditLogTypePayoutRequest:
		err = a.auditLogUsecase.LogPayoutRequestEvent(r.Context(), userIDInt, ipAddressPtr, userAgentPtr, payoutID)
	case models.AuditLogTypePayoutApproval:
		err = a.auditLogUsecase.LogPayoutApprovalEvent(r.Context(), userIDInt, ipAddressPtr, userAgentPtr, payoutID)
	case models.AuditLogTypePayoutReject:
		err = a.auditLogUsecase.LogPayoutRejectEvent(r.Context(), userIDInt, ipAddressPtr, userAgentPtr, payoutID)
	case models.AuditLogType2FAEnable:
		err = a.auditLogUsecase.Log2FAEnableEvent(r.Context(), userIDInt, ipAddressPtr, userAgentPtr, targetUserID)
	case models.AuditLogType2FADisable:
		err = a.auditLogUsecase.Log2FADisableEvent(r.Context(), userIDInt, ipAddressPtr, userAgentPtr, targetUserID)
	case models.AuditLogTypeManualPayinImport:
		err = a.auditLogUsecase.LogManualPayinImportEvent(r.Context(), userIDInt, ipAddressPtr, userAgentPtr, payinID)
	case models.AuditLogTypePayoutResend:
		err = a.auditLogUsecase.LogPayoutResendEvent(r.Context(), userIDInt, ipAddressPtr, userAgentPtr, payoutID)
	case models.AuditLogTypePayoutMarkSent:
		err = a.auditLogUsecase.LogPayoutMarkSentEvent(r.Context(), userIDInt, ipAddressPtr, userAgentPtr, payoutID)
	case models.AuditLogTypeMerchantStatusUpload:
		err = a.auditLogUsecase.LogMerchantStatusUploadEvent(r.Context(), userIDInt, ipAddressPtr, userAgentPtr)
	case models.AuditLogTypeExternalAPIAccess:
		err = a.auditLogUsecase.LogExternalAPIAccessEvent(r.Context(), userIDInt, ipAddressPtr, userAgentPtr)
	default:
		err = a.auditLogUsecase.LogEventByType(r.Context(), userIDInt, ipAddressPtr, userAgentPtr, a.options.AuditLogType, nil)
	}

	if err != nil {
		fmt.Printf("Failed to log audit event: %v\n", err)
	}
}

type auditLogResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *auditLogResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
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
