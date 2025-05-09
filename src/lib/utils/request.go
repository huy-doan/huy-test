package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/huydq/test/src/lib/i18n"
)

// ParseJSONBody decodes a JSON request body into the provided struct
func ParseJSONBody(r *http.Request, dst interface{}) error {
	if r.Body == nil {
		return errors.New("request body is empty")
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

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
		return 0, fmt.Errorf(i18n.T(r.Context(), "params.not_found", paramName))
	}

	value, err := strconv.Atoi(rawValue)
	if err != nil {
		return 0, fmt.Errorf(i18n.T(r.Context(), "params.invalid_number", paramName))
	}

	if value <= 0 {
		return 0, fmt.Errorf(i18n.T(r.Context(), "params.must_be_positive", paramName))
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
