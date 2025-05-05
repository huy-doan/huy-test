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

// PayoutHandler handles payout-related requests
type PayoutHandler struct {
	payoutUsecase *usecase.PayoutUsecase
}

// NewPayoutHandler creates a new PayoutHandler
func NewPayoutHandler(payoutUsecase *usecase.PayoutUsecase) *PayoutHandler {
	return &PayoutHandler{
		payoutUsecase: payoutUsecase,
	}
}

// ListPayouts handles the request to list payouts
// @Summary List all payouts
// @Description Get a list of all payouts in the system with optional filtering
// @Tags Admin Payout Management
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Param sort_field query string false "Field to sort by (default: id)"
// @Param sort_order query string false "Sort order: ascend or descend (default: ascend)"
// @Param created_at query string false "Filter by created date (RFC3339, RFC1123, or YYYY-MM-DD format)"
// @Param sending_date query string false "Filter by sending date (RFC3339, RFC1123, or YYYY-MM-DD format)"
// @Param sent_date query string false "Filter by sent date (RFC3339, RFC1123, or YYYY-MM-DD format)"
// @Param payout_status query int false "Filter by payout status"
// @Security BearerAuth
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /admin/payouts [get]
func (h *PayoutHandler) ListPayouts(w http.ResponseWriter, r *http.Request) {
	// Check admin role
	roleCode, ok := r.Context().Value(middleware.RoleCodeKey).(string)
	if !ok || roleCode != string(models.RoleCodeAdmin) {
		response.Forbidden(w, i18n.T(r.Context(), "account.unauthorized"))
		return
	}
	payoutFilter := filter.NewPayoutFilter()

	baseFilter := utils.ExtractPaginationAndSorting(r)

	payoutFilter.Pagination = baseFilter.Pagination
	payoutFilter.Sort = baseFilter.Sort

	err := applyPayoutFilterParams(r, payoutFilter)
	if err != nil {
		message := i18n.T(r.Context(), "payout.invalid_filter_params")
		response.BadRequest(w, message, nil)
		return
	}

	payoutFilter.ApplyFilters()

	payouts, totalPages, total, err := h.payoutUsecase.GetAll(r.Context(), payoutFilter)
	if err != nil {
		message := i18n.T(r.Context(), "payout.list_failed")
		response.Error(w, errors.InternalError(message))
		return
	}

	responseData := map[string]any{
		"payouts":     serializers.SerializePayoutCollection(payouts),
		"page":        payoutFilter.Pagination.Page,
		"page_size":   payoutFilter.Pagination.PageSize,
		"total_pages": totalPages,
		"total":       total,
	}

	response.Success(w, responseData, i18n.T(r.Context(), "payout.list_success"))
}

func applyPayoutFilterParams(r *http.Request, payoutFilter *filter.PayoutFilter) error {
	createdAt, err := utils.ExtractDateParam(r, "created_at")
	if err != nil {
		return err
	}
	payoutFilter.CreatedAt = createdAt

	sendingDate, err := utils.ExtractDateParam(r, "sending_date")
	if err != nil {
		return err
	}
	payoutFilter.SendingDate = sendingDate

	sentDate, err := utils.ExtractDateParam(r, "sent_date")
	if err != nil {
		return err
	}
	payoutFilter.SentDate = sentDate

	payoutStatus, err := utils.ExtractIntParam(r, "payout_status")
	if err != nil {
		return err
	}
	payoutFilter.PayoutStatus = payoutStatus

	return nil
}
