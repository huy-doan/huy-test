package persistence

import (
	"fmt"

	domainFilter "github.com/huydq/test/internal/domain/model/util"
	"gorm.io/gorm"
)

// GormFilterBuilder applies domain filters to GORM queries
type GormFilterBuilder struct{}

// ApplyBaseFilter applies base filter conditions to a GORM query
func (b *GormFilterBuilder) ApplyBaseFilter(query *gorm.DB, baseFilter *domainFilter.BaseFilter) *gorm.DB {
	// Apply standard conditions
	for _, condition := range baseFilter.Conditions {
		switch condition.Operator {
		case domainFilter.Equal:
			query = query.Where(fmt.Sprintf("%s = ?", condition.Field), condition.Value)
		case domainFilter.NotEqual:
			query = query.Where(fmt.Sprintf("%s != ?", condition.Field), condition.Value)
		case domainFilter.GreaterThan:
			query = query.Where(fmt.Sprintf("%s > ?", condition.Field), condition.Value)
		case domainFilter.GreaterThanOrEqual:
			query = query.Where(fmt.Sprintf("%s >= ?", condition.Field), condition.Value)
		case domainFilter.LessThan:
			query = query.Where(fmt.Sprintf("%s < ?", condition.Field), condition.Value)
		case domainFilter.LessThanOrEqual:
			query = query.Where(fmt.Sprintf("%s <= ?", condition.Field), condition.Value)
		case domainFilter.Like:
			query = query.Where(fmt.Sprintf("%s LIKE ?", condition.Field), condition.Value)
		case domainFilter.In:
			query = query.Where(fmt.Sprintf("%s IN (?)", condition.Field), condition.Value)
		}
	}

	// Apply date filters
	for _, dateFilter := range baseFilter.DateFilters {
		if dateFilter.Date != nil {
			switch dateFilter.Operator {
			case domainFilter.Equal:
				query = query.Where("DATE("+dateFilter.Field+") = DATE(?)", dateFilter.Date)
			case domainFilter.NotEqual:
				query = query.Where("DATE("+dateFilter.Field+") != DATE(?)", dateFilter.Date)
			case domainFilter.GreaterThan:
				query = query.Where("DATE("+dateFilter.Field+") > DATE(?)", dateFilter.Date)
			case domainFilter.GreaterThanOrEqual:
				query = query.Where("DATE("+dateFilter.Field+") >= DATE(?)", dateFilter.Date)
			case domainFilter.LessThan:
				query = query.Where("DATE("+dateFilter.Field+") < DATE(?)", dateFilter.Date)
			case domainFilter.LessThanOrEqual:
				query = query.Where("DATE("+dateFilter.Field+") <= DATE(?)", dateFilter.Date)
			}
		}
	}

	// Apply sorting
	if len(baseFilter.Sort) > 0 {
		for _, sort := range baseFilter.Sort {
			query = query.Order(fmt.Sprintf("%s %s", sort.Field, sort.Direction))
		}
	} else {
		// Default sort by ID if not specified
		query = query.Order("id DESC")
	}

	return query
}

// ApplyPagination applies pagination to a GORM query
func (b *GormFilterBuilder) ApplyPagination(query *gorm.DB, pagination domainFilter.Pagination) *gorm.DB {
	offset := (pagination.Page - 1) * pagination.PageSize
	return query.Offset(offset).Limit(pagination.PageSize)
}

// NewGormFilterBuilder creates a new GormFilterBuilder
func NewGormFilterBuilder() *GormFilterBuilder {
	return &GormFilterBuilder{}
}
