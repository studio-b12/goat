package errs

import (
	"fmt"
	"strings"
)

// Errors adds untility functionalities to
// an array of errors.
type Errors []error

func (t Errors) Error() string {
	var sb strings.Builder

	if len(t) == 0 {
		return ""
	}

	fmt.Fprintf(&sb, "%d errors occured:\n", len(t))

	for i, err := range t {
		lines := strings.Split(err.Error(), "\n")
		for j, line := range lines {
			if j == 0 {
				continue
			}
			lines[j] = "       " + line
		}

		fmt.Fprintf(&sb, "  [%02d] %s\n", i, strings.Join(lines, "\n"))
	}

	return sb.String()
}

// HasSome returns true if the inner error
// array is not empty.
func (t Errors) HasSome() bool {
	return len(t) != 0
}

// Condense returns nil if the inner error
// array has no errors in it. Otherwise, the
// error array will be returned as Errors type.
func (t Errors) Condense() error {
	if t.HasSome() {
		return t
	}
	return nil
}

// Append adds the given error to the errors
// array. If the passed error is an Errors
// array, the errors contained are added one
// by one to the list of errors.
func (t Errors) Append(err error) Errors {
	if errs, ok := err.(Errors); ok {
		for _, e := range errs {
			t = t.Append(e)
		}
		return t
	}

	return append(t, err)
}
