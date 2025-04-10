package handlers

import (
	"net/http"

	"github.com/vnlab/makeshop-payment/src/api/http/errors"
	"github.com/vnlab/makeshop-payment/src/api/http/response"
	"github.com/vnlab/makeshop-payment/src/api/http/serializers"
	"github.com/vnlab/makeshop-payment/src/lib/utils"
	"github.com/vnlab/makeshop-payment/src/usecase"
	"github.com/vnlab/makeshop-payment/src/lib/i18n"
	validator "github.com/vnlab/makeshop-payment/src/api/http/validator/user"
	"fmt"
)

type LockedAccountHandler struct {
	lockedAccountUsecase *usecase.LockedAccountUsecase
}

func NewLockedAccountHandler(lockedAccountUsecase *usecase.LockedAccountUsecase) *LockedAccountHandler {
	return &LockedAccountHandler{
		lockedAccountUsecase: lockedAccountUsecase,
	}
}

// UpdateOrCreateLockedAccount godoc
// @Summary Update or create locked account record
// @Description Update or create a locked account record for a user
// @Tags locked-accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body validator.UpdateLockedAccountRequest true "Locked account details"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden"
// @Failure 404 {object} response.Response "Not Found"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /admin/locked-account [post]
func (h *LockedAccountHandler) UpdateOrCreateLockedAccount(w http.ResponseWriter, r *http.Request) {
	var req validator.UpdateLockedAccountRequest
	if err := utils.ParseJSONBody(r, &req); err != nil {
		response.ValidationError(w, err)
		return
	}

	if err := req.Validate(); err != nil {
		response.ValidationError(w, err)
		return
	}

	lockedAccount, err := h.lockedAccountUsecase.UpdateOrCreateLockedAccount(r.Context(), req)
	if err != nil {
		response.Error(w, errors.InternalError(i18n.T(r.Context(), "common.error")))
		return
	}

	successMsg := fmt.Sprintf(i18n.T(r.Context(), "common.success"), "ロックアカウント")

	response.Success(w, serializers.NewLockedAccountSerializer(lockedAccount).Serialize(), successMsg)
} 
