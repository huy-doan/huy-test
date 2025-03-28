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

// GetPathParam extracts a parameter from the URL path
// Example: for path "/api/users/123", GetPathParam(r, "id", "/api/users/") would return "123"
func GetPathParam(r *http.Request, paramName, prefix string) string {
	path := r.URL.Path
	if len(path) <= len(prefix) {
		return ""
	}
	
	// Remove the prefix from the path
	paramPath := path[len(prefix):]
	
	// In a simple case, the rest is the parameter
	return paramPath
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
