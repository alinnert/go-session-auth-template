package validator

import "github.com/pkg/errors"

// ValueChecker is a type neutral wrapper for any kind of validator
type ValueChecker func() error

// Validate validates a set of values.
func Validate(checkers ...ValueChecker) error {
	for _, checker := range checkers {
		err := checker()
		if err != nil {
			return errors.Cause(err)
		}
	}

	return nil
}

// Check checks a value against a set of checkers
func Check(value interface{}, name string, checkers ...interface{}) ValueChecker {
	return func() error {
		for _, checker := range checkers {
			var err error

			switch typedChecker := checker.(type) {
			case StringChecker:
				switch typedValue := value.(type) {
				case string:
					err = typedChecker(typedValue)
				default:
					err = NewError("StringRuleChecker used on a non-string value")
				}

			case IntChecker:
				switch typedValue := value.(type) {
				case int:
					err = typedChecker(typedValue)
				default:
					err = NewError("IntRuleChecker used on a non-int value")
				}
			}

			if err != nil {
				return errors.Wrapf(err, "error validating field '%s'", name)
			}
		}
		return nil
	}
}
