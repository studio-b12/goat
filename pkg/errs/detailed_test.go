package errs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithSuffix(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		inner := errors.New("inner error")
		err := WithSuffix(inner, "suffix")

		assert.ErrorIs(t, err, inner)
		assert.Equal(t,
			inner.Error()+" suffix",
			err.Error())
	})

	t.Run("error", func(t *testing.T) {
		inner := errors.New("inner error")
		suffixErr := errors.New("suffix error")
		err := WithSuffix(inner, suffixErr)

		assert.ErrorIs(t, err, inner)
		assert.Equal(t,
			inner.Error()+" "+suffixErr.Error(),
			err.Error())
	})

	t.Run("stringer", func(t *testing.T) {

		inner := errors.New("inner error")
		err := WithSuffix(inner, stringer("suffix stringer"))

		assert.ErrorIs(t, err, inner)
		assert.Equal(t,
			inner.Error()+" suffix stringer",
			err.Error())
	})

	t.Run("empty", func(t *testing.T) {

		inner := errors.New("inner error")
		err := WithSuffix(inner, "")

		assert.ErrorIs(t, err, inner)
		assert.Equal(t,
			inner.Error(),
			err.Error())
	})

	t.Run("nil", func(t *testing.T) {

		inner := errors.New("inner error")
		err := WithSuffix(inner, nil)

		assert.ErrorIs(t, err, inner)
		assert.Equal(t,
			inner.Error(),
			err.Error())
	})
}

func TestWithPrefix(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		inner := errors.New("inner error")
		err := WithPrefix("prefix", inner)

		assert.ErrorIs(t, err, inner)
		assert.Equal(t,
			"prefix "+inner.Error(),
			err.Error())
	})

	t.Run("error", func(t *testing.T) {
		inner := errors.New("inner error")
		prefixErr := errors.New("prefix error")
		err := WithPrefix(prefixErr, inner)

		assert.ErrorIs(t, err, inner)
		assert.Equal(t,
			prefixErr.Error()+" "+inner.Error(),
			err.Error())
	})

	t.Run("stringer", func(t *testing.T) {

		inner := errors.New("inner error")
		err := WithPrefix(stringer("prefix stringer"), inner)

		assert.ErrorIs(t, err, inner)
		assert.Equal(t,
			"prefix stringer "+inner.Error(),
			err.Error())
	})

	t.Run("empty", func(t *testing.T) {

		inner := errors.New("inner error")
		err := WithPrefix("", inner)

		assert.ErrorIs(t, err, inner)
		assert.Equal(t,
			inner.Error(),
			err.Error())
	})

	t.Run("empty", func(t *testing.T) {

		inner := errors.New("inner error")
		err := WithPrefix(nil, inner)

		assert.ErrorIs(t, err, inner)
		assert.Equal(t,
			inner.Error(),
			err.Error())
	})
}

type stringer string

func (t stringer) String() string { return string(t) }
