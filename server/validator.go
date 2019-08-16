package server

import (
	"reflect"
	"strings"

	"gopkg.in/go-playground/validator.v9"
)

// GetValidator creates an Validator instance
func GetValidator() *validator.Validate {
	validator := validator.New()

	// With this `validator.ValidationErrors[i].Field()`
	// returns the name used in the json field tag.
	validator.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return validator
}
