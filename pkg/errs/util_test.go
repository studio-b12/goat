package errs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsTypeOf(t *testing.T) {
	assert.True(t,
		IsOfType[stringError](stringError("test")))
	assert.True(t,
		IsOfType[structError](structError{"test"}))
	assert.True(t,
		IsOfType[stringError](InnerError{Inner: stringError("test")}))

	assert.False(t,
		IsOfType[structError](stringError("test")))
	assert.False(t,
		IsOfType[stringError](structError{"test"}))
	assert.False(t,
		IsOfType[structError](InnerError{Inner: stringError("test")}))
}

type stringError string

func (t stringError) Error() string { return string(t) }

type structError struct {
	msg string
}

func (t structError) Error() string { return t.msg }
