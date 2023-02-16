package engine

import "github.com/studio-b12/goat/pkg/errs"

// Exception wraps an engine execution error
// and holds a simple message concluding the
// error.
type Exception struct {
	errs.InnerError

	Msg string
}

func (t Exception) Error() string {
	if t.Msg == "" {
		return "<unknown exception>"
	}
	return t.Msg
}
