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

func TestParse_BlockRaw(t *testing.T) {
	t.Run("unescaped-empty", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[Body]
`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t, []byte(nil), res.Tests[0].Body)
	})

	t.Run("unescaped-EOF", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[Body]
some body content
some more content
`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t,
			[]byte("some body content\nsome more content\n"),
			res.Tests[0].Body)
	})

	t.Run("unescaped-newblock", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[Body]
some body content
some more content

[QueryParams]
`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t,
			[]byte("some body content\nsome more content\n"),
			res.Tests[0].Body)
	})

	t.Run("unescaped-newrequest", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[Body]
some body content
some more content

---
`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t,
			[]byte("some body content\nsome more content\n"),
			res.Tests[0].Body)
	})

	t.Run("unescaped-section", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[Body]
some body content
some more content

### tests
`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t,
			[]byte("some body content\nsome more content\n"),
			res.Tests[0].Body)
	})

	t.Run("escaped-empty", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[Body]
´´´
´´´
`

		p := stringParser(swapTicks(raw))
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t, []byte(nil), res.Tests[0].Body)
	})

	t.Run("escaped-EOF", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[Body]
´´´
some body content
some more content
´´´
`

		p := stringParser(swapTicks(raw))
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t,
			[]byte("some body content\nsome more content\n"),
			res.Tests[0].Body)
	})

	t.Run("escaped-newblock", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[Body]
´´´
some body content

[QueryParams]
some more content
´´´

[QueryParams]
`

		p := stringParser(swapTicks(raw))
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t,
			[]byte("some body content\n\n[QueryParams]\nsome more content\n"),
			res.Tests[0].Body)
	})

	t.Run("escaped-newrequest", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[Body]
´´´
some body content

---

some more content
´´´

---
`

		p := stringParser(swapTicks(raw))
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t,
			[]byte("some body content\n\n---\n\nsome more content\n"),
			res.Tests[0].Body)
	})

	t.Run("escaped-section", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[Body]
´´´
some body content

### setup

some more content
´´´

### tests
`

		p := stringParser(swapTicks(raw))
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t,
			[]byte("some body content\n\n### setup\n\nsome more content\n"),
			res.Tests[0].Body)
	})

	t.Run("escaped-open", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[Body]
´´´
some body content

---
`

		p := stringParser(swapTicks(raw))
		_, err := p.Parse()

		assert.ErrorIs(t, err, ErrOpenEscapeBlock, err)
	})
}

func TestParse_BlockValues(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[QueryParams]
		
		`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t, map[string]any{}, res.Tests[0].QueryParams)
	})

	t.Run("value-strings", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[QueryParams]
string1 = "some string 1"
string2 =     "some string 2"
string3 = 		"some string 3" 
		`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t, map[string]any{
			"string1": "some string 1",
			"string2": "some string 2",
			"string3": "some string 3",
		}, res.Tests[0].QueryParams)
	})

	t.Run("value-integer", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[QueryParams]
int1 = 1
int2 = 1_000
int3 = -123
		`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t, map[string]any{
			"int1": int64(1),
			"int2": int64(1000),
			"int3": int64(-123),
		}, res.Tests[0].QueryParams)
	})

	t.Run("value-float", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[QueryParams]
float1 = 1.234
float2 = 1_000.234
float3 = 0.12
float4 = -12.34
		`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t, map[string]any{
			"float1": float64(1.234),
			"float2": float64(1000.234),
			"float3": float64(0.12),
			"float4": float64(-12.34),
		}, res.Tests[0].QueryParams)
	})

	t.Run("value-boolean", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[QueryParams]
bool1 = true
bool2 = false
		`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t, map[string]any{
			"bool1": true,
			"bool2": false,
		}, res.Tests[0].QueryParams)
	})

	t.Run("value-array", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[QueryParams]
arrayEmpty1 = []
arrayEmpty2 = [  	]

arrayString1 = ["some string"]
arrayString2 = ["some string", "another string","and another one"]

arrayInt1 = [1]
arrayInt2 = [1, 2,-3,	4_000]

arrayFloat1 = [1.23]
arrayFloat2 = [1.0, -1.1,1.234]

arrayMixed = ["a string", 2, 3.456, true]

arrayNested = [[1,2], [[true, false], "foo"]]

arrayMultiline = [
	"foo",
	"bar"
]

arrayLeadingComma1 = [true, false,]
arrayLeadingComma2 = [
	true, 
	false,
]
		`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t, map[string]any{
			"arrayEmpty1":        []any(nil),
			"arrayEmpty2":        []any(nil),
			"arrayString1":       []any{"some string"},
			"arrayString2":       []any{"some string", "another string", "and another one"},
			"arrayInt1":          []any{int64(1)},
			"arrayInt2":          []any{int64(1), int64(2), int64(-3), int64(4_000)},
			"arrayFloat1":        []any{1.23},
			"arrayFloat2":        []any{1.0, -1.1, 1.234},
			"arrayMixed":         []any{"a string", int64(2), 3.456, true},
			"arrayNested":        []any{[]any{int64(1), int64(2)}, []any{[]any{true, false}, "foo"}},
			"arrayMultiline":     []any{"foo", "bar"},
			"arrayLeadingComma1": []any{true, false},
			"arrayLeadingComma2": []any{true, false},
		}, res.Tests[0].QueryParams)
	})

	t.Run("value-invalid-entry", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[QueryParams]
invalid
		`

		p := stringParser(raw)
		_, err := p.Parse()

		assert.ErrorIs(t, err, ErrInvalidBlockEntryAssignment, err)
	})

	t.Run("value-invalid-assignment", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[QueryParams]
invalid = 
		`

		p := stringParser(raw)
		_, err := p.Parse()

		assert.ErrorIs(t, err, ErrInvalidToken, err)
	})

	t.Run("value-invalid-string", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[QueryParams]
invalid = "
		`

		p := stringParser(raw)
		_, err := p.Parse()

		assert.ErrorIs(t, err, ErrInvalidToken, err)
	})

	t.Run("value-invalid-array-1", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[QueryParams]
invalid = [
		`

		p := stringParser(raw)
		_, err := p.Parse()

		assert.ErrorIs(t, err, ErrInvalidToken, err)
	})

	t.Run("value-invalid-array-2", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[QueryParams]
invalid = [1, 2
		`

		p := stringParser(raw)
		_, err := p.Parse()

		assert.ErrorIs(t, err, ErrInvalidToken, err)
	})
}

// --- Helpers --------------------------------------------

func stringParser(raw string) *Parser {
	return NewParser(strings.NewReader(raw))
}

func swapTicks(v string) string {
	return strings.ReplaceAll(v, "´", "`")
}
