package executor

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/studio-b12/goat/pkg/errs"
	"github.com/zekrotja/rogu/log"
)

var (
	ErrCanceled = errors.New("canceled")
)

type BatchExecutionError struct {
	errs.InnerError
	Path string
}

func wrapBatchExecutionError(err error, path string) *BatchExecutionError {
	var batchErr BatchExecutionError
	batchErr.Inner = err
	batchErr.Path = path
	return &batchErr
}

// BatchResultError holds and array of errors
// as well as the total amount of executions.
type BatchResultError struct {
	Inner errs.Errors
	Total int
}

func (t *BatchResultError) Error() string {
	return fmt.Sprintf("%02d of %02d batches failed", len(t.Inner), t.Total)
}

func (t *BatchResultError) Unwrap() error {
	return t.Inner.Condense()
}

// FailedFiles returns the list of files that have failed.
func (t *BatchResultError) FailedFiles() (files []string) {
	files = make([]string, 0, len(t.Inner))
	for _, err := range t.Inner {
		if batchErr, ok := errs.As[*BatchExecutionError](err); ok {
			bePath := batchErr.Path
			if !filepath.IsAbs(bePath) {
				absBePath, err := filepath.Abs(bePath)
				if err == nil {
					bePath = absBePath
				} else {
					log.Error().Err(err).Field("dir", bePath).Msg("Failed getting absolute path to dir")
				}
			}
			files = append(files, bePath)
		}
	}
	return files
}

// ErrorMessages returns a list of the inner errors
// as strings assambled from the path and error
// message.
func (t *BatchResultError) ErrorMessages() []string {
	errMsgs := make([]string, 0, len(t.Inner))
	for _, err := range t.Inner {
		if batchErr, ok := errs.As[*BatchExecutionError](err); ok {
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
