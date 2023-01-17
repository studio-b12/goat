package gurlfile

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRequest(t *testing.T) {

	t.Run("simple", func(t *testing.T) {
		const raw = `
GET https://example.com
		`

		exp := newRequest()
		exp.raw = strings.TrimSpace(raw)
		exp.Method = "GET"
		exp.URI = "https://example.com"

		res, err := testCtx().parseRequest(raw, nil)
		assert.Nil(t, err)
		assert.Equal(t, exp, res)
	})

	t.Run("headers", func(t *testing.T) {
		const raw = `
GET https://example.com
Content-Type: application/json
X-XSRF-Token:	 some token
		`

		exp := newRequest()
		exp.raw = strings.TrimSpace(raw)
		exp.Method = "GET"
		exp.URI = "https://example.com"
		exp.Header = http.Header{
			"Content-Type": []string{"application/json"},
			"X-Xsrf-Token": []string{"some token"},
		}

		res, err := testCtx().parseRequest(raw, nil)
		assert.Nil(t, err)
		assert.Equal(t, exp, res)
	})

	t.Run("body", func(t *testing.T) {
		const raw = `
GET https://example.com
{
	"hello": "world"
}
		`

		exp := newRequest()
		exp.raw = strings.TrimSpace(raw)
		exp.Method = "GET"
		exp.URI = "https://example.com"
		exp.Body = []byte("{\n\t\"hello\": \"world\"\n}")

		res, err := testCtx().parseRequest(raw, nil)
		assert.Nil(t, err)
		assert.Equal(t, exp, res)
	})

	t.Run("headers", func(t *testing.T) {
		const raw = `
GET https://example.com
Content-Type: application/json
X-XSRF-Token:	 some token
{
	"hello": "world"
}
		`

		exp := newRequest()
		exp.raw = strings.TrimSpace(raw)
		exp.Method = "GET"
		exp.URI = "https://example.com"
		exp.Header = http.Header{
			"Content-Type": []string{"application/json"},
			"X-Xsrf-Token": []string{"some token"},
		}
		exp.Body = []byte("{\n\t\"hello\": \"world\"\n}")

		res, err := testCtx().parseRequest(raw, nil)
		assert.Nil(t, err)
		assert.Equal(t, exp, res)
	})

	t.Run("options", func(t *testing.T) {
		const raw = `
GET https://example.com

[QueryParams]
page = 1
sortBy = "date"
filter = [1, 2, 3]
		`

		exp := newRequest()
		exp.raw = strings.TrimSpace(raw)
		exp.Method = "GET"
		exp.URI = "https://example.com"
		exp.QueryParams = map[string]any{
			"page":   int64(1),
			"sortBy": "date",
			"filter": []any{int64(1), int64(2), int64(3)},
		}

		res, err := testCtx().parseRequest(raw, nil)
		assert.Nil(t, err)
		assert.Equal(t, exp, res)
	})

	t.Run("script", func(t *testing.T) {
		const raw = `
GET https://example.com

assert(true);
var a = 1;
		`

		exp := newRequest()
		exp.raw = strings.TrimSpace(raw)
		exp.Method = "GET"
		exp.URI = "https://example.com"
		exp.Script = "assert(true);\nvar a = 1;"

		res, err := testCtx().parseRequest(raw, nil)
		assert.Nil(t, err)
		assert.Equal(t, exp, res)
	})

	t.Run("all-together", func(t *testing.T) {
		const raw = `
GET https://example.com
Content-Type: application/json
X-XSRF-Token:	 some token
{
	"hello": "world"
}

[QueryParams]
page = 1
sortBy = "date"
filter = [1, 2, 3]

assert(true);
var a = 1;
		`

		exp := newRequest()
		exp.raw = strings.TrimSpace(raw)
		exp.Method = "GET"
		exp.URI = "https://example.com"
		exp.Header = http.Header{
			"Content-Type": []string{"application/json"},
			"X-Xsrf-Token": []string{"some token"},
		}
		exp.Body = []byte("{\n\t\"hello\": \"world\"\n}")
		exp.QueryParams = map[string]any{
			"page":   int64(1),
			"sortBy": "date",
			"filter": []any{int64(1), int64(2), int64(3)},
		}
		exp.Script = "assert(true);\nvar a = 1;"

		res, err := testCtx().parseRequest(raw, nil)
		assert.Nil(t, err)
		assert.Equal(t, exp, res)
	})

	t.Run("params", func(t *testing.T) {
		const raw = `
GET https://example.com/{{.Path}}
Content-Type: {{.ContentType}}
X-XSRF-Token:	 some token
{
	"hello": "{{.Data}}"
}

[QueryParams]
page = {{.Page}}
sortBy = "date"
filter = [1, 2, 3]

assert(true);
var a = {{.A}};
		`

		exp := newRequest()
		exp.raw = strings.TrimSpace(raw)
		exp.Method = "GET"
		exp.URI = "https://example.com/{{.Path}}"
		exp.Header = http.Header{
			"Content-Type": []string{"{{.ContentType}}"},
			"X-Xsrf-Token": []string{"some token"},
		}
		exp.Body = []byte("{\n\t\"hello\": \"{{.Data}}\"\n}")
		exp.QueryParams = map[string]any{
			"page":   "{{.Page}}",
			"sortBy": "date",
			"filter": []any{int64(1), int64(2), int64(3)},
		}
		exp.Script = "assert(true);\nvar a = {{.A}};"

		res, err := testCtx().parseRequest(raw, nil)
		assert.Nil(t, err)
		assert.Equal(t, exp, res)
	})

	t.Run("error-empty", func(t *testing.T) {
		const raw = `

		`

		exp := Request{}

		res, err := testCtx().parseRequest(raw, nil)
		assert.ErrorIs(t, err, ErrEmptyRequest)
		assert.Equal(t, exp, res)
	})
}

func TestRemoveComments(t *testing.T) {
	var raw, exp, res string

	raw = `
some text // some comments
// another comment
text//close to comment
https://example.com

test
	`

	exp = `
some text

text//close to comment
https://example.com

test
	`

	res = removeComments(raw)
	assert.Equal(t, exp, res)
}

func testCtx() context {
	return context{}
}
