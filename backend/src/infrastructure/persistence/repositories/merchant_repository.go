package repositories

import (
	validator "github.com/huydq/test/src/api/http/validator/merchant"
	"github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/domain/repositories"
	"gorm.io/gorm"
)

// MerchantRepository is the implementation of repositories.MerchantRepository
type MerchantRepository struct {
	db *gorm.DB
}

// NewMerchantRepository creates a new instance of MerchantRepository
func NewMerchantRepository(db *gorm.DB) repositories.MerchantRepository {
	return &MerchantRepository{
		db: db,
	}
}

// ListMerchants retrieves merchants with optional filtering
func (r *MerchantRepository) ListMerchants(params validator.MerchantListFilter) ([]models.Merchant, int, error) {
	var merchants []models.Merchant
	var total int64
	query := r.db.Model(&models.Merchant{})

	query = query.Preload("MerchantPaymentProviderReview", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC").Limit(1)
	})

	if params.CreatedAtStart != nil {
		query = query.Where("merchant.created_at >= ?", *params.CreatedAtStart)
	}

	if params.CreatedAtEnd != nil {
		query = query.Where("merchant.created_at <= ?", *params.CreatedAtEnd)
	}

	if len(params.ReviewStatus) > 0 {
		query = query.Joins("LEFT JOIN payment_provider_review ON merchant.id = payment_provider_review.merchant_id")
		query = query.Where("payment_provider_review.merchant_review_status IN ?", params.ReviewStatus)
		query = query.Order("payment_provider_review.created_at DESC")
	}

	if params.Search != "" {
		query = query.Where(
			"merchant.payment_merchant_id LIKE ? OR merchant.merchant_name LIKE ? OR merchant.shop_url LIKE ?",
			"%"+params.Search+"%", "%"+params.Search+"%", "%"+params.Search+"%",
		)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if params.Page > 0 && params.PageSize > 0 {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	sortOrder := "DESC"
	if params.SortOrder != "" {
		if params.SortOrder == "asc" {
			sortOrder = "ASC"
		} else {
			sortOrder = "DESC"
		}
	}

	allowedSortFields := map[string]string{
		"id":                  "merchant.id",
		"created_at":          "merchant.created_at",
		"merchant_name":       "merchant.merchant_name",
		"payment_merchant_id": "merchant.payment_merchant_id",
		"payment_provider_id": "merchant.payment_provider_id",
		"shop_id":             "merchant.shop_id",
		"shop_url":            "merchant.shop_url",
	}

	sortField := "merchant.id"
	if params.SortField != "" {
		if dbField, ok := allowedSortFields[params.SortField]; ok {
			sortField = dbField
		}
	}

	orderBy := sortField + " " + sortOrder
	query = query.Order(orderBy)

	err = query.Find(&merchants).Error
	if err != nil {
		return nil, 0, err
	}

	return merchants, int(total), nil
}
