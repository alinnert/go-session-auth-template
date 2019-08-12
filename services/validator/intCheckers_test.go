package validator_test

import (
	"auth-server/services/validator"
	"testing"
)

func TestIntIsMin(t *testing.T) {
	tests := []struct {
		name         string
		minValue     int
		testingValue int
		errType      errorType
	}{
		{"value is greater than min", 10, 20, errorTypeNone},
		{"value is equal to min", 10, 10, errorTypeNone},
		{"value is less than min", 10, 5, errorTypeValidationError},
	}

	for _, test := range tests {
		checker := validator.IntIsMin(test.minValue)
		err := checker(test.testingValue)
		testCheckerResult(t, err, test.errType, test.name)
	}
}

func TestIntIsMax(t *testing.T) {
	tests := []struct {
		name         string
		maxValue     int
		testingValue int
		errType      errorType
	}{
		{"value is greater than max", 10, 20, errorTypeValidationError},
		{"value is equal to max", 10, 10, errorTypeNone},
		{"value is less than max", 10, 5, errorTypeNone},
	}

	for _, test := range tests {
		checker := validator.IntIsMax(test.maxValue)
		err := checker(test.testingValue)
		testCheckerResult(t, err, test.errType, test.name)
	}
}
