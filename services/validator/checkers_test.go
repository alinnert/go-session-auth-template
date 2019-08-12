package validator_test

import (
	"auth-server/services/validator"
	"testing"
)

type errorType int

const (
	errorTypeNone errorType = iota
	errorTypeValidationError
	errorTypeBaseError
)

func testCheckerResult(
	t *testing.T, err error, errType errorType, name string,
) {
	switch err.(type) {
	case validator.ValidationError:
		if errType != errorTypeValidationError {
			t.Errorf("Test '%s' failed: Expected checker to return ValidationError", name)
		}

	case validator.Error:
		if errType != errorTypeBaseError {
			t.Errorf("Test '%s' failed: Expected checker to return Error", name)
		}

	case nil:
		if errType != errorTypeNone {
			t.Errorf("Test '%s' failed: Expected checker to return no error", name)
		}
	}
}
