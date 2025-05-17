package middleware

import (
	"github.com/labstack/echo/v4"
)

// MiddlewareFunc is now an Echo middleware function
type MiddlewareFunc func(echo.HandlerFunc) echo.HandlerFunc

// MiddlewareBuilder helps with chaining Echo middleware functions
type MiddlewareBuilder struct {
	middlewares []any // can be either echo.MiddlewareFunc or func(echo.Context) error
}

// Build creates a new middleware builder
func Build(middlewares ...any) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		middlewares: middlewares,
	}
}

// With adds more middleware to the chain
func (mb *MiddlewareBuilder) With(middlewares ...any) *MiddlewareBuilder {
	mb.middlewares = append(mb.middlewares, middlewares...)
	return mb
}

// Handle applies the middleware chain to the given handler function and returns the final handler
func (mb *MiddlewareBuilder) Handle(handler echo.HandlerFunc) echo.HandlerFunc {
	// If there are no middlewares, just return the handler
	if len(mb.middlewares) == 0 {
		return handler
	}

	// Apply middlewares in reverse order (last middleware first)
	h := handler
	for i := len(mb.middlewares) - 1; i >= 0; i-- {
		middleware := mb.middlewares[i]

		if echoMiddleware, ok := middleware.(echo.MiddlewareFunc); ok {
			// If it's already an Echo middleware func, use it directly
			h = echoMiddleware(h)
		} else if fn, ok := middleware.(func(echo.HandlerFunc) echo.HandlerFunc); ok {
			// If it's our custom middleware func type
			h = fn(h)
		}
	}

	return h
}

// CreateResponseMiddleware creates a middleware that runs after the handler completes
func CreateResponseMiddleware(fn func(c echo.Context)) MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Execute the handler first
			err := next(c)
			// Then run the response logic
			fn(c)
			return err
		}
	}
}
