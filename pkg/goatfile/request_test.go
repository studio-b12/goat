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
