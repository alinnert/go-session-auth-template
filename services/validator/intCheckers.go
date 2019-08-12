package validator

// IntChecker checks a string.
type IntChecker func(str int) error

// IntIsMin checks if an int has the provided min value.
func IntIsMin(min int) IntChecker {
	return func(num int) error {
		if num < min {
			return NewValidationError("int is less than min value %d", min)
		}
		return nil
	}
}

// IntIsMax checks if an int has the provided max value.
func IntIsMax(max int) IntChecker {
	return func(num int) error {
		if num > max {
			return NewValidationError("int is greater than max value %d", max)
		}
		return nil
	}
}
