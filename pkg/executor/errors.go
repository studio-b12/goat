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

// ParamsParsingError wraps an error occurred during
// parameter parsing.
type ParamsParsingError struct {
	errs.InnerError
}

func NewParamsParsingError(err error) error {
	return ParamsParsingError{
		InnerError: errs.InnerError{
			Inner: err,
		},
	}
}

// TeardownError wraps an error occurred in a teardown step.
type TeardownError struct {
	errs.InnerError
}

func NewTeardownError(err error) error {
	return TeardownError{
		InnerError: errs.InnerError{
			Inner: err,
		},
	}
}

type NoAbortError struct {
	errs.InnerError
}

func NewNoAbortError(err error) error {
	return NoAbortError{
		InnerError: errs.InnerError{
			Inner: err,
		},
	}
}
