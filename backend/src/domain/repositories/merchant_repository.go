package repositories

import (
	validator "github.com/huydq/test/src/api/http/validator/merchant"
	"github.com/huydq/test/src/domain/models"
)

// MerchantRepository defines the interface for merchant data operations
type MerchantRepository interface {
	ListMerchants(params validator.MerchantListFilter) ([]models.Merchant, int, error)
}
