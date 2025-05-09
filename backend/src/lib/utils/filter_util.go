package utils

import (
	"net/http"
	"strconv"
	"time"

	validator "github.com/huydq/test/src/api/http/validator/common"
	"github.com/huydq/test/src/domain/repositories/filter"
)

// ExtractPaginationAndSorting extracts pagination and sorting parameters from a request
// and creates an initialized BaseFilter
func ExtractPaginationAndSorting(r *http.Request) *filter.BaseFilter {
	page, pageSize := ExtractPaginationParams(r)
	sortField := r.URL.Query().Get("sort_field")
	sortOrder := r.URL.Query().Get("sort_order")

	baseFilter := &filter.BaseFilter{}
	baseFilter.SetPagination(page, pageSize)

	if sortField != "" {
		direction := filter.MapSortDirection(sortOrder)
		baseFilter.SetSort(sortField, direction)
	}

	return baseFilter
}

// ExtractDateParam extracts a date parameter from a request query
func ExtractDateParam(r *http.Request, paramName string) (*time.Time, error) {
	if dateStr := r.URL.Query().Get(paramName); dateStr != "" {
		parsedTime, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			// Try alternative formats
			parsedTime, err = time.Parse(time.RFC1123, dateStr)
			if err != nil {
				parsedTime, err = time.Parse("2006-01-02", dateStr)
				if err != nil {
					return nil, err
				}
			}
		}
		return &parsedTime, nil
	}
	return nil, nil
}

// ExtractIntParam extracts an integer parameter from a request query
func ExtractIntParam(r *http.Request, paramName string) (*int, error) {
	if valueStr := r.URL.Query().Get(paramName); valueStr != "" {
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			return nil, err
		}
		return &value, nil
	}
	return nil, nil
}

// ExtractStringParam extracts a string parameter from a request query
func ExtractStringParam(r *http.Request, paramName string) *string {
	if value := r.URL.Query().Get(paramName); value != "" {
		return &value
	}
	return nil
}

// AdjustSortingRequest adjusts a SortingRequest from HTTP to domain format
func AdjustSortingRequest(sortReq *validator.SortingRequest) (string, filter.SortDirection) {
	if sortReq == nil || sortReq.SortField == "" {
		return "id", filter.Ascending
	}

	return sortReq.SortField, filter.MapSortDirection(sortReq.SortOrder)
}
