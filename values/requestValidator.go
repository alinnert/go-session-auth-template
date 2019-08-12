package values

import (
	"reflect"
	"strings"

	"gopkg.in/go-playground/validator.v9"
)

// ValidateRequest is a validator for request bodies.
var ValidateRequest *validator.Validate

func init() {
	ValidateRequest = validator.New()

	// With this `validator.ValidationErrors[i].Field()`
	// returns the name used in the json field tag.
	ValidateRequest.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}
