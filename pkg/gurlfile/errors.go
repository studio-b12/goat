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
	ErrInvalidUse         = errors.New("invalid use statement")
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
