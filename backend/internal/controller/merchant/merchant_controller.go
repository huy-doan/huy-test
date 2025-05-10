package merchant

import (
	"github.com/huydq/test/internal/controller/base"
	"github.com/huydq/test/internal/controller/merchant/mapper"
	"github.com/huydq/test/internal/usecase/merchant"
	"github.com/labstack/echo/v4"
)

type MerchantController struct {
	base.BaseController
	merchantUsecase merchant.MerchantManagementUsecase
}

func NewMerchantController(merchantUsecase merchant.MerchantManagementUsecase) *MerchantController {
	return &MerchantController{
		BaseController:  *base.NewBaseController(),
		merchantUsecase: merchantUsecase,
	}
}

// ListMerchants handles the request to list merchants with pagination and filtering
func (c *MerchantController) ListMerchants(ctx echo.Context) error {
	var request mapper.MerchantListRequest
	if err := c.BindAndValidate(ctx, &request); err != nil {
		return c.SendError(ctx, err)
	}

	if request.Page <= 0 {
		request.Page = 1
	}
	if request.PageSize <= 0 {
		request.PageSize = 10
	}

	inputData := request.ToMerchantListInputData()
	merchants, totalPages, totalCount, err := c.merchantUsecase.ListMerchants(ctx.Request().Context(), inputData)
	if err != nil {
		return c.SendError(ctx, err)
	}

	merchantListSuccessMapper := mapper.NewMerchantListSuccessMapper(ctx)
	response := merchantListSuccessMapper.ToMerchantListSuccessResponse(merchants, totalPages, totalCount, request.Page, request.PageSize)

	return c.SendOK(ctx, response)
}
