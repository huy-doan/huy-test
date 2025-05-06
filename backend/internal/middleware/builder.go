package middleware

import (
	"net/http"
)

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

type Middleware func(http.Handler) http.Handler

type MiddlewareBuilder struct {
	middlewares []any
}

func Build(middlewares ...any) *MiddlewareBuilder {
	return &MiddlewareBuilder{middlewares: middlewares}
}

func (mb *MiddlewareBuilder) With(middlewares ...any) *MiddlewareBuilder {
	mb.middlewares = append(mb.middlewares, middlewares...)
	return mb
}

// Handle applies the middleware chain to the given handler function and returns the final handler
func (mb *MiddlewareBuilder) Handle(handler http.HandlerFunc) http.HandlerFunc {
	finalHandler := handler

	for i := len(mb.middlewares) - 1; i >= 0; i-- {
		middleware := mb.middlewares[i]
		switch m := middleware.(type) {
		case MiddlewareFunc:
			finalHandler = m(finalHandler)
		case func(http.HandlerFunc) http.HandlerFunc:
			finalHandler = m(finalHandler)
		case Middleware:
			finalHandler = wrapStandardMiddleware(m, finalHandler)
		case func(http.Handler) http.Handler:
			finalHandler = wrapStandardMiddleware(m, finalHandler)
		default:
			panic("Unsupported middleware type")
		}
	}

	return finalHandler
}

// wrapStandardMiddleware wraps standard http.Handler middleware to work with http.HandlerFunc
func wrapStandardMiddleware(middleware func(http.Handler) http.Handler, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		middleware(handler).ServeHTTP(w, r)
	}
}

// CreateResponseMiddleware creates a middleware that runs after the handler completes
func CreateResponseMiddleware(fn func(w http.ResponseWriter, r *http.Request)) MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Execute the handler first
			next(w, r)
			// Then run the response logic
			fn(w, r)
		}
	}
}
