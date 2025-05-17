package util

import (
	"time"
)

// FilterOperator defines comparison operators for filtering
type FilterOperator string

const (
	Equal              FilterOperator = "="
	NotEqual           FilterOperator = "!="
	GreaterThan        FilterOperator = ">"
	GreaterThanOrEqual FilterOperator = ">="
	LessThan           FilterOperator = "<"
	LessThanOrEqual    FilterOperator = "<="
	Like               FilterOperator = "LIKE"
	In                 FilterOperator = "IN"
	Between            FilterOperator = "BETWEEN"
)

// FilterCondition represents a single filter condition
type FilterCondition struct {
	Field    string
	Operator FilterOperator
	Value    any
}

// DateFilter represents a filter for date fields
type DateFilter struct {
	Field    string
	Operator FilterOperator
	Date     *time.Time
}

// SortDirection defines sort direction
type SortDirection string

const (
	Ascending  SortDirection = "ASC"
	Descending SortDirection = "DESC"
)

// SortOption represents sorting configuration
type SortOption struct {
	Field     string
	Direction SortDirection
}

// Pagination represents pagination parameters
type Pagination struct {
	Page     int
	PageSize int
}

// BaseFilter provides common filtering capabilities
type BaseFilter struct {
	Conditions      []FilterCondition
	DateFilters     []DateFilter
	Sort            []SortOption
	Pagination      Pagination
	ValidSortFields map[string]bool
}

// AddCondition adds a filter condition
func (f *BaseFilter) AddCondition(field string, operator FilterOperator, value any) {
	f.Conditions = append(f.Conditions, FilterCondition{
		Field:    field,
		Operator: operator,
		Value:    value,
	})
}

// AddDateFilter adds a date filter
func (f *BaseFilter) AddDateFilter(field string, operator FilterOperator, date *time.Time) {
	if date != nil {
		f.DateFilters = append(f.DateFilters, DateFilter{
			Field:    field,
			Operator: operator,
			Date:     date,
		})
	}
}

// SetSort sets the sort options
func (f *BaseFilter) SetSort(field string, direction SortDirection) {
	// Validate the sort field if valid fields are specified
	if len(f.ValidSortFields) > 0 {
		if !f.ValidSortFields[field] {
			// Default to "id" if invalid field is provided
			field = "id"
		}
	}

	f.Sort = []SortOption{{Field: field, Direction: direction}}
}

// SetPagination sets pagination parameters
func (f *BaseFilter) SetPagination(page, pageSize int) {
	// Default values if invalid
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	f.Pagination = Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

// MapSortDirection maps common sort direction strings to SortDirection
func MapSortDirection(direction string) SortDirection {
	switch direction {
	case "ascend", "asc", "ASC":
		return Ascending
	case "descend", "desc", "DESC":
		return Descending
	default:
		return Ascending
	}
}
