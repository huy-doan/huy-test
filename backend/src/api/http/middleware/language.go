package middleware

import (
	"fmt"
	"net/http"

	"github.com/huydq/test/src/lib/i18n"
)

var initialized bool

// LanguageMiddleware adds language to context
func LanguageMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !initialized {
			if err := i18n.Init(); err != nil {
				next.ServeHTTP(w, r)
				return
			}

			initialized = true
		}

		lang := r.Header.Get("Accept-Language")
		if lang == "" {
			lang = "ja"
		}

		ctx := i18n.WithLanguage(r.Context(), lang)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// InitI18n initializes the i18n system
func InitI18n() error {
	if err := i18n.Init(); err != nil {
		return fmt.Errorf("failed to initialize i18n: %w", err)
	}
	initialized = true

	return nil
}
