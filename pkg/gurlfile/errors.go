package gurlfile

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyRequest       = errors.New("empty request")
	ErrInvalidHead        = errors.New("invalid request head")
	ErrInvalidSectionName = errors.New("invalid section name")
	ErrAlreadyParsed      = errors.New("request has already been parsed")
)

type InnerError struct {
	Inner error
}

func (t InnerError) Unwrap() error {
	return t.Inner
}

type ParseError struct {
	InnerError

	RequestNumber int
}

func newParseError(reqNumb int, inner error) (t ParseError) {
	t.RequestNumber = reqNumb
	t.Inner = inner
	return t
}

func (t ParseError) Error() string {
	return fmt.Sprintf("Request #%02d: %s",
		t.RequestNumber, t.Inner.Error())
}

type DetailedError struct {
	InnerError

	Details string
}

func newDetailedError(inner error, details string, params ...any) (t DetailedError) {
	t.Details = fmt.Sprintf(details, params...)
	t.Inner = inner
	return t
}

func (t DetailedError) Error() string {
	return fmt.Sprintf("%s\n%s",
		t.Inner.Error(), t.Details)
}
