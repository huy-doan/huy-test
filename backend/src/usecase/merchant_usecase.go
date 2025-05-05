package usecase

import (
	validator "github.com/vnlab/makeshop-payment/src/api/http/validator/merchant"
	"github.com/vnlab/makeshop-payment/src/domain/models"
	"github.com/vnlab/makeshop-payment/src/domain/repositories"
)

// MerchantUsecase defines the business logic for merchant operations
type MerchantUsecase struct {
	merchantRepo repositories.MerchantRepository
}

// NewMerchantUsecase creates a new merchant usecase instance
func NewMerchantUsecase(merchantRepo repositories.MerchantRepository) *MerchantUsecase {
	return &MerchantUsecase{
		merchantRepo: merchantRepo,
	}
}

type ListMerchantsResponse struct {
	Merchants []models.Merchant `json:"merchants"`
	Total     int               `json:"total"`
	Page      int               `json:"page"`
	PageSize  int               `json:"page_size"`
}

// ListMerchants retrieves merchants according to the filter parameters
func (u *MerchantUsecase) ListMerchants(filter validator.MerchantListFilter) (ListMerchantsResponse, error) {
	const (
		defaultPage     = 1
		defaultPageSize = 10
	)

	if filter.Page <= 0 {
		filter.Page = defaultPage
	}

	if filter.PageSize <= 0 {
		filter.PageSize = defaultPageSize
	}

	repoParams := validator.MerchantListFilter{
		Page:           filter.Page,
		PageSize:       filter.PageSize,
		Search:         filter.Search,
		ReviewStatus:   filter.ReviewStatus,
		CreatedAtStart: filter.CreatedAtStart,
		CreatedAtEnd:   filter.CreatedAtEnd,
		SortField:      filter.SortField,
		SortOrder:      filter.SortOrder,
	}

	merchants, total, err := u.merchantRepo.ListMerchants(repoParams)
	if err != nil {
		return ListMerchantsResponse{}, err
	}

	response := ListMerchantsResponse{
		Merchants: merchants,
		Total:     total,
		Page:      filter.Page,
		PageSize:  filter.PageSize,
	}

	return response, nil
}
