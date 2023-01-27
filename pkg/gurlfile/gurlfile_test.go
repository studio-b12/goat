package gurlfile

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToHttpRequest(t *testing.T) {
	req := newRequest()
	req.URI = "https://example.com/api/v1"
	req.Method = "POST"
	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-foo", "bar")
	req.QueryParams = map[string]any{
		"foo":  "bar",
		"bazz": 2,
		"arr":  []any{3, 4},
	}

	httpReq, err := req.ToHttpRequest()
	assert.Nil(t, err, err)
	assert.Equal(t, "POST", httpReq.Method)
	assert.Equal(t, "https://example.com/api/v1?arr=3&arr=4&bazz=2&foo=bar", httpReq.URL.String())
	assert.Equal(t, http.Header{
		"Content-Type": []string{"application/json"},
		"X-Foo":        []string{"bar"},
	}, httpReq.Header)
}

func TestMerge(t *testing.T) {
	getA := func() Gurlfile {
		return Gurlfile{
			Setup: []Request{
				testRequest("A", "1"),
				testRequest("A", "2"),
				testRequest("A", "3"),
			},
			Tests: []Request{
				testRequest("A", "4"),
				testRequest("A", "5"),
			},
			Teardown: []Request{},
		}
	}

	getB := func() Gurlfile {
		return Gurlfile{
			Setup: []Request{
				testRequest("B", "1"),
				testRequest("B", "2"),
			},
			SetupEach: []Request{
				testRequest("B", "3"),
			},
			Tests: []Request{},
			Teardown: []Request{
				testRequest("B", "4"),
				testRequest("B", "5"),
			},
			TeardownEach: []Request{
				testRequest("B", "6"),
			},
		}
	}

	t.Run("a-into-b", func(t *testing.T) {
		a := getA()
		b := getB()

		a.Merge(b)

		assert.Equal(t, getB(), b)
		assert.Equal(t, Gurlfile{
			Setup: []Request{
				testRequest("A", "1"),
				testRequest("A", "2"),
				testRequest("A", "3"),
				testRequest("B", "1"),
				testRequest("B", "2"),
			},
			SetupEach: []Request{
				testRequest("B", "3"),
			},
			Tests: []Request{
				testRequest("A", "4"),
				testRequest("A", "5"),
			},
			Teardown: []Request{
				testRequest("B", "4"),
				testRequest("B", "5"),
			},
			TeardownEach: []Request{
				testRequest("B", "6"),
			},
		}, a)
	})

	t.Run("b-into-a", func(t *testing.T) {
		a := getA()
		b := getB()

		b.Merge(a)

		assert.Equal(t, getA(), a)
		assert.Equal(t, Gurlfile{
			Setup: []Request{
				testRequest("B", "1"),
				testRequest("B", "2"),
				testRequest("A", "1"),
				testRequest("A", "2"),
				testRequest("A", "3"),
			},
			SetupEach: []Request{
				testRequest("B", "3"),
			},
			Tests: []Request{
				testRequest("A", "4"),
				testRequest("A", "5"),
			},
			Teardown: []Request{
				testRequest("B", "4"),
				testRequest("B", "5"),
			},
			TeardownEach: []Request{
				testRequest("B", "6"),
			},
		}, b)
	})
}

func TestParseWithParams(t *testing.T) {
	getReq := func() Request {
		r := newRequest()
		r.Method = "{{.method}}"
		r.URI = "{{.instance}}/api/v1/login"
		r.Header.Set("Content-Type", "{{.contentType}}")
		r.Header.Set("Authorization", "bearer {{.token}}")
		r.QueryParams = map[string]any{"page": "{{.page}}"}
		r.Options = map[string]any{"condition": "{{.condition}}"}
		r.Body = []byte(`{"username": "{{.creds.username}}", "password": "{{.creds.password}}"}`)
		r.Script = `var foo = "{{.foo}}"`
		return r
	}

	t.Run("general", func(t *testing.T) {
		params := map[string]any{
			"method":      "GET",
			"instance":    "https://example.com",
			"contentType": "application/json",
			"token":       "some-token",
			"page":        2,
			"condition":   true,
			"creds": map[string]any{
				"username": "some-username",
				"password": "some-password",
			},
			"foo": "bar",
		}

		r := getReq()
		err := r.ParseWithParams(params)
		assert.Nil(t, err, err)
	})
}

// --- Helpers ---

func testRequest(method, uri string) Request {
	r := newRequest()
	r.Method = method
	r.URI = uri
	return r
}
