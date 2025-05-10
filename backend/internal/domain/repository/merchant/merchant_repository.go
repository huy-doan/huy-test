package merchant

import (
	"context"

	"github.com/huydq/test/internal/datastructure/inputdata"
	"github.com/huydq/test/internal/domain/model/merchant"
)

// MerchantRepository defines the interface for merchant data operations
type MerchantRepository interface {
	ListMerchants(ctx context.Context, params *inputdata.MerchantListInputData) ([]*merchant.Merchant, int, int, error)
}
