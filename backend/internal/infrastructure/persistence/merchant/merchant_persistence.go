package merchant

import (
	"context"
	"math"

	"github.com/huydq/test/internal/datastructure/inputdata"
	model "github.com/huydq/test/internal/domain/model/merchant"
	repository "github.com/huydq/test/internal/domain/repository/merchant"
	"github.com/huydq/test/internal/infrastructure/persistence/merchant/dto"
	"github.com/huydq/test/internal/pkg/database"
	"gorm.io/gorm"
)

type MerchantRepositoryImpl struct {
	db *gorm.DB
}

func NewMerchantRepository(db *gorm.DB) repository.MerchantRepository {
	return &MerchantRepositoryImpl{
		db: db,
	}
}

// ListMerchants retrieves merchants with optional filtering
func (r *MerchantRepositoryImpl) ListMerchants(ctx context.Context, params *inputdata.MerchantListInputData) ([]*model.Merchant, int, int, error) {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return nil, 0, 0, err
	}

	query := db.WithContext(ctx).Model(&dto.Merchant{})
	query = r.applyFilters(query, params)

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, 0, err
	}

	query = db.WithContext(ctx).Model(&dto.Merchant{})
	query = r.applyFilters(query, params)
	query = r.applyPagination(query, params)
	query = r.applySorting(query, params)
	query = query.Preload("MerchantPaymentProviderReview", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC").Limit(1)
	})

	var merchantDTOs []dto.Merchant
	if err := query.Find(&merchantDTOs).Error; err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(count) / float64(params.PageSize)))

	merchants := make([]*model.Merchant, len(merchantDTOs))
	for i, dto := range merchantDTOs {
		merchants[i] = dto.ToMerchantModel()
	}

	return merchants, totalPages, int(count), nil
}

// applyFilters applies search and review status filters to the query
func (r *MerchantRepositoryImpl) applyFilters(query *gorm.DB, params *inputdata.MerchantListInputData) *gorm.DB {
	if params.CreatedAtStart != nil {
		query = query.Where("merchant.created_at >= ?", params.CreatedAtStart)
	}

	if params.CreatedAtEnd != nil {
		query = query.Where("merchant.created_at <= ?", params.CreatedAtEnd)
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

	return query
}

// applyPagination applies pagination to the query
func (r *MerchantRepositoryImpl) applyPagination(query *gorm.DB, params *inputdata.MerchantListInputData) *gorm.DB {
	offset := (params.Page - 1) * params.PageSize
	return query.Offset(offset).Limit(params.PageSize)
}

// applySorting applies sorting to the query
func (r *MerchantRepositoryImpl) applySorting(query *gorm.DB, params *inputdata.MerchantListInputData) *gorm.DB {
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

	return query.Order(sortField + " " + sortOrder)
}
