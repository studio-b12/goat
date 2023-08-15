package goatfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	def := newRequest()
	def.Header.Add("foo", "bar")

	getA := func() Goatfile {
		return Goatfile{
			Defaults: &def,
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
			Tests: []Action{},
			Teardown: []Action{
				testRequest("B", "4"),
				testRequest("B", "5"),
			},
		}
	}

	t.Run("a-into-b", func(t *testing.T) {
		a := getA()
		b := getB()

		a.Merge(b)

		assert.Equal(t, getB(), b)
		assert.Equal(t, Goatfile{
			Defaults: &def,
			Setup: []Action{
				testRequest("A", "1"),
				testRequest("A", "2"),
				testRequest("A", "3"),
				testRequest("B", "1"),
				testRequest("B", "2"),
			},
			Tests: []Action{
				testRequest("A", "4"),
				testRequest("A", "5"),
			},
			Teardown: []Action{
				testRequest("B", "4"),
				testRequest("B", "5"),
			},
		}, a)
	})

	t.Run("b-into-a", func(t *testing.T) {
		a := getA()
		b := getB()

		b.Merge(a)

		assert.Equal(t, getA(), a)
		assert.Equal(t, Goatfile{
			Defaults: &def,
			Setup: []Action{
				testRequest("B", "1"),
				testRequest("B", "2"),
				testRequest("A", "1"),
				testRequest("A", "2"),
				testRequest("A", "3"),
			},
			Tests: []Action{
				testRequest("A", "4"),
				testRequest("A", "5"),
			},
			Teardown: []Action{
				testRequest("B", "4"),
				testRequest("B", "5"),
			},
		}, b)
	})
}

// --- Helpers ---

func testRequest(method, uri string, opt ...int) Request {
	r := newRequest()
	r.Method = method
	r.URI = uri

	if len(opt) > 0 {
		r.PosLine = opt[0]
	}

	return r
}

func testRequestWithPath(method, uri string, path string, opt ...int) Request {
	r := newRequest()
	r.Method = method
	r.URI = uri
	r.Path = path

	if len(opt) > 0 {
		r.PosLine = opt[0]
	}

	return r
}
