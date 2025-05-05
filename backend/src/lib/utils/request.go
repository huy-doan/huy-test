package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/vnlab/makeshop-payment/src/lib/i18n"
)

// ParseJSONBody decodes a JSON request body into the provided struct
func ParseJSONBody(r *http.Request, dst any) error {
	if r.Body == nil {
		return errors.New("request body is empty")
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	// Always close the body, but check for errors
	defer func() {
		if cerr := r.Body.Close(); cerr != nil {
			if err == nil {
				err = cerr
			}
		}
	}()

	// If the body is empty, return an error
	if len(body) == 0 {
		return errors.New("request body is empty")
	}

	// Decode the JSON
	return json.Unmarshal(body, dst)
}

// ExtractParamFromPath extracts a parameter from the URL path
func ExtractParamFromPath(r *http.Request, paramName string) (int, error) {
	rawValue := r.PathValue(paramName)
	if rawValue == "" {
		errorMsg := i18n.T(r.Context(), "params.not_found", paramName)
		return 0, fmt.Errorf("%s", errorMsg)
	}

	value, err := strconv.Atoi(rawValue)
	if err != nil {
		errorMsg := i18n.T(r.Context(), "params.invalid_number", paramName)
		return 0, fmt.Errorf("%s", errorMsg)
	}

	if value <= 0 {
		errorMsg := i18n.T(r.Context(), "params.must_be_positive", paramName)
		return 0, fmt.Errorf("%s", errorMsg)
	}

	return value, nil
}

// GetQueryParam retrieves a query parameter with default value
func GetQueryParam(r *http.Request, name string, defaultValue string) string {
	value := r.URL.Query().Get(name)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetQueryParamInt retrieves an integer query parameter with default value
func GetQueryParamInt(r *http.Request, name string, defaultValue int) int {
	strValue := r.URL.Query().Get(name)
	if strValue == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		return defaultValue
	}

	return intValue
}

func ExtractPaginationParams(r *http.Request) (int, int) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	page := 1
	if pageStr != "" {
		if pageVal, err := strconv.Atoi(pageStr); err == nil && pageVal > 0 {
			page = pageVal
		}
	}

	pageSize := 10
	if pageSizeStr != "" {
		if pageSizeVal, err := strconv.Atoi(pageSizeStr); err == nil && pageSizeVal > 0 && pageSizeVal <= 100 {
			pageSize = pageSizeVal
		}
	}

	return page, pageSize
}

// GetQueryParamTime retrieves a time query parameter with default value
func GetQueryParamTime(r *http.Request, name string, layout string) *time.Time {
	value := r.URL.Query().Get(name)
	if value == "" {
		return nil
	}

	parsed, err := time.Parse(layout, value)
	if err != nil {
		return nil
	}

	return &parsed
}

// GetQueryParamIntSlice retrieves multiple integer values from query parameters
func GetQueryParamIntSlice(r *http.Request, name string) []int {
	var values []int
	queryValues := r.URL.Query()[name]
	for _, value := range queryValues {
		if intValue, err := strconv.Atoi(value); err == nil && intValue > 0 {
			values = append(values, intValue)
		}
	}
	return values
}
