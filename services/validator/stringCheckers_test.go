package validator_test

import (
	"auth-server/services/validator"
	"testing"
)

func TestStringIsNotEmpty(t *testing.T) {
	tests := []struct {
		name         string
		testingValue string
		errType      errorType
	}{
		{"string is empty", "", errorTypeValidationError},
		{"string is not empty", "hello world", errorTypeNone},
	}

	for _, test := range tests {
		checker := validator.StringIsNotEmpty()
		err := checker(test.testingValue)
		testCheckerResult(t, err, test.errType, test.name)
	}
}

func TestStringIsEmail(t *testing.T) {
	tests := []struct {
		name         string
		testingValue string
		errType      errorType
	}{
		{"string is a valid e-mail", "john.doe@example.com", errorTypeNone},
		{"string is not a valid e-mail", "john.doe", errorTypeValidationError},
	}

	for _, test := range tests {
		checker := validator.StringIsEmail()
		err := checker(test.testingValue)
		testCheckerResult(t, err, test.errType, test.name)
	}
}

func TestStringIsEqualTo(t *testing.T) {
	tests := []struct {
		name           string
		testingValue   string
		comparingValue string
		errType        errorType
	}{
		{"strings are equal", "foo", "foo", errorTypeNone},
		{"strings are not equal", "foo", "bar", errorTypeValidationError},
	}

	for _, test := range tests {
		checker := validator.StringIsEqualTo(test.comparingValue, "other string")
		err := checker(test.testingValue)
		testCheckerResult(t, err, test.errType, test.name)
	}
}

func TestStringMinLength(t *testing.T) {
	tests := []struct {
		name         string
		testingValue string
		minLength    int
		errType      errorType
	}{
		{"string = 3, minlength = 3", "foo", 3, errorTypeNone},
		{"string = 6, minlength = 3", "foobar", 3, errorTypeNone},
		{"string = 3, minlength = 5", "foo", 5, errorTypeValidationError},
	}

	for _, test := range tests {
		checker := validator.StringMinLength(test.minLength)
		err := checker(test.testingValue)
		testCheckerResult(t, err, test.errType, test.name)
	}
}

func TestStringMaxLength(t *testing.T) {
	tests := []struct {
		name         string
		testingValue string
		maxLength    int
		errType      errorType
	}{
		{"string = 3, maxlength = 3", "foo", 3, errorTypeNone},
		{"string = 6, maxlength = 3", "foobar", 3, errorTypeValidationError},
		{"string = 3, maxlength = 5", "foo", 5, errorTypeNone},
	}

	for _, test := range tests {
		checker := validator.StringMaxLength(test.maxLength)
		err := checker(test.testingValue)
		testCheckerResult(t, err, test.errType, test.name)
	}
}
