package executor

import (
	"fmt"

	"github.com/studio-b12/goat/pkg/errs"
)

type BatchExecutionError struct {
	errs.InnerError
	Path string
}

func wrapBatchExecutionError(err error, path string) BatchExecutionError {
	var batchErr BatchExecutionError
	batchErr.Inner = err
	batchErr.Path = path
	return batchErr
}

// BatchResultError holds and array of errors
// as well as the total amount of executions.
type BatchResultError struct {
	Inner errs.Errors
	Total int
}

func (t BatchResultError) Error() string {
	return fmt.Sprintf("%02d of %02d batches failed", len(t.Inner), t.Total)
}

func (t BatchResultError) Unwrap() error {
	return t.Inner.Condense()
}

// ErrorMessages returns a list of the inner errors
// as strings assambled from the path and error
// message.
func (t BatchResultError) ErrorMessages() []string {
	errMsgs := make([]string, 0, len(t.Inner))
	for _, err := range t.Inner {
		if batchErr, ok := err.(BatchExecutionError); ok {
			errMsgs = append(errMsgs, fmt.Sprintf("%s: %s", batchErr.Path, batchErr.Error()))
		}
	}
	return errMsgs
}

// ParamsParsingError is an alias for error
// to identify that an error originates from
// the parameter parsing step.
type ParamsParsingError error
