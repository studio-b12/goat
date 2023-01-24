package executor

import (
	"fmt"

	"github.com/studio-b12/gurl/pkg/errs"
)

// BatchExecutionError holds and array of errors
// as well as the total amount of executions.
type BatchExecutionError struct {
	Inner errs.Errors
	Total int
}

func (t BatchExecutionError) Error() string {
	return fmt.Sprintf("%02d of %02d batches failed", len(t.Inner), t.Total)
}

func (t BatchExecutionError) Unwrap() error {
	return t.Inner.Condense()
}

// ParamsParsingError is an alias for error
// to identify that an error originates from
// the parameter parsing step.
type ParamsParsingError error
