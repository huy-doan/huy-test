package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/huydq/test/internal/pkg/logger"
	"github.com/labstack/echo/v4"
)

type responseRecorder struct {
	echo.Response
	status int
	body   *bytes.Buffer
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	return r.body.Write(b)
}

func (r *responseRecorder) WriteHeader(code int) {
	r.status = code
	r.Response.WriteHeader(code)
}

type requestLoggerConfig struct {
	Skipper                 func(c echo.Context) bool
	LogRequestBody          bool
	LogResponseBody         bool
	MaxRequestBodyLogSize   int
	MaxResponseBodyLogSize  int
	SanitizeRequestHeaders  []string
	SanitizeResponseHeaders []string
	LogRequestHeaders       bool
	LogResponseHeaders      bool
}

func (m *MiddlewareManager) RequestLoggerMiddleware() echo.MiddlewareFunc {
	config := requestLoggerConfig{
		Skipper:                 func(c echo.Context) bool { return false },
		LogRequestBody:          true,
		LogResponseBody:         false,
		MaxRequestBodyLogSize:   4096, // 4KB
		MaxResponseBodyLogSize:  4096, // 4KB
		SanitizeRequestHeaders:  []string{"Authorization", "Cookie"},
		SanitizeResponseHeaders: []string{"Set-Cookie"},
		LogRequestHeaders:       true,
		LogResponseHeaders:      true,
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()

			start := time.Now()

			// Get or generate a trace ID
			traceID := req.Header.Get("X-Trace-ID")
			if traceID == "" {
				traceID = logger.GenerateTraceID()
				req.Header.Set("X-Trace-ID", traceID)
				res.Header().Set("X-Trace-ID", traceID)
			}

			requestLogger := m.logger.WithTraceID(traceID)
			c.Set("logger", requestLogger)

			// Prepare request logging fields
			fields := map[string]any{
				"trace_id":    traceID,
				"method":      req.Method,
				"path":        req.URL.Path,
				"query":       req.URL.RawQuery,
				"remote_addr": req.RemoteAddr,
				"user_agent":  req.UserAgent(),
			}

			// Log request headers if enabled
			if config.LogRequestHeaders {
				fields["headers"] = sanitizeHeaders(req.Header, config.SanitizeRequestHeaders)
			}

			// Capture and log request body for applicable methods
			if config.LogRequestBody && (req.Method == http.MethodPost || req.Method == http.MethodPut || req.Method == http.MethodPatch) {
				var bodyBytes []byte
				if req.Body != nil {
					bodyBytes, req.Body, _ = readAndReplaceBody(req.Body, config.MaxRequestBodyLogSize)
					if len(bodyBytes) > 0 {
						fields["request_body"] = sanitizeBody(bodyBytes)
					}
				}
			}

			// Log incoming request
			requestLogger.Info("API Request", fields)

			// If we want to capture the response body, replace the response writer
			var responseBodyBuffer *bytes.Buffer
			if config.LogResponseBody {
				responseBodyBuffer = &bytes.Buffer{}
				recorder := &responseRecorder{
					Response: *res,
					body:     responseBodyBuffer,
				}
				c.Response().Writer = recorder
			}

			// Process the request
			err := next(c)

			// Log response
			responseFields := map[string]any{
				"trace_id":     traceID,
				"method":       req.Method,
				"path":         req.URL.Path,
				"status":       res.Status,
				"duration_ms":  time.Since(start).Milliseconds(),
				"content_type": res.Header().Get(echo.HeaderContentType),
			}

			// Log response headers if enabled
			if config.LogResponseHeaders {
				responseFields["headers"] = sanitizeHeaders(res.Header(), config.SanitizeResponseHeaders)
			}

			// Log captured response body if enabled
			if config.LogResponseBody && responseBodyBuffer != nil && responseBodyBuffer.Len() > 0 {
				responseBodyBytes := responseBodyBuffer.Bytes()
				if len(responseBodyBytes) > config.MaxResponseBodyLogSize {
					responseBodyBytes = responseBodyBytes[:config.MaxResponseBodyLogSize]
				}
				responseFields["response_body"] = sanitizeBody(responseBodyBytes)
			}

			// Log the error if present
			if err != nil {
				responseFields["error"] = err.Error()
				requestLogger.Error("API Response", responseFields)
			} else {
				// Log based on status code
				if res.Status >= 400 {
					requestLogger.Error("API Response", responseFields)
				} else if res.Status >= 300 {
					requestLogger.Warn("API Response", responseFields)
				} else {
					requestLogger.Info("API Response", responseFields)
				}
			}

			return err
		}
	}
}

func sanitizeHeaders(headers http.Header, sensitiveHeaders []string) map[string]string {
	result := make(map[string]string)
	for k, v := range headers {
		// Skip sensitive headers
		shouldSanitize := false
		k = http.CanonicalHeaderKey(k)

		for _, sensitive := range sensitiveHeaders {
			if k == http.CanonicalHeaderKey(sensitive) {
				shouldSanitize = true
				break
			}
		}

		// Also check for common sensitive header patterns
		lowerK := strings.ToLower(k)
		if strings.Contains(lowerK, "token") ||
			strings.Contains(lowerK, "password") ||
			strings.Contains(lowerK, "secret") ||
			strings.Contains(lowerK, "key") {
			shouldSanitize = true
		}

		if shouldSanitize {
			result[k] = "[REDACTED]"
		} else if len(v) > 0 {
			result[k] = v[0]
		}
	}
	return result
}

func sanitizeBody(bodyBytes []byte) any {
	// Try to parse as JSON to sanitize sensitive fields
	var bodyData map[string]any
	if json.Unmarshal(bodyBytes, &bodyData) == nil {
		// Remove sensitive fields if it's JSON
		for k := range bodyData {
			lowerK := strings.ToLower(k)
			if strings.Contains(lowerK, "password") ||
				strings.Contains(lowerK, "token") ||
				strings.Contains(lowerK, "secret") ||
				strings.Contains(lowerK, "key") ||
				strings.Contains(lowerK, "auth") {
				bodyData[k] = "[REDACTED]"
			}
		}
		return bodyData
	}

	// If not JSON, return as string (truncated if needed)
	bodyStr := string(bodyBytes)
	return bodyStr
}

func readAndReplaceBody(body io.ReadCloser, maxSize int) ([]byte, io.ReadCloser, error) {
	// Read up to maxSize bytes
	limitReader := io.LimitReader(body, int64(maxSize))
	bodyBytes, err := io.ReadAll(limitReader)
	if err != nil {
		return nil, body, err
	}

	// Close original body
	if err := body.Close(); err != nil {
		return bodyBytes, io.NopCloser(bytes.NewReader(bodyBytes)), err
	}

	// Create new reader from captured bytes
	return bodyBytes, io.NopCloser(bytes.NewReader(bodyBytes)), nil
}
