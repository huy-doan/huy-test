package handlers

import (
	"net/http"

	"github.com/vnlab/makeshop-payment/src/api/http/errors"
	"github.com/vnlab/makeshop-payment/src/api/http/middleware"
	"github.com/vnlab/makeshop-payment/src/api/http/response"
	"github.com/vnlab/makeshop-payment/src/api/http/serializers"
	models "github.com/vnlab/makeshop-payment/src/domain/models"
	"github.com/vnlab/makeshop-payment/src/domain/repositories/filter"
	"github.com/vnlab/makeshop-payment/src/lib/i18n"
	"github.com/vnlab/makeshop-payment/src/lib/utils"
	"github.com/vnlab/makeshop-payment/src/usecase"
)

// AuditLogHandler handles audit log-related requests
type AuditLogHandler struct {
	auditLogUsecase *usecase.AuditLogUsecase
}

// NewAuditLogHandler creates a new AuditLogHandler
func NewAuditLogHandler(auditLogUsecase *usecase.AuditLogUsecase) *AuditLogHandler {
	return &AuditLogHandler{
		auditLogUsecase: auditLogUsecase,
	}
}

// ListAuditLogs handles the request to list audit logs
// @Summary List all audit logs
// @Description Get a list of all audit logs in the system with optional filtering
// @Tags Admin Audit Log Management
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Param sort_field query string false "Field to sort by (default: id)"
// @Param sort_order query string false "Sort order: ascend or descend (default: ascend)"
// @Param created_at query string false "Filter by created date (RFC3339, RFC1123, or YYYY-MM-DD format)"
// @Param user_id query int false "Filter by user ID"
// @Param keyword query string false "Filter by keyword in description"
// @Param audit_log_type query string false "Filter by audit log type"
// @Security BearerAuth
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /admin/audit-logs [get]
func (h *AuditLogHandler) ListAuditLogs(w http.ResponseWriter, r *http.Request) {
	// Check admin role
	roleCode, ok := r.Context().Value(middleware.RoleCodeKey).(string)
	if !ok || roleCode != string(models.RoleCodeAdmin) {
		response.Forbidden(w, i18n.T(r.Context(), "account.unauthorized"))
		return
	}

	auditLogFilter := filter.NewAuditLogFilter()

	baseFilter := utils.ExtractPaginationAndSorting(r)

	auditLogFilter.Pagination = baseFilter.Pagination
	auditLogFilter.Sort = baseFilter.Sort

	err := applyAuditLogFilterParams(r, auditLogFilter)
	if err != nil {
		message := i18n.T(r.Context(), "audit_log.invalid_filter_params")
		response.BadRequest(w, message, nil)
		return
	}

	auditLogFilter.ApplyFilters()

	auditLogs, totalPages, total, err := h.auditLogUsecase.ListAuditLogs(r.Context(), auditLogFilter)
	if err != nil {
		message := i18n.T(r.Context(), "audit_log.list_failed")
		response.Error(w, errors.InternalError(message))
		return
	}

	responseData := map[string]any{
		"audit_logs":  serializers.SerializeAuditLogCollection(auditLogs),
		"page":        auditLogFilter.Pagination.Page,
		"page_size":   auditLogFilter.Pagination.PageSize,
		"total_pages": totalPages,
		"total":       total,
	}

	response.Success(w, responseData, i18n.T(r.Context(), "audit_log.list_success"))
}

// GetAuditLogUsers returns a list of distinct users who have audit log entries
// @Summary Get distinct users with audit logs
// @Description Retrieve a list of distinct users who have audit log entries
// @Tags Admin Audit Log Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response "Success"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /admin/audit-logs/users [get]
func (h *AuditLogHandler) GetAuditLogUsers(w http.ResponseWriter, r *http.Request) {
	// Check admin role
	roleCode, ok := r.Context().Value(middleware.RoleCodeKey).(string)
	if !ok || roleCode != string(models.RoleCodeAdmin) {
		response.Forbidden(w, i18n.T(r.Context(), "account.unauthorized"))
		return
	}

	users, err := h.auditLogUsecase.GetAuditLogUsers(r.Context())
	if err != nil {
		message := i18n.T(r.Context(), "audit_log.get_users_failed")
		response.Error(w, errors.InternalError(message))
		return
	}

	responseData := map[string]any{
		"users": serializers.SerializeUserCollection(users),
	}

	response.Success(w, responseData, i18n.T(r.Context(), "audit_log.get_users_success"))
}

func applyAuditLogFilterParams(r *http.Request, auditLogFilter *filter.AuditLogFilter) error {
	createdAt, err := utils.ExtractDateParam(r, "created_at")
	if err != nil {
		return err
	}
	auditLogFilter.CreatedAt = createdAt

	userID, err := utils.ExtractIntParam(r, "user_id")
	if err != nil {
		return err
	}
	auditLogFilter.UserID = userID

	auditLogType := r.URL.Query().Get("audit_log_type")
	if auditLogType != "" {
		auditLogFilter.AuditLogType = &auditLogType
	}

	description := r.URL.Query().Get("description")
	if description != "" {
		auditLogFilter.Description = &description
	}

	return nil
}
