package goatfile

import (
	"errors"
	"fmt"

	"github.com/studio-b12/goat/pkg/errs"
)

var (
	ErrTemplateAlreadyParsed       = errors.New("request template has already been parsed")
	ErrTemplateAlreadyPreParsed    = errors.New("request template has already been pre-parsed")
	ErrIllegalCharacter            = errors.New("illegal character")
	ErrUnexpected                  = errors.New("unexpected error")
	ErrInvalidStringLiteral        = errors.New("invalid string literal")
	ErrEmptyUsePath                = errors.New("empty use path")
	ErrEmptyCallPath               = errors.New("empty call path")
	ErrInvalidSection              = errors.New("invalid section")
	ErrInvalidLogSection           = errors.New("invalid log section")
	ErrInvalidRequestMethod        = errors.New("invalid request method")
	ErrNoRequestURI                = errors.New("method must be followed by the request URI")
	ErrInvalidToken                = errors.New("invalid token")
	ErrInvalidLiteral              = errors.New("invalid literal")
	ErrInvalidBlockHeader          = errors.New("invalid block header")
	ErrInvalidBlockEntryAssignment = errors.New("block entry must start with an assignment")
	ErrInvalidHeaderKey            = errors.New("header values must start with a key")
	ErrInvalidHeaderSeparator      = errors.New("header key and value must be separated by a colon (:)")
	ErrNoHeaderValue               = errors.New("no header value")
	ErrFollowingImport             = errors.New("failed following import")
	ErrOpenEscapeBlock             = errors.New("open escape block")
	ErrBlockOutOfRequest           = errors.New("blocks must follow after a request head")
	ErrMultiImport                 = errors.New("the impoted resource has already been imported")
	ErrInvalidFileDescriptor       = errors.New("invalid file descriptor")
	ErrSectionDefinedMultiple      = errors.New("the section has been already defined")
	ErrUnclosedGroup               = errors.New("group has not been closed")
	ErrMissingGroup                = errors.New("missing group definition")
	ErrVarNotFound                 = errors.New("variable not found")
	ErrNotAByteArray               = errors.New("not a byte array")
)

// ParseError wraps an inner error with
// additional parsing context.
type ParseError struct {
	errs.InnerError

	Line    int
	LinePos int
}

func (t ParseError) Error() string {
	return fmt.Sprintf("%d:%d: %s",
		t.Line+1, t.LinePos+1, t.Inner.Error())
}
