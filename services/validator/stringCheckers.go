package validator

import (
	"regexp"
)

// StringChecker checks a string.
type StringChecker func(str string) error

// StringIsNotEmpty checks if a string is not empty.
func StringIsNotEmpty() StringChecker {
	return func(str string) error {
		if str == "" {
			return NewValidationError("string must not be empty")
		}
		return nil
	}
}

// StringIsEmail checks if a string is an e-mail address.
func StringIsEmail() StringChecker {
	return func(str string) error {
		ok, err := regexp.Match(`^.+@.+\..+$`, []byte(str))
		if err != nil {
			return err
		}
		if !ok {
			return NewValidationError("string is not a valid e-mail address")
		}
		return nil
	}
}

// StringIsEqualTo checks if two strings are equal
func StringIsEqualTo(other, otherName string) StringChecker {
	return func(str string) error {
		if str != other {
			return NewValidationError("string is not equal to '%s'", otherName)
		}
		return nil
	}
}

// StringMinLength checks a string against a min length
func StringMinLength(min int) StringChecker {
	return func(str string) error {
		if len(str) < min {
			return NewValidationError("string is shorter than min length %d", min)
		}
		return nil
	}
}

// StringMaxLength checks a string against a max length
func StringMaxLength(max int) StringChecker {
	return func(str string) error {
		if len(str) > max {
			return NewValidationError("string is longer than max length %d", max)
		}
		return nil
	}
}
