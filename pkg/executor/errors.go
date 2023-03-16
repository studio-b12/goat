package executor

import (
	"fmt"

	"github.com/studio-b12/goat/pkg/errs"
	"github.com/studio-b12/goat/pkg/goatfile"
)

// BatchExecutionError holds and array of errors
// as well as the total amount of executions.
type BatchExecutionError struct {
	Inner errs.Errors
	Files []goatfile.Goatfile
}

func (t BatchExecutionError) Error() string {
	return fmt.Sprintf("%02d of %02d batches failed", len(t.Inner), len(t.Files))
}

func (t BatchExecutionError) Pathes() []string {
	pathes := make([]string, 0, len(t.Files))
	for _, gf := range t.Files {
		pathes = append(pathes, gf.Path)
	}
	return pathes
}

func (t BatchExecutionError) Unwrap() error {
	return t.Inner.Condense()
}

// ParamsParsingError is an alias for error
// to identify that an error originates from
// the parameter parsing step.
type ParamsParsingError error
