package serializers

import (
	"time"

	"github.com/vnlab/makeshop-payment/src/domain/models"
)

// MerchantResponse represents the serialized merchant for API responses and documentation
type MerchantResponse struct {
	ID                int    `json:"id"`
	PaymentProviderID int    `json:"payment_provider_id"`
	PaymentMerchantID string `json:"payment_merchant_id"`
	MerchantName      string `json:"merchant_name"`
	ShopID            int    `json:"shop_id"`
	ShopURL           string `json:"shop_url"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
	ReviewStatus      *int   `json:"review_status,omitempty"`
}

// MerchantListResponse is the API response for the merchant list endpoint
type MerchantListResponse struct {
	Merchants []MerchantResponse `json:"merchants"`
	Total     int                `json:"total"`
	Page      int                `json:"page"`
	PageSize  int                `json:"page_size"`
}

// MerchantSerializer handles serialization of merchant model to response format
type MerchantSerializer struct {
	Merchant *models.Merchant
}

// NewMerchantSerializer creates a new merchant serializer instance
func NewMerchantSerializer(merchant *models.Merchant) *MerchantSerializer {
	return &MerchantSerializer{Merchant: merchant}
}

// Serialize converts a merchant model to response format
func (s *MerchantSerializer) Serialize() *MerchantResponse {
	if s.Merchant == nil {
		return nil
	}

	response := &MerchantResponse{
		ID:                s.Merchant.ID,
		PaymentProviderID: s.Merchant.PaymentProviderID,
		PaymentMerchantID: s.Merchant.PaymentMerchantID,
		MerchantName:      s.Merchant.MerchantName,
		ShopID:            s.Merchant.ShopID,
		ShopURL:           s.Merchant.ShopURL,
		CreatedAt:         s.Merchant.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         s.Merchant.UpdatedAt.Format(time.RFC3339),
	}

	if s.Merchant.MerchantPaymentProviderReview != nil {
		response.ReviewStatus = &s.Merchant.MerchantPaymentProviderReview.MerchantReviewStatus
	}

	return response
}

// SerializeMerchantCollection converts a slice of merchant models to response format
func SerializeMerchantCollection(merchants []models.Merchant) []*MerchantResponse {
	result := make([]*MerchantResponse, len(merchants))

	for i, merchant := range merchants {
		merchantCopy := merchant // Create a copy to avoid pointer issues
		result[i] = NewMerchantSerializer(&merchantCopy).Serialize()
	}

	return result
}

// NewMerchantListResponse creates a response for the merchant list endpoint
func NewMerchantListResponse(merchants []models.Merchant, total, page, pageSize int) MerchantListResponse {
	return MerchantListResponse{
		Merchants: MerchantResponseArrayToSlice(SerializeMerchantCollection(merchants)),
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
	}
}

// MerchantResponseArrayToSlice converts an array of pointers to a slice of values
func MerchantResponseArrayToSlice(merchantResponses []*MerchantResponse) []MerchantResponse {
	result := make([]MerchantResponse, len(merchantResponses))
	for i, m := range merchantResponses {
		if m != nil {
			result[i] = *m
		}
	}
	return result
}
