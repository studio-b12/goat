package goatfile

import (
	"github.com/studio-b12/goat/pkg/goatfile/ast"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestFromAst(t *testing.T) {
	astR := ast.Request{
		Head: ast.RequestHead{
			Method: "GET",
			Url:    "https://foo.bar",
		},
		Blocks: []ast.RequestBlock{
			ast.RequestOptions{ast.KVList[any]{ast.KV[any]{Key: "a", Value: "b"}}},
			ast.RequestBody{ast.TextBlock{"body stuff"}},
			ast.RequestScript{ast.TextBlock{"script stuff"}},
		},
	}

	gf, err := RequestFromAst(&astR, "somepath")
	assert.Nil(t, err, err)

	assert.Equal(t, "somepath", gf.Path)
	assert.Equal(t, "GET", gf.Method)
	assert.Equal(t, "https://foo.bar", gf.URI)
	assert.Equal(t, "b", gf.Options["a"])
	assert.Equal(t, StringContent("body stuff"), gf.Body)
	assert.Equal(t, StringContent("script stuff"), gf.Script)

}

func TestPartialRequestFromAst(t *testing.T) {
	astR := ast.PartialRequest{
		Blocks: []ast.RequestBlock{
			ast.RequestOptions{ast.KVList[any]{ast.KV[any]{Key: "a", Value: "b"}}},
			ast.RequestBody{ast.TextBlock{"body stuff"}},
			ast.RequestScript{ast.TextBlock{"script stuff"}},
		},
	}

	gf, err := PartialRequestFromAst(astR, "somepath")
	assert.Nil(t, err, err)

	assert.Equal(t, "somepath", gf.Path)
	assert.Equal(t, "", gf.Method)
	assert.Equal(t, "", gf.URI)
	assert.Equal(t, "b", gf.Options["a"])
	assert.Equal(t, StringContent("body stuff"), gf.Body)
	assert.Equal(t, StringContent("script stuff"), gf.Script)

}

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

func TestPreSubstituteWithParams(t *testing.T) {
	getReq := func() Request {
		r := newRequest()
		r.URI = "{{.instance}}/api/v1/login"
		r.Header.Set("Content-Type", "{{.contentType}}")
		r.Header.Set("Authorization", "bearer {{.token}}")
		r.QueryParams = map[string]any{"page": "{{.page}}"}
		r.Options = map[string]any{"condition": "{{.condition}}"}
		r.Body = StringContent(`{"username": "{{.creds.username}}", "password": "{{.creds.password}}"}`)
		r.PreScript = StringContent(`var bar = "{{.bar}}"`)
		r.Script = StringContent(`var foo = "{{.foo}}"`)
		return r
	}

	t.Run("general", func(t *testing.T) {
		params := map[string]any{
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
			"bar": "bazz",
		}

		r := getReq()
		err := r.PreSubstitudeWithParams(params)
		assert.Nil(t, err, err)

		assert.Equal(t,
			"{{.instance}}/api/v1/login",
			r.URI)
		assert.Equal(t,
			"{{.contentType}}",
			r.Header.Get("Content-Type"))
		assert.Equal(t,
			"bearer {{.token}}",
			r.Header.Get("Authorization"))
		assert.Equal(t,
			"{{.page}}",
			r.QueryParams["page"])
		assert.Equal(t,
			"{{.condition}}",
			r.Options["condition"])
		assert.Equal(t,
			StringContent(`{"username": "{{.creds.username}}", "password": "{{.creds.password}}"}`),
			r.Body)
		assert.Equal(t,
			StringContent(`var bar = "bazz"`),
			r.PreScript)
		assert.Equal(t,
			StringContent(`var foo = "{{.foo}}"`),
			r.Script)
	})
}

func TestSubstituteWithParams(t *testing.T) {
	getReq := func() Request {
		r := newRequest()
		r.URI = "{{.instance}}/api/v1/login"
		r.Header.Set("Content-Type", "{{.contentType}}")
		r.Header.Set("Authorization", "bearer {{.token}}")
		r.QueryParams = map[string]any{"page": "{{.page}}"}
		r.Options = map[string]any{"condition": "{{.condition}}"}
		r.Body = StringContent(`{"username": "{{.creds.username}}", "password": "{{.creds.password}}"}`)
		r.PreScript = StringContent(`var bar = "{{.bar}}"`)
		r.Script = StringContent(`var foo = "{{.foo}}"`)
		return r
	}

	t.Run("general", func(t *testing.T) {
		params := map[string]any{
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
			"bar": "bazz",
		}

		r := getReq()
		err := r.SubstitudeWithParams(params)
		assert.Nil(t, err, err)

		assert.Equal(t,
			"https://example.com/api/v1/login",
			r.URI)
		assert.Equal(t,
			"application/json",
			r.Header.Get("Content-Type"))
		assert.Equal(t,
			"bearer some-token",
			r.Header.Get("Authorization"))
		assert.Equal(t,
			"2",
			r.QueryParams["page"])
		assert.Equal(t,
			"true",
			r.Options["condition"])
		assert.Equal(t,
			StringContent(`{"username": "some-username", "password": "some-password"}`),
			r.Body)
		assert.Equal(t,
			StringContent(`var bar = "{{.bar}}"`),
			r.PreScript)
		assert.Equal(t,
			StringContent(`var foo = "bar"`),
			r.Script)
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
