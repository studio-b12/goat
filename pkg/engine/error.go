package engine

import "github.com/studio-b12/goat/pkg/errs"

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
