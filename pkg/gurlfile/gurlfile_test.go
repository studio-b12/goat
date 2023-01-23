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
