package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
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

// ExtractIDFromPath extracts an ID from the URL path
// Example: for path "/api/v1/users/123", ExtractIDFromPath(r, "/api/v1/users/") would return 123
func ExtractIDFromPath(r *http.Request, prefix string) (int, error) {
	path := r.URL.Path
	if len(path) <= len(prefix) {
		return 0, errors.New("invalid path")
	}

	// Remove the prefix from the path
	idStr := path[len(prefix):]

	// Convert to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New("invalid ID format")
	}

	return id, nil
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
