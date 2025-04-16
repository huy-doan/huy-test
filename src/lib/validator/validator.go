package validator

import (
	"reflect"
	"regexp"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	once     sync.Once
	validate *validator.Validate
)

// Initialize validator with custom validations
func init() {
	once.Do(func() {
		validate = validator.New()

		// Register custom validations
		validate.RegisterValidation("kana", validateKana)
		validate.RegisterValidation("password_policy", PasswordPolicy)
		// Use JSON tag names for validation errors
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get("json")
			if name == "" {
				name = fld.Name
			}
			return name
		})
	})
}

// GetValidate returns the validator engine
func GetValidate() *validator.Validate {
	return validate
}

// Setup re-initializes the validator if needed (e.g. for testing)
func Setup() {
	// This is now a no-op since init() already sets up the validator
	// But we keep it for compatibility with existing code
}

// Validates if string contains only Katakana characters
func validateKana(fl validator.FieldLevel) bool {
	kanaStr := fl.Field().String()
	if kanaStr == "" {
		return true
	}

	kanaPattern := regexp.MustCompile(`^[0-9０-９ァ-ヶｦ-ﾟー]+$`)
	return kanaPattern.MatchString(kanaStr)
}

// PasswordPolicy validates that a password meets security requirements:
// - Minimum 12 characters
// - At least 1 uppercase letter
// - At least 1 lowercase letter
// - At least 1 number
// - At least 1 special character
func PasswordPolicy(fl validator.FieldLevel) bool {
	var (
		upperPattern   = regexp.MustCompile(`[A-Z]`)
		numberPattern  = regexp.MustCompile(`[0-9]`)
		lowerPattern   = regexp.MustCompile(`[a-z]`)
		specialPattern = regexp.MustCompile(`[^a-zA-Z0-9]`)
		password       = fl.Field().String()
	)

	return len(password) >= 12 &&
		upperPattern.MatchString(password) &&
		lowerPattern.MatchString(password) &&
		numberPattern.MatchString(password) &&
		specialPattern.MatchString(password)
}
