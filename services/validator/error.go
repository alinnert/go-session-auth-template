package validator

import "fmt"

// Error occurs if the validator is used in a wrong or undefined way.
type Error struct {
	err string
}

func (err Error) Error() string {
	return fmt.Sprintf(err.err)
}

// NewError creates a new Error with fmt.Sprintf().
func NewError(msg string, params ...interface{}) Error {
	return Error{fmt.Sprintf(msg, params...)}
}
