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
}

// newDetailedErr returns a new ErrorWithDetails
// with the given inner error and details.
func NewDetailedErr(inner error, details any) error {
	var t ErrorWithDetails

	t.Inner = inner
	t.Details = details

	return t
}

func (t ErrorWithDetails) Error() string {
	msg := t.Inner.Error()

	if t.Details != nil {
		var detailsString string

		if stringer, ok := t.Details.(interface{ String() string }); ok {
			detailsString = stringer.String()
		} else if err, ok := t.Details.(error); ok {
			detailsString = err.Error()
		} else {
			detailsString = fmt.Sprintf("%v", t.Details)
		}

		if detailsString != "" {
			msg += " " + detailsString
		}
	}

	return msg
}
