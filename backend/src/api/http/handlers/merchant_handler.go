package handlers

import (
	"net/http"
	"time"

	"fmt"

	"github.com/vnlab/makeshop-payment/src/api/http/response"
	"github.com/vnlab/makeshop-payment/src/api/http/serializers"
	validator "github.com/vnlab/makeshop-payment/src/api/http/validator/merchant"
	"github.com/vnlab/makeshop-payment/src/lib/i18n"
	"github.com/vnlab/makeshop-payment/src/lib/utils"
	"github.com/vnlab/makeshop-payment/src/usecase"
)

// MerchantHandler handles HTTP requests related to merchants
type MerchantHandler struct {
	merchantUsecase *usecase.MerchantUsecase
}

// NewMerchantHandler creates a new merchant handler instance
func NewMerchantHandler(merchantUsecase *usecase.MerchantUsecase) *MerchantHandler {
	return &MerchantHandler{
		merchantUsecase: merchantUsecase,
	}
}

// ListMerchants handles the GET request to list merchants with optional filtering
// @Summary List merchants
// @Description Get a list of merchants with optional filtering and pagination
// @Tags Merchant Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query integer false "Page number (default: 1)"
// @Param page_size query integer false "Page size (default: 10)"
// @Param search query string false "Search query for merchant name, ID, etc."
// @Param review_status query array false "Filter by review status (values 1-3, can use multiple times: review_status=1&review_status=2)"
// @Param created_at_start query string false "Filter by creation date start (RFC3339 format)"
// @Param created_at_end query string false "Filter by creation date end (RFC3339 format)"
// @Param sort_field query string false "Field to sort by (e.g., created_at, merchant_name, etc.)"
// @Param sort_order query string false "Sort order (asc or desc, default: desc)"
// @Success 200 {object} response.Response{data=serializers.MerchantListResponse} "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /admin/merchants [get]
func (h *MerchantHandler) ListMerchants(w http.ResponseWriter, r *http.Request) {
	page := utils.GetQueryParamInt(r, "page", 1)
	pageSize := utils.GetQueryParamInt(r, "page_size", 10)
	search := r.URL.Query().Get("search")
	createdAtStart := utils.GetQueryParamTime(r, "created_at_start", time.RFC3339)
	createdAtEnd := utils.GetQueryParamTime(r, "created_at_end", time.RFC3339)
	sortField := r.URL.Query().Get("sort_field")
	sortOrder := r.URL.Query().Get("sort_order")
	reviewStatuses := utils.GetQueryParamIntSlice(r, "review_status")

	filter := validator.MerchantListFilter{
		Page:           page,
		PageSize:       pageSize,
		Search:         search,
		ReviewStatus:   reviewStatuses,
		CreatedAtStart: createdAtStart,
		CreatedAtEnd:   createdAtEnd,
		SortField:      sortField,
		SortOrder:      sortOrder,
	}

	result, err := h.merchantUsecase.ListMerchants(filter)
	if err != nil {
		response.Error(w, err)
		return
	}

	responseData := serializers.NewMerchantListResponse(
		result.Merchants,
		result.Total,
		result.Page,
		result.PageSize,
	)

	successMsg := fmt.Sprintf(i18n.T(r.Context(), "common.success"), "加盟店")

	response.Success(w, responseData, successMsg)
}
