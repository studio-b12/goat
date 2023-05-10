package goatfile

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

func TestParseWithParams(t *testing.T) {
	getReq := func() Request {
		r := newRequest()
		r.Method = "{{.method}}"
		r.URI = "{{.instance}}/api/v1/login"
		r.Header.Set("Content-Type", "{{.contentType}}")
		r.Header.Set("Authorization", "bearer {{.token}}")
		r.QueryParams = map[string]any{"page": "{{.page}}"}
		r.Options = map[string]any{"condition": "{{.condition}}"}
		r.Body = StringContent(`{"username": "{{.creds.username}}", "password": "{{.creds.password}}"}`)
		r.Script = StringContent(`var foo = "{{.foo}}"`)
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

func TestMerge_request(t *testing.T) {
	t.Run("headers", func(t *testing.T) {
		def := newRequest()
		def.Header.Add("foo", "bar")
		def.Header.Add("hello", "world")

		req := newRequest()
		req.Header.Add("bazz", "fuzz")
		req.Header.Add("hello", "moon")

		req.Merge(&def)

		assert.Equal(t, "bar",
			req.Header.Get("foo"))
		assert.Equal(t, "fuzz",
			req.Header.Get("bazz"))
		assert.Equal(t, "moon",
			req.Header.Get("hello"))
	})

	t.Run("queryParams", func(t *testing.T) {
		def := newRequest()
		def.QueryParams = make(map[string]any)
		def.QueryParams["foo"] = "bar"
		def.QueryParams["hello"] = "world"

		req := newRequest()
		req.QueryParams = make(map[string]any)
		req.QueryParams["bazz"] = "fuzz"
		req.QueryParams["hello"] = "moon"

		req.Merge(&def)

		assert.Equal(t, "bar",
			req.QueryParams["foo"])
		assert.Equal(t, "fuzz",
			req.QueryParams["bazz"])
		assert.Equal(t, "moon",
			req.QueryParams["hello"])
	})

	t.Run("options", func(t *testing.T) {
		def := newRequest()
		def.Options = make(map[string]any)
		def.Options["foo"] = "bar"
		def.Options["hello"] = "world"

		req := newRequest()
		req.Options = make(map[string]any)
		req.Options["bazz"] = "fuzz"
		req.Options["hello"] = "moon"

		req.Merge(&def)

		assert.Equal(t, "bar",
			req.Options["foo"])
		assert.Equal(t, "fuzz",
			req.Options["bazz"])
		assert.Equal(t, "moon",
			req.Options["hello"])
	})

	t.Run("body", func(t *testing.T) {
		def := newRequest()
		def.Body = StringContent("foo bar")
		req := newRequest()
		req.Merge(&def)
		assert.Equal(t, StringContent("foo bar"), req.Body)

		def = newRequest()
		req = newRequest()
		req.Body = StringContent("foo bar")
		req.Merge(&def)
		assert.Equal(t, StringContent("foo bar"), req.Body)

		def = newRequest()
		def.Body = StringContent("hello world")
		req = newRequest()
		req.Body = StringContent("foo bar")
		req.Merge(&def)
		assert.Equal(t, StringContent("foo bar"), req.Body)

		def = newRequest()
		req = newRequest()
		req.Merge(&def)
		assert.Equal(t, NoContent{}, req.Body)
	})

	t.Run("preScript", func(t *testing.T) {
		def := newRequest()
		def.PreScript = StringContent("foo bar")
		req := newRequest()
		req.Merge(&def)
		assert.Equal(t, StringContent("foo bar"), req.PreScript)

		def = newRequest()
		req = newRequest()
		req.PreScript = StringContent("foo bar")
		req.Merge(&def)
		assert.Equal(t, StringContent("foo bar"), req.PreScript)

		def = newRequest()
		def.PreScript = StringContent("hello world")
		req = newRequest()
		req.PreScript = StringContent("foo bar")
		req.Merge(&def)
		assert.Equal(t, StringContent("foo bar"), req.PreScript)

		def = newRequest()
		req = newRequest()
		req.Merge(&def)
		assert.Equal(t, NoContent{}, req.PreScript)
	})

	t.Run("script", func(t *testing.T) {
		def := newRequest()
		def.Script = StringContent("foo bar")
		req := newRequest()
		req.Merge(&def)
		assert.Equal(t, StringContent("foo bar"), req.Script)

		def = newRequest()
		req = newRequest()
		req.Script = StringContent("foo bar")
		req.Merge(&def)
		assert.Equal(t, StringContent("foo bar"), req.Script)

		def = newRequest()
		def.Script = StringContent("hello world")
		req = newRequest()
		req.Script = StringContent("foo bar")
		req.Merge(&def)
		assert.Equal(t, StringContent("foo bar"), req.Script)

		def = newRequest()
		req = newRequest()
		req.Merge(&def)
		assert.Equal(t, NoContent{}, req.Script)
	})
}
