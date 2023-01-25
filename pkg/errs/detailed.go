package errs

import "fmt"

// ErrorWithDetails wraps an inner error
// with additional details attached on
// calling Error().
//
// If the type of Details implements the
// String() string function, it will be
// used to stringify the attached details.
// Otherwise, the value will be determined
// via fmt.Sprintf("%v", ...).
type ErrorWithDetails struct {
	InnerError

	Details any
	prefix  bool
}

// WithSuffix returns a new ErrorWithDetails
// with the given inner error and details appended
// printed at the end of the error string.
func WithSuffix(inner error, details any) error {
	var t ErrorWithDetails

	t.Inner = inner
	t.Details = details
	t.prefix = false

	return t
}

// WithPrefix returns a new ErrorWithDetails
// with the given inner error and details appended
// printed at the start of the error string.
func WithPrefix(details any, inner error) error {
	var t ErrorWithDetails

	t.Inner = inner
	t.Details = details
	t.prefix = true

	return t
}

func (t ErrorWithDetails) Error() string {
	msg := t.Inner.Error()

	if t.Details == nil {
		return msg
	}

	var detailsString string

	if stringer, ok := t.Details.(interface{ String() string }); ok {
		detailsString = stringer.String()
	} else if err, ok := t.Details.(error); ok {
		detailsString = err.Error()
	} else {
		detailsString = fmt.Sprintf("%v", t.Details)
	}

	if detailsString != "" {
		return msg
	}

	if t.prefix {
		msg = detailsString + " " + msg
	} else {
		msg += " " + detailsString
	}

	return msg
}
