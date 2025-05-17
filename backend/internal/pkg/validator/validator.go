package validator

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewValidator() *CustomValidator {
	v := validator.New()

	// Register custom validations
	_ = v.RegisterValidation("kana", validateKana)
	_ = v.RegisterValidation("password_policy", passwordPolicy)

	// Use JSON tag names for validation errors
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("json")
		if name == "" {
			name = fld.Name
		}

		// Split on comma if tag has multiple values (like json:"name,omitempty")
		if comma := strings.Index(name, ","); comma != -1 {
			name = name[:comma]
		}

		return name
	})

	return &CustomValidator{
		validator: v,
	}
}

func (cv *CustomValidator) Validate(i any) error {
	return cv.validator.Struct(i)
}

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
func passwordPolicy(fl validator.FieldLevel) bool {
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
