package errs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoin(t *testing.T) {
	t.Run("nil > single", func(t *testing.T) {
		err1 := (error)(nil)
		err2 := errors.New("err 2")

		joined := Join(err1, err2)
		assert.Equal(t, err2, joined)
	})

	t.Run("single > nil", func(t *testing.T) {
		err1 := errors.New("err 1")
		err2 := (error)(nil)

		joined := Join(err1, err2)
		assert.Equal(t, err1, joined)
	})

	t.Run("nil > multi", func(t *testing.T) {
		err1 := (error)(nil)
		err2 := Errors{
			errors.New("err sub 1"),
			errors.New("err sub 2"),
		}

		joined := Join(err1, err2)
		assert.ElementsMatch(t, Errors{
			errors.New("err sub 1"),
			errors.New("err sub 2"),
		}, joined)
	})

	t.Run("multi > nil", func(t *testing.T) {
		err1 := Errors{
			errors.New("err sub 1"),
			errors.New("err sub 2"),
		}
		err2 := (error)(nil)

		joined := Join(err1, err2)
		assert.ElementsMatch(t, Errors{
			errors.New("err sub 1"),
			errors.New("err sub 2"),
		}, joined)
	})

	t.Run("single > multi", func(t *testing.T) {
		err1 := errors.New("err 1")
		err2 := Errors{
			errors.New("err sub 1"),
			errors.New("err sub 2"),
		}

		joined := Join(err1, err2)
		assert.ElementsMatch(t, joined, append(Errors{err1}, err2...))
	})

	t.Run("multi > single", func(t *testing.T) {
		err1 := Errors{
			errors.New("err sub 1"),
			errors.New("err sub 2"),
		}
		err2 := errors.New("err 2")

		joined := Join(err1, err2)
		assert.ElementsMatch(t, joined, append(Errors{err2}, err1...))
	})

	t.Run("multi > multi", func(t *testing.T) {
		err1 := Errors{
			errors.New("err sub 1"),
			errors.New("err sub 2"),
		}
		err2 := Errors{
			errors.New("err sub 3"),
			errors.New("err sub 4"),
		}

		joined := Join(err1, err2)
		assert.ElementsMatch(t, joined, append(err1, err2...))
	})
}
