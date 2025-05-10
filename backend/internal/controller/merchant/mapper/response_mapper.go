package mapper

import (
	"time"

	"github.com/huydq/test/internal/domain/model/merchant"
	"github.com/labstack/echo/v4"
)

type MerchantResponse struct {
	ID                int    `json:"id"`
	MerchantName      string `json:"merchant_name"`
	PaymentMerchantID string `json:"payment_merchant_id"`
	PaymentProviderID int    `json:"payment_provider_id"`
	ShopID            int    `json:"shop_id"`
	ShopURL           string `json:"shop_url"`
	CreatedAt         string `json:"created_at"`
}

type MerchantListData struct {
	Merchants []MerchantResponse `json:"merchants"`
	Page      int                `json:"page"`
	PageSize  int                `json:"page_size"`
	Total     int                `json:"total"`
}

type MerchantListResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    MerchantListData `json:"data"`
}

type MerchantListSuccessMapper struct {
	ctx echo.Context
}

func NewMerchantListSuccessMapper(ctx echo.Context) *MerchantListSuccessMapper {
	return &MerchantListSuccessMapper{
		ctx: ctx,
	}
}

// ToMerchantListSuccessResponse maps the domain model to the response
func (m *MerchantListSuccessMapper) ToMerchantListSuccessResponse(
	merchants []*merchant.Merchant,
	totalPages int,
	totalCount int,
	currentPage int,
	pageSize int,
) *MerchantListResponse {
	merchantResponses := make([]MerchantResponse, len(merchants))
	for i, merchant := range merchants {
		merchantResponses[i] = MerchantResponse{
			ID:                merchant.ID,
			MerchantName:      merchant.MerchantName,
			PaymentMerchantID: merchant.PaymentMerchantID,
			PaymentProviderID: merchant.PaymentProviderID,
			ShopID:            merchant.ShopID,
			ShopURL:           merchant.ShopURL,
			CreatedAt:         merchant.CreatedAt.Format(time.RFC3339),
		}
	}

	response := &MerchantListResponse{
		Success: true,
		Message: "加盟店成功しました",
		Data: MerchantListData{
			Merchants: merchantResponses,
			Page:      currentPage,
			PageSize:  pageSize,
			Total:     totalCount,
		},
	}

	return response
}
