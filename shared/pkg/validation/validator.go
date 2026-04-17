package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validate *validator.Validate
}

func NewValidator() *CustomValidator {
	v := validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.Split(fld.Tag.Get("json"), ",")[0]
		if name == "-" {
			return ""
		}
		return name
	})

	v.RegisterValidation("document", func(fl validator.FieldLevel) bool {
		doc := fl.Field().String()
		return len(doc) == 11 || len(doc) == 14
	})

	return &CustomValidator{validate: v}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validate.Struct(i)
}
