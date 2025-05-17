package i18n

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Translation struct {
	ID          string            `json:"id"`
	Translation map[string]string `json:"translation"`
}

type I18n struct {
	translations map[string]map[string]Translation
	defaultLang  string
}

var instance *I18n

// Init initializes the i18n system
func Init() error {
	if instance == nil {
		instance = &I18n{
			translations: make(map[string]map[string]Translation),
			defaultLang:  "ja",
		}
	}

	dir := "src/lib/i18n/translations"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("translations directory does not exist: %s", dir)
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read translations directory: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			lang := file.Name()[:len(file.Name())-5]

			translations := make(map[string]Translation)
			data, err := os.ReadFile(filepath.Join(dir, file.Name()))
			if err != nil {
				return fmt.Errorf("failed to read translation file %s: %w", file.Name(), err)
			}

			if err := json.Unmarshal(data, &translations); err != nil {
				return fmt.Errorf("failed to parse translation file %s: %w", file.Name(), err)
			}

			instance.translations[lang] = translations
		}
	}

	return nil
}

// getLanguageFromContext extracts language from context or returns default
func (i *I18n) getLanguageFromContext(ctx context.Context) string {
	if ctx == nil {
		return i.defaultLang
	}

	if val := ctx.Value("language"); val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return i.defaultLang
}

// getTranslation retrieves the translation for a given key and language
func (i *I18n) getTranslation(key, lang string) (string, bool) {
	translations, ok := i.translations[lang]
	if !ok {
		translations = i.translations[i.defaultLang]
	}

	translation, ok := translations[key]
	if !ok {
		return "", false
	}

	text, ok := translation.Translation[lang]
	if !ok {
		text = translation.Translation[i.defaultLang]
	}

	return text, true
}

// T translates a key to the specified language and formats it with the given parameters
func T(ctx context.Context, key string, args ...interface{}) string {
	if instance == nil {
		if err := Init(); err != nil {
			return key
		}
	}

	lang := instance.getLanguageFromContext(ctx)
	text, found := instance.getTranslation(key, lang)

	if !found {
		return key
	}

	if len(args) == 0 {
		return text
	}

	return fmt.Sprintf(text, args...)
}

// WithLanguage adds language to context
func WithLanguage(ctx context.Context, lang string) context.Context {
	return context.WithValue(ctx, "language", lang)
}
