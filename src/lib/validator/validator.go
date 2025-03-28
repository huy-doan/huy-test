package validator

import (
	"reflect"
	"regexp"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// CustomValidator is a validator wrapper for Gin
type CustomValidator struct {
	once     sync.Once
	validate *validator.Validate
}

// Init initializes the validator
func (v *CustomValidator) Init() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")

		// Add any custom validations here
		// For example:
		// v.validate.RegisterValidation("custom_rule", customRuleFunc)
		// Add kana validation
		v.validate.RegisterValidation("kana", validateKana)

		// Use JSON tag names for validation errors
		v.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get("json")
			if name == "" {
				name = fld.Name
			}
			return name
		})
	})
}

// Engine returns the underlying validator engine
func (v *CustomValidator) Engine() interface{} {
	v.Init()
	return v.validate
}

// ValidateStruct validates a struct
func (v *CustomValidator) ValidateStruct(obj interface{}) error {
	v.Init()
	return v.validate.Struct(obj)
}

// Setup sets up the validator for Gin
func Setup() {
	binding.Validator = &CustomValidator{}
}

// GetValidate returns the validator engine
func GetValidate() *validator.Validate {
	v := &CustomValidator{}
	v.Init()
	return v.validate
}

func validateKana(fl validator.FieldLevel) bool {
	kanaStr := fl.Field().String()
	if kanaStr == "" {
		return true
	}

	kanaPattern := regexp.MustCompile(`^[0-9０-９ァ-ヶｦ-ﾟー]+$`)
	return kanaPattern.MatchString(kanaStr)
}
