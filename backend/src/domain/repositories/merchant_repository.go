package repositories

import (
	validator "github.com/vnlab/makeshop-payment/src/api/http/validator/merchant"
	"github.com/vnlab/makeshop-payment/src/domain/models"
)

// MerchantRepository defines the interface for merchant data operations
type MerchantRepository interface {
	ListMerchants(params validator.MerchantListFilter) ([]models.Merchant, int, error)
}
