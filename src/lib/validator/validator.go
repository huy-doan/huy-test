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

		// Register validators
		validate.RegisterValidation("kana", validateKana)

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

// ValidateStruct validates a struct
func ValidateStruct(obj interface{}) error {
	return validate.Struct(obj)
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
