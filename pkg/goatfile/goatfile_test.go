package goatfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	getA := func() Goatfile {
		return Goatfile{
			Setup: []Action{
				testRequest("A", "1"),
				testRequest("A", "2"),
				testRequest("A", "3"),
			},
			Tests: []Action{
				testRequest("A", "4"),
				testRequest("A", "5"),
			},
			Teardown: []Action{},
		}
	}

	getB := func() Goatfile {
		return Goatfile{
			Setup: []Action{
				testRequest("B", "1"),
				testRequest("B", "2"),
			},
			SetupEach: []Action{
				testRequest("B", "3"),
			},
			Tests: []Action{},
			Teardown: []Action{
				testRequest("B", "4"),
				testRequest("B", "5"),
			},
			TeardownEach: []Action{
				testRequest("B", "6"),
			},
		}
	}

	t.Run("a-into-b", func(t *testing.T) {
		a := getA()
		b := getB()

		a.Merge(b)

		assert.Equal(t, getB(), b)
		assert.Equal(t, Goatfile{
			Setup: []Action{
				testRequest("A", "1"),
				testRequest("A", "2"),
				testRequest("A", "3"),
				testRequest("B", "1"),
				testRequest("B", "2"),
			},
			SetupEach: []Action{
				testRequest("B", "3"),
			},
			Tests: []Action{
				testRequest("A", "4"),
				testRequest("A", "5"),
			},
			Teardown: []Action{
				testRequest("B", "4"),
				testRequest("B", "5"),
			},
			TeardownEach: []Action{
				testRequest("B", "6"),
			},
		}, a)
	})

	t.Run("b-into-a", func(t *testing.T) {
		a := getA()
		b := getB()

		b.Merge(a)

		assert.Equal(t, getA(), a)
		assert.Equal(t, Goatfile{
			Setup: []Action{
				testRequest("B", "1"),
				testRequest("B", "2"),
				testRequest("A", "1"),
				testRequest("A", "2"),
				testRequest("A", "3"),
			},
			SetupEach: []Action{
				testRequest("B", "3"),
			},
			Tests: []Action{
				testRequest("A", "4"),
				testRequest("A", "5"),
			},
			Teardown: []Action{
				testRequest("B", "4"),
				testRequest("B", "5"),
			},
			TeardownEach: []Action{
				testRequest("B", "6"),
			},
		}, b)
	})
}

// --- Helpers ---

func testRequest(method, uri string) Request {
	r := newRequest()
	r.Method = method
	r.URI = uri
	return r
}
