package validator

import "fmt"

// ValidationError occurs if on of the given values is not valid.
type ValidationError struct {
	err string
}

func (err ValidationError) Error() string {
	return fmt.Sprintf(err.err)
}

// NewValidationError creates a new ValidationError with fmt.Sprintf().
func NewValidationError(msg string, params ...interface{}) ValidationError {
	return ValidationError{fmt.Sprintf(msg, params...)}
}
