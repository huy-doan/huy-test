package payout

import (
	"github.com/huydq/test/internal/controller/base"
	"github.com/huydq/test/internal/controller/payout/mapper"
	"github.com/huydq/test/internal/pkg/api/generated"
	"github.com/huydq/test/internal/pkg/common/response"
	messages "github.com/huydq/test/internal/pkg/utils/messages"
	usecase "github.com/huydq/test/internal/usecase/payout"
	"github.com/labstack/echo/v4"
)

type PayoutController struct {
	base.BaseController
	payoutUsecase usecase.PayoutUsecase
}

func NewPayoutController(payoutUsecase usecase.PayoutUsecase) *PayoutController {
	return &PayoutController{
		BaseController: *base.NewBaseController(),
		payoutUsecase:  payoutUsecase,
	}
}

func (c *PayoutController) ListPayouts(ctx echo.Context) error {
	var request generated.PayoutListRequest
	if err := c.BindAndValidate(ctx, &request); err != nil {
		return response.SendError(ctx, err)
	}

	filter := mapper.ToPayoutFilter(&request)

	payouts, totalPages, totalCount, err := c.payoutUsecase.ListPayouts(ctx.Request().Context(), filter)
	if err != nil {
		return response.SendError(ctx, err)
	}

	payoutListSuccessMapper := mapper.ToPayoutListSuccessResponse(
		payouts,
		request.Page,
		request.PageSize,
		totalPages,
		totalCount,
	)

	return response.SendOK(ctx, messages.MsgListPayoutsSuccess, payoutListSuccessMapper)
}
