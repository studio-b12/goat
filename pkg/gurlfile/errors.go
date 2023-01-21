package gurlfile

import (
	"errors"
	"fmt"
)

var (
	ErrTemplateAlreadyParsed       = errors.New("request template has already been parsed")
	ErrIllegalCharacter            = errors.New("illegal character")
	ErrUnexpected                  = errors.New("unexpected error")
	ErrInvalidStringLiteral        = errors.New("invalid string literal")
	ErrEmptyUsePath                = errors.New("empty use path")
	ErrInvalidSection              = errors.New("invalid section")
	ErrInvalidRequestMethod        = errors.New("invalid request method")
	ErrNoRequestURI                = errors.New("method must be followed by the request URI")
	ErrInvalidToken                = errors.New("invalid token")
	ErrInvalidBlockHeader          = errors.New("invalid block header")
	ErrInvalidBlockEntryAssignment = errors.New("block entry must start with an assignment")
	ErrInvalidHeaderKey            = errors.New("header values must start with a key")
	ErrInvalidHeaderSeparator      = errors.New("header key and value must be separated by a colon (:)")
	ErrFollowingImport             = errors.New("failed following import")
)

type InnerError struct {
	Inner error
}

func (t InnerError) Unwrap() error {
	return t.Inner
}

type ParseError struct {
	InnerError

	Line    int
	LinePos int
}

func (t ParseError) Error() string {
	return fmt.Sprintf("%d:%d: %s",
		t.Line+1, t.LinePos, t.Inner.Error())
}

type ErrorWithDetails struct {
	InnerError

	Details any
}

func newDetailedErr(inner error, details any) error {
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
