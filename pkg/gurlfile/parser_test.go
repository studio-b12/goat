package gurlfile

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: Add more unit tests

func TestParse_Simple(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		const raw = `GET https://example.com`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t, 1, len(res.Tests))
		assert.Equal(t, "GET", res.Tests[0].Method)
		assert.Equal(t, "https://example.com", res.Tests[0].URI)
	})

	t.Run("multi", func(t *testing.T) {
		const raw = `
GET https://example1.com

---

POST https://example2.com
---
LOGIN https://example3.com
-----------------------
		
CHECK https://example4.com
		
---
		`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)

		assert.Equal(t, 4, len(res.Tests))

		assert.Equal(t, "GET", res.Tests[0].Method)
		assert.Equal(t, "https://example1.com", res.Tests[0].URI)

		assert.Equal(t, "POST", res.Tests[1].Method)
		assert.Equal(t, "https://example2.com", res.Tests[1].URI)

		assert.Equal(t, "LOGIN", res.Tests[2].Method)
		assert.Equal(t, "https://example3.com", res.Tests[2].URI)

		assert.Equal(t, "CHECK", res.Tests[3].Method)
		assert.Equal(t, "https://example4.com", res.Tests[3].URI)
	})
}

func TestParse_Blocks(t *testing.T) {
	t.Run("single-single-block", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[Header]
Key-1: value 1
key-2: value 2
		
		`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t, 1, len(res.Tests))
		assert.Equal(t, "GET", res.Tests[0].Method)
		assert.Equal(t, "https://example.com", res.Tests[0].URI)
		assert.Equal(t, http.Header{
			"Key-1": []string{"value 1"},
			"Key-2": []string{"value 2"},
		}, res.Tests[0].Header)
	})

	t.Run("single-multi-block", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[Header]
Key-1: value 1
key-2: value 2

[Body]
some
body

[queryparams]
keyInt = 2
keyString = "some string"
		
		`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t, 1, len(res.Tests))
		assert.Equal(t, "GET", res.Tests[0].Method)
		assert.Equal(t, "https://example.com", res.Tests[0].URI)
		assert.Equal(t, http.Header{
			"Key-1": []string{"value 1"},
			"Key-2": []string{"value 2"},
		}, res.Tests[0].Header)
		assert.Equal(t, []byte("some\nbody\n"), res.Tests[0].Body)
		assert.Equal(t, map[string]any{
			"keyInt":    int64(2),
			"keyString": "some string",
		}, res.Tests[0].QueryParams)
	})

	t.Run("single-invalidblockheader", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[invalidblock]
Key-1: value 1
key-2: value 2
		
		`

		p := stringParser(raw)
		_, err := p.Parse()

		assert.ErrorIs(t, err, ErrInvalidBlockHeader, err)
	})

	t.Run("single-emptyblockheader", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[]
Key-1: value 1
key-2: value 2
		
		`

		p := stringParser(raw)
		_, err := p.Parse()

		assert.ErrorIs(t, err, ErrInvalidBlockHeader)
	})

	t.Run("single-openblockheader", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[QueryParams
Key-1: value 1
key-2: value 2
		
		`

		p := stringParser(raw)
		_, err := p.Parse()

		assert.ErrorIs(t, err, ErrInvalidBlockHeader, err)
	})
}

func TestParse_BlockHeaders(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[Header]
		
		`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t, http.Header{}, res.Tests[0].Header)
	})

	t.Run("values", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[Header]
key: value
key-2:  value 2
Some-Key-3: 		some value 3
SOME_KEY_4: 		§$%&/()=!§

multiple-1: value 1
multiple-1: value 2

		`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t, http.Header{
			"Key":        []string{"value"},
			"Key-2":      []string{"value 2"},
			"Some-Key-3": []string{"some value 3"},
			"Some_key_4": []string{"§$%&/()=!§"},
			"Multiple-1": []string{"value 1", "value 2"},
		}, res.Tests[0].Header)
	})

	t.Run("no-separator", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[Header]
invalid
		
		`

		p := stringParser(raw)
		_, err := p.Parse()

		assert.ErrorIs(t, err, ErrInvalidHeaderSeparator, err)
	})

	t.Run("invalid-key-format", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[Header]
some key: value
		
		`

		p := stringParser(raw)
		_, err := p.Parse()

		assert.ErrorIs(t, err, ErrInvalidHeaderSeparator, err)
	})

	t.Run("no-value", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[Header]
some-key:
some-key-2:
		
		`

		p := stringParser(raw)
		_, err := p.Parse()

		assert.ErrorIs(t, err, ErrNoHeaderValue, err)
	})
}

func stringParser(raw string) *Parser {
	return NewParser(strings.NewReader(raw))
}
