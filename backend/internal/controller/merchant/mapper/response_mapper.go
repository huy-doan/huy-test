package mapper

import (
	model "github.com/huydq/test/internal/domain/model/merchant"
	generated "github.com/huydq/test/internal/pkg/api/generated"
	"github.com/labstack/echo/v4"
)

type MerchantListSuccessMapper struct {
	ctx echo.Context
}

func NewMerchantListSuccessMapper(ctx echo.Context) *MerchantListSuccessMapper {
	return &MerchantListSuccessMapper{
		ctx: ctx,
	}
}

func (m *MerchantListSuccessMapper) ToMerchantListData(
	merchants []*model.Merchant,
	totalPages int,
	totalCount int,
	currentPage int,
	pageSize int,
) *generated.MerchantListResponse {
	merchantResponses := make([]generated.Merchant, len(merchants))
	for i, merchant := range merchants {
		merchantResponses[i] = generated.Merchant{
			Id:                merchant.ID,
			MerchantName:      merchant.MerchantName,
			PaymentMerchantId: merchant.PaymentMerchantID,
			PaymentProviderId: merchant.PaymentProviderID,
			ShopId:            merchant.ShopID,
			ShopUrl:           merchant.ShopURL,
			CreatedAt:         merchant.CreatedAt,
			UpdatedAt:         merchant.UpdatedAt,
		}
	}

	return &generated.MerchantListResponse{
		Merchants: merchantResponses,
		Page:      currentPage,
		PageSize:  pageSize,
		Total:     totalCount,
	}
}
