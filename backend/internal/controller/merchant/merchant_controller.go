package merchant

import (
	"github.com/huydq/test/internal/controller/base"
	"github.com/huydq/test/internal/controller/merchant/mapper"
	generated "github.com/huydq/test/internal/pkg/api/generated"
	response "github.com/huydq/test/internal/pkg/common/response"
	messages "github.com/huydq/test/internal/pkg/utils/messages"
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
	var request generated.MerchantListRequest
	if err := c.BindAndValidate(ctx, &request); err != nil {
		return response.SendError(ctx, err)
	}

	inputData := mapper.ToMerchantListInputData(&request)
	merchants, totalPages, totalCount, err := c.merchantUsecase.ListMerchants(ctx.Request().Context(), inputData)
	if err != nil {
		return response.SendError(ctx, err)
	}

	merchantListSuccessMapper := mapper.NewMerchantListSuccessMapper(ctx)
	merchantListData := merchantListSuccessMapper.ToMerchantListData(merchants, totalPages, totalCount, request.Page, request.PageSize)

	return response.SendOK(ctx, messages.MsgListMerchantsSuccess, merchantListData)
}
