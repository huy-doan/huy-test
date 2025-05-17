package merchant

import (
	"context"

	"github.com/huydq/test/internal/datastructure/inputdata"
	merchantModel "github.com/huydq/test/internal/domain/model/merchant"
	merchantRepo "github.com/huydq/test/internal/domain/repository/merchant"
)

type MerchantManagementUsecase interface {
	ListMerchants(ctx context.Context, input *inputdata.MerchantListInputData) ([]*merchantModel.Merchant, int, int, error)
}

type ManageMerchantsUsecase struct {
	merchantRepo merchantRepo.MerchantRepository
}

func NewManageMerchantsUsecase(merchantRepo merchantRepo.MerchantRepository) *ManageMerchantsUsecase {
	return &ManageMerchantsUsecase{
		merchantRepo: merchantRepo,
	}
}

// ListMerchants lists merchants with optional filtering and pagination
func (uc *ManageMerchantsUsecase) ListMerchants(ctx context.Context, input *inputdata.MerchantListInputData) ([]*merchantModel.Merchant, int, int, error) {
	const (
		defaultPage     = 1
		defaultPageSize = 10
	)

	if input.Page <= 0 {
		input.Page = defaultPage
	}

	if input.PageSize <= 0 {
		input.PageSize = defaultPageSize
	}

	return uc.merchantRepo.ListMerchants(ctx, input)
}
