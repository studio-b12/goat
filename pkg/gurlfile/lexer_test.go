package gurlfile

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanString(t *testing.T) {
	t.Run("empty-eof", func(t *testing.T) {
		r := strings.NewReader(``)
		tk, lit := newScanner(r).scanString()
		assert.Equal(t, tokSTRING, tk)
		assert.Equal(t, "", lit)
	})

	t.Run("empty-lf", func(t *testing.T) {
		r := strings.NewReader("\n")
		tk, lit := newScanner(r).scanString()
		assert.Equal(t, tokSTRING, tk)
		assert.Equal(t, "", lit)
	})

	t.Run("unquoted-eof", func(t *testing.T) {
		r := strings.NewReader(`foo`)
		tk, lit := newScanner(r).scanString()
		assert.Equal(t, tokSTRING, tk)
		assert.Equal(t, "foo", lit)
	})

	t.Run("unquoted-lf", func(t *testing.T) {
		r := strings.NewReader("foo\n")
		tk, lit := newScanner(r).scanString()
		assert.Equal(t, tokSTRING, tk)
		assert.Equal(t, "foo", lit)
	})

	t.Run("unquoted-space", func(t *testing.T) {
		r := strings.NewReader("foo bar")
		tk, lit := newScanner(r).scanString()
		assert.Equal(t, tokSTRING, tk)
		assert.Equal(t, "foo", lit)
	})

	t.Run("quoted-double", func(t *testing.T) {
		r := strings.NewReader(`"hello world"`)
		tk, lit := newScanner(r).scanString()
		assert.Equal(t, tokSTRING, tk)
		assert.Equal(t, "hello world", lit)
	})

	t.Run("quoted-single", func(t *testing.T) {
		r := strings.NewReader(`'hello world'`)
		tk, lit := newScanner(r).scanString()
		assert.Equal(t, tokSTRING, tk)
		assert.Equal(t, "hello world", lit)
	})

	t.Run("quoted-mixed", func(t *testing.T) {
		r := strings.NewReader(`'hello "world"'`)
		tk, lit := newScanner(r).scanString()
		assert.Equal(t, tokSTRING, tk)
		assert.Equal(t, `hello "world"`, lit)
	})
}

func TestScanNumber(t *testing.T) {
	t.Run("integer-eof", func(t *testing.T) {
		r := strings.NewReader(`123`)
		tk, lit := newScanner(r).scanNumber()
		assert.Equal(t, tokINTEGER, tk)
		assert.Equal(t, "123", lit)
	})

	t.Run("integer-lf", func(t *testing.T) {
		r := strings.NewReader("123\n")
		tk, lit := newScanner(r).scanNumber()
		assert.Equal(t, tokINTEGER, tk)
		assert.Equal(t, "123", lit)
	})

	t.Run("integer-sub", func(t *testing.T) {
		r := strings.NewReader("123someident")
		tk, lit := newScanner(r).scanNumber()
		assert.Equal(t, tokINTEGER, tk)
		assert.Equal(t, "123", lit)
	})

	t.Run("float-eof", func(t *testing.T) {
		r := strings.NewReader(`123.456`)
		tk, lit := newScanner(r).scanNumber()
		assert.Equal(t, tokFLOAT, tk)
		assert.Equal(t, "123.456", lit)
	})

	t.Run("float-lf", func(t *testing.T) {
		r := strings.NewReader("123.456\n")
		tk, lit := newScanner(r).scanNumber()
		assert.Equal(t, tokFLOAT, tk)
		assert.Equal(t, "123.456", lit)
	})

	t.Run("float-sub", func(t *testing.T) {
		r := strings.NewReader("123.456someident")
		tk, lit := newScanner(r).scanNumber()
		assert.Equal(t, tokFLOAT, tk)
		assert.Equal(t, "123.456", lit)
	})
}
