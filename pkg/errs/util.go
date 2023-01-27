package errs

import "errors"

// IsOfType returns true when the given
// error is of the type of T.
//
// If not and the error can be unwrapped,
// the unwrapped error will be checked
// until it either matches the type T or
// can not be further unwrapped.
func IsOfType[T any](err error) bool {
	_, ok := err.(T)
	if ok {
		return true
	}

	err = errors.Unwrap(err)
	if err != nil {
		return IsOfType[T](err)
	}

	return false
}
