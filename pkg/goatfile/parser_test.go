package goatfile

import (
	"github.com/studio-b12/goat/pkg/goatfile/ast"
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
		assert.Equal(t, 1, len(res.Actions))
		assert.Equal(t, "GET", res.Actions[0].(*ast.Request).Head.Method)
		assert.Equal(t, "https://example.com", res.Actions[0].(*ast.Request).Head.Url)
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

[Body]
abc
		
---

CHECK https://example5.com

[Body]
abc
		
------
		`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)

		assert.Equal(t, 5, len(res.Actions))

		assert.Equal(t, "GET", res.Actions[0].(*ast.Request).Head.Method)
		assert.Equal(t, "https://example1.com", res.Actions[0].(*ast.Request).Head.Url)

		assert.Equal(t, "POST", res.Actions[1].(*ast.Request).Head.Method)
		assert.Equal(t, "https://example2.com", res.Actions[1].(*ast.Request).Head.Url)

		assert.Equal(t, "LOGIN", res.Actions[2].(*ast.Request).Head.Method)
		assert.Equal(t, "https://example3.com", res.Actions[2].(*ast.Request).Head.Url)

		assert.Equal(t, "CHECK", res.Actions[3].(*ast.Request).Head.Method)
		assert.Equal(t, "https://example4.com", res.Actions[3].(*ast.Request).Head.Url)

		assert.Equal(t, "CHECK", res.Actions[4].(*ast.Request).Head.Method)
		assert.Equal(t, "https://example5.com", res.Actions[4].(*ast.Request).Head.Url)
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
		assert.Equal(t, 1, len(res.Actions))
		assert.Equal(t, "GET", res.Actions[0].(*ast.Request).Head.Method)
		assert.Equal(t, "https://example.com", res.Actions[0].(*ast.Request).Head.Url)
		assert.Equal(t, ast.RequestHeader{
			"Key-1": []string{"value 1"},
			"key-2": []string{"value 2"},
		}, res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestHeader))
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

[Auth]
username = "foo"
password = "{{.creds.password}}"
		
		`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t, 1, len(res.Actions))
		assert.Equal(t, "GET", res.Actions[0].(*ast.Request).Head.Method)
		assert.Equal(t, "https://example.com", res.Actions[0].(*ast.Request).Head.Url)
		assert.Equal(t, ast.RequestHeader{
			"Key-1": []string{"value 1"},
			"key-2": []string{"value 2"},
		}, res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestHeader))
		assert.Equal(t, ast.TextBlock("some\nbody\n"), res.Actions[0].(*ast.Request).Blocks[1].(ast.RequestBody))
		assert.Equal(t, ast.RequestQueryParams{
			"keyInt":    int64(2),
			"keyString": "some string",
		}, res.Actions[0].(*ast.Request).Blocks[2].(ast.RequestQueryParams))
		assert.Equal(t, ast.RequestAuth{
			"username": "foo",
			"password": "{{.creds.password}}",
		}, res.Actions[0].(*ast.Request).Blocks[3].(ast.RequestAuth))
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
		assert.Equal(t, ast.RequestHeader{}, res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestHeader))
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
		assert.Equal(t, ast.RequestHeader{
			"key":        []string{"value"},
			"key-2":      []string{"value 2"},
			"Some-Key-3": []string{"some value 3"},
			"SOME_KEY_4": []string{"§$%&/()=!§"},
			"multiple-1": []string{"value 1", "value 2"},
		}, res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestHeader))
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
		assert.Equal(t, ast.NoContent{}, res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestBody))
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
			ast.TextBlock("some body content\nsome more content\n"),
			res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestBody))
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
			ast.TextBlock("some body content\nsome more content\n"),
			(res.Actions[0].(*ast.Request).Blocks[0]).(ast.RequestBody))
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
			ast.TextBlock("some body content\nsome more content\n"),
			res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestBody))
	})

	t.Run("unescaped-finaldelim", func(t *testing.T) {
		const raw = `

GET https://example.com

[Body]
some body content
some more content
---`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t,
			ast.TextBlock("some body content\nsome more content"),
			res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestBody))
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
			ast.TextBlock("some body content\nsome more content"),
			res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestBody))
	})

	t.Run("unescaped-logsection", func(t *testing.T) {
		const raw = `

GET https://example.com

[Body]
some body content
some more content
##### some log section
`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t,
			ast.TextBlock("some body content\nsome more content"),
			res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestBody))
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
		assert.Equal(t, ast.NoContent{}, res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestBody))
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
			ast.TextBlock("some body content\nsome more content\n"),
			res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestBody))
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
			ast.TextBlock("some body content\n\n[QueryParams]\nsome more content\n"),
			res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestBody))
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
			ast.TextBlock("some body content\n\n---\n\nsome more content\n"),
			res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestBody))
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
			ast.TextBlock("some body content\n\n### setup\n\nsome more content\n"),
			res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestBody))
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

	t.Run("script general", func(t *testing.T) {
		const raw = `

GET https://example.com

[Script]
assert(response.StatusCode == 200, "invalid status code");
var id = response.Body.id;

---

`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t,
			ast.TextBlock(`assert(response.StatusCode == 200, "invalid status code");`+
				"\nvar id = response.Body.id;\n"),
			res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestScript))
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
		assert.Equal(t, ast.RequestQueryParams{}, res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestQueryParams))
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
		assert.Equal(t, ast.RequestQueryParams{
			"string1": "some string 1",
			"string2": "some string 2",
			"string3": "some string 3",
		}, res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestQueryParams))
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
		assert.Equal(t, ast.RequestQueryParams{
			"int1": int64(1),
			"int2": int64(1000),
			"int3": int64(-123),
		}, res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestQueryParams))
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
		assert.Equal(t, ast.RequestQueryParams{
			"float1": float64(1.234),
			"float2": float64(1000.234),
			"float3": float64(0.12),
			"float4": float64(-12.34),
		}, res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestQueryParams))
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
		assert.Equal(t, ast.RequestQueryParams{
			"bool1": true,
			"bool2": false,
		}, res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestQueryParams))
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
		assert.Equal(t, ast.RequestQueryParams{
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
		}, res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestQueryParams))
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

func TestParse_Comments(t *testing.T) {
	t.Run("uri", func(t *testing.T) {
		const raw = `
// Some comment
   // Some comment
GET https://example.com //another comment
// comment
// heyo
			`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t, "GET", res.Actions[0].(*ast.Request).Head.Method)
		assert.Equal(t, "https://example.com", res.Actions[0].(*ast.Request).Head.Url)
	})

	t.Run("blocks", func(t *testing.T) {
		const raw = `
GET https://example.com

// some comment
[QueryParams] // block hader comment
key1 = "value" // another comment
key2 = 1.23 // comment
// in betweeny
arr = [ // comment
	1, // comment
	// another comment
	2 // comment
] // comment
			`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t, ast.RequestQueryParams{
			"key1": "value",
			"key2": 1.23,
			"arr":  []any{int64(1), int64(2)},
		}, res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestQueryParams))
	})

	t.Run("invlid-1", func(t *testing.T) {
		const raw = `
GET https://example.com

/
			`

		p := stringParser(raw)
		_, err := p.Parse()

		assert.ErrorIs(t, err, ErrInvalidToken)
	})

	t.Run("invlid", func(t *testing.T) {
		const raw = `
GET https://example.com / foo
			`

		p := stringParser(raw)
		_, err := p.Parse()

		assert.ErrorIs(t, err, ErrInvalidToken)
	})
}

func TestParse_Sections(t *testing.T) {
	t.Run("general", func(t *testing.T) {
		const raw = `
### Setup

GET https://example1.com
---
GET https://example2.com
---
GET https://example4.com

###   	tests

GET https://example5.com

---

GET https://example6.com

	### teardown

GET https://example7.com
---
GET https://example8.com

			`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)

		assert.Equal(t, "https://example1.com", res.Sections[0].(ast.SectionSetup).Actions[0].(*ast.Request).Head.Url)
		assert.Equal(t, "https://example2.com", res.Sections[0].(ast.SectionSetup).Actions[1].(*ast.Request).Head.Url)

		assert.Equal(t, "https://example5.com", res.Sections[1].(ast.SectionTests).Actions[0].(*ast.Request).Head.Url)
		assert.Equal(t, "https://example6.com", res.Sections[1].(ast.SectionTests).Actions[1].(*ast.Request).Head.Url)

		assert.Equal(t, "https://example7.com", res.Sections[2].(ast.SectionTeardown).Actions[0].(*ast.Request).Head.Url)
		assert.Equal(t, "https://example8.com", res.Sections[2].(ast.SectionTeardown).Actions[1].(*ast.Request).Head.Url)
	})

	t.Run("invalid-1", func(t *testing.T) {
		const raw = `
## Tests

GET https://example.com
			`

		p := stringParser(raw)
		_, err := p.Parse()

		assert.ErrorIs(t, err, ErrIllegalCharacter, err)
	})

	t.Run("invalid-2", func(t *testing.T) {
		const raw = `
###

GET https://example.com
			`

		p := stringParser(raw)
		_, err := p.Parse()

		assert.ErrorIs(t, err, ErrInvalidSection, err)
	})

	t.Run("invalid-3", func(t *testing.T) {
		const raw = `
### invalid-section

GET https://example.com
			`

		p := stringParser(raw)
		_, err := p.Parse()

		assert.ErrorIs(t, err, ErrInvalidSection, err)
	})

	t.Run("invalid-4", func(t *testing.T) {
		const raw = `
### Tests Invalid

GET https://example.com
			`

		p := stringParser(raw)
		_, err := p.Parse()

		assert.ErrorIs(t, err, ErrInvalidSection, err)
	})
}

func TestParse_Use(t *testing.T) {
	t.Run("general", func(t *testing.T) {
		const raw = `
use file1

use file2
use ../file3 // hey, a comment!

use "some file"

use 	  ../another/file
		`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t, []string{
			"file1",
			"file2",
			"../file3",
			"some file",
			"../another/file",
		}, importsToPaths(res.Imports))
	})

	t.Run("invalid-inclomplete", func(t *testing.T) {
		const raw = `
use
		`

		p := stringParser(raw)
		_, err := p.Parse()

		assert.ErrorIs(t, err, ErrInvalidStringLiteral, err)
	})

	t.Run("invalid-empty-1", func(t *testing.T) {
		const raw = `
use
		`

		p := stringParser(raw)
		_, err := p.Parse()

		assert.ErrorIs(t, err, ErrInvalidStringLiteral, err)
	})

	t.Run("invalid-empty-2", func(t *testing.T) {
		const raw = `
use ""
		`

		p := stringParser(raw)
		_, err := p.Parse()

		assert.ErrorIs(t, err, ErrEmptyUsePath, err)
	})

	t.Run("invalid-openstring", func(t *testing.T) {
		const raw = `
use "
		`

		p := stringParser(raw)
		_, err := p.Parse()

		assert.ErrorIs(t, err, ErrInvalidStringLiteral, err)
	})

	t.Run("invalid-keyword", func(t *testing.T) {
		const raw = `
use"test"
		`

		p := stringParser(raw)
		_, err := p.Parse()

		assert.ErrorIs(t, err, ErrInvalidStringLiteral, err)
	})
}

// --- Special Tests --------------------------------------

func TestParse_BlockTemplateValues(t *testing.T) {
	t.Run("variable-1", func(t *testing.T) {
		const raw = `
GET https://example.com

[Options]
someoption = {{.param}}
		`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t, ParameterValue(".param"), res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestOptions)["someoption"])
	})

	t.Run("variable-2", func(t *testing.T) {
		const raw = `
GET https://example.com

[Options]
someoption = {{ .param }}
		`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t, ParameterValue(" .param "), res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestOptions)["someoption"])
	})

	t.Run("wrapped", func(t *testing.T) {
		const raw = `
GET https://example.com

[Options]
someoption1 = {{ print {{.param1}} {{.param2}} }}
someoption2 = {{ print {{if .param1}}true{{else}}false{{end}} }}
		`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t, ParameterValue(" print {{.param1}} {{.param2}} "), res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestOptions)["someoption1"])
		assert.Equal(t, ParameterValue(" print {{if .param1}}true{{else}}false{{end}} "), res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestOptions)["someoption2"])
	})

	t.Run("instring-1", func(t *testing.T) {
		const raw = `
GET https://example.com

[Options]
someoption = {{ print "}}" }}
		`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t, ParameterValue(` print "}}" `), res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestOptions)["someoption"])
	})

	t.Run("instring-2", func(t *testing.T) {
		const raw = `
GET https://example.com

[Options]
someoption = {{ print ´}}´ }}
		`

		p := stringParser(swapTicks(raw))
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t, ParameterValue(" print `}}` "), res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestOptions)["someoption"])
	})

	t.Run("instring-wrapped", func(t *testing.T) {
		const raw = `
GET https://example.com

[Options]
someoption = {{ print {{ "}}" }} }}
		`

		p := stringParser(swapTicks(raw))
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t, ParameterValue(` print {{ "}}" }} `), res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestOptions)["someoption"])
	})
}

// TODO: This validation should now be handled on the AST
// See https://github.com/studio-b12/goat/issues/19
//func TestParseMultipleSectionsCheck(t *testing.T) {
//	t.Run("multiple-options", func(t *testing.T) {
//		const raw = `
//GET https://example.com
//
//[Options]
//someoption = "a"
//
//[Header]
//some: header
//
//[Options]
//anotheroption = "b"
//		`
//
//		p := stringParser(raw)
//		_, err := p.Parse()
//
//		assert.ErrorIs(t, err, ErrSectionDefinedMultiple, err)
//		var ewd errs.ErrorWithDetails
//		assert.True(t, errors.As(err, &ewd))
//		assert.Equal(t, fmt.Sprintf("[%s]:", optionNameOptions), ewd.Details.(string))
//	})
//
//	t.Run("multiple-header", func(t *testing.T) {
//		const raw = `
//GET https://example.com
//[Header]
//some: header
//
//[Options]
//someoption = "a"
//
//[Header]
//		`
//
//		p := stringParser(raw)
//		_, err := p.Parse()
//
//		assert.ErrorIs(t, err, ErrSectionDefinedMultiple, err)
//		var ewd errs.ErrorWithDetails
//		assert.True(t, errors.As(err, &ewd))
//		assert.Equal(t, fmt.Sprintf("[%s]:", optionNameHeader), ewd.Details.(string))
//	})
//
//	t.Run("multiple-script", func(t *testing.T) {
//		const raw = `
//GET https://example.com
//[Header]
//some: header
//
//[Options]
//someoption = "a"
//
//[Script]
//
//[Script]
//
//		`
//
//		p := stringParser(raw)
//		_, err := p.Parse()
//
//		assert.ErrorIs(t, err, ErrSectionDefinedMultiple, err)
//		var ewd errs.ErrorWithDetails
//		assert.True(t, errors.As(err, &ewd))
//		assert.Equal(t, fmt.Sprintf("[%s]:", optionNameScript), ewd.Details.(string))
//	})
//
//	t.Run("multiple-body", func(t *testing.T) {
//		const raw = `
//GET https://example.com
//[Header]
//some: header
//
//[Options]
//someoption = "a"
//
//[Body]
//foobar
//
//[Script]
//
//[Body]
//barbazz
//		`
//
//		p := stringParser(raw)
//		_, err := p.Parse()
//
//		assert.ErrorIs(t, err, ErrSectionDefinedMultiple, err)
//		var ewd errs.ErrorWithDetails
//		assert.True(t, errors.As(err, &ewd))
//		assert.Equal(t, fmt.Sprintf("[%s]:", optionNameBody), ewd.Details.(string))
//	})
//}

func TestLogSections(t *testing.T) {
	t.Run("general", func(t *testing.T) {
		const raw = `
##### Log section 1

GET https://example.com

[Script]
script stuff 1

---

##### Log section 2

GET https://example.com

[Script]
script stuff 2

##### Log section 3

GET https://example.com

[Script]
script stuff 3`

		p := stringParser(raw)
		gf, err := p.Parse()
		assert.Nil(t, err, err)
		assert.Equal(t, "Log section 1", gf.Actions[0].(ast.LogSection).Content)
		assert.Equal(t, ast.TextBlock("script stuff 1\n"), gf.Actions[1].(*ast.Request).Blocks[0].(ast.RequestScript))
		assert.Equal(t, "Log section 2", gf.Actions[2].(ast.LogSection).Content)
		assert.Equal(t, ast.TextBlock("script stuff 2\n"), gf.Actions[3].(*ast.Request).Blocks[0].(ast.RequestScript))
		assert.Equal(t, "Log section 3", gf.Actions[4].(ast.LogSection).Content)
		assert.Equal(t, ast.TextBlock("script stuff 3"), gf.Actions[5].(*ast.Request).Blocks[0].(ast.RequestScript))
	})
}

func TestDefaults(t *testing.T) {
	t.Run("general", func(t *testing.T) {
		const raw = `
### Defaults

// Some headers
[Header]
foo: bar
hello: world

[Options]
answer = 42 // some comment
some = "value"

[Auth]
username = "foo"
password = "bar"

[Body]
hello
world

[PreScript]
some pre script

[Script]
some script
`

		p := stringParser(raw)
		gf, err := p.Parse()
		assert.Nil(t, err, err)
		assert.Equal(t, []string{"bar"}, gf.Sections[0].(ast.SectionDefaults).Request.Blocks[0].(ast.RequestHeader)["foo"])
		assert.Equal(t, []string{"world"}, gf.Sections[0].(ast.SectionDefaults).Request.Blocks[0].(ast.RequestHeader)["hello"])
		assert.Equal(t, int64(42), gf.Sections[0].(ast.SectionDefaults).Request.Blocks[1].(ast.RequestOptions)["answer"])
		assert.Equal(t, "value", gf.Sections[0].(ast.SectionDefaults).Request.Blocks[1].(ast.RequestOptions)["some"])
		assert.Equal(t, "foo", gf.Sections[0].(ast.SectionDefaults).Request.Blocks[2].(ast.RequestAuth)["username"])
		assert.Equal(t, "bar", gf.Sections[0].(ast.SectionDefaults).Request.Blocks[2].(ast.RequestAuth)["password"])
		assert.Equal(t, ast.TextBlock("hello\nworld\n"), gf.Sections[0].(ast.SectionDefaults).Request.Blocks[3].(ast.RequestBody))
		assert.Equal(t, ast.TextBlock("some pre script\n"), gf.Sections[0].(ast.SectionDefaults).Request.Blocks[4].(ast.RequestPreScript))
		assert.Equal(t, ast.TextBlock("some script\n"), gf.Sections[0].(ast.SectionDefaults).Request.Blocks[5].(ast.RequestScript))
	})

	t.Run("with-following-section", func(t *testing.T) {
		const raw = `
### Defaults

// Some headers
[Header]
foo: bar
hello: world

[Script]
some script

### Tests

GET https://exmaple.com
---
GET https://exmaple.com
`

		p := stringParser(raw)
		gf, err := p.Parse()
		assert.Nil(t, err, err)
		assert.Equal(t, []string{"bar"}, gf.Sections[0].(ast.SectionDefaults).Request.Blocks[0].(ast.RequestHeader)["foo"])
		assert.Equal(t, []string{"world"}, gf.Sections[0].(ast.SectionDefaults).Request.Blocks[0].(ast.RequestHeader)["hello"])
		assert.Equal(t, ast.TextBlock("some script\n"), gf.Sections[0].(ast.SectionDefaults).Request.Blocks[1].(ast.RequestScript))
	})

	t.Run("with-following-splitter+section", func(t *testing.T) {
		const raw = `
### Defaults

// Some headers
[Header]
foo: bar
hello: world

[Script]
some script

---

### Tests

GET https://exmaple.com
---
GET https://exmaple.com
`

		p := stringParser(raw)
		gf, err := p.Parse()
		assert.Nil(t, err, err)
		assert.Equal(t, []string{"bar"}, gf.Sections[0].(ast.SectionDefaults).Request.Blocks[0].(ast.RequestHeader)["foo"])
		assert.Equal(t, []string{"world"}, gf.Sections[0].(ast.SectionDefaults).Request.Blocks[0].(ast.RequestHeader)["hello"])
		assert.Equal(t, ast.TextBlock("some script\n"), gf.Sections[0].(ast.SectionDefaults).Request.Blocks[1].(ast.RequestScript))
	})

	t.Run("with-following-splitter", func(t *testing.T) {
		const raw = `
### Defaults

// Some headers
[Header]
foo: bar
hello: world

[Script]
some script

---

GET https://exmaple.com
---
GET https://exmaple.com
`

		p := stringParser(raw)
		gf, err := p.Parse()
		assert.Nil(t, err, err)
		assert.Equal(t, []string{"bar"}, gf.Sections[0].(ast.SectionDefaults).Request.Blocks[0].(ast.RequestHeader)["foo"])
		assert.Equal(t, []string{"world"}, gf.Sections[0].(ast.SectionDefaults).Request.Blocks[0].(ast.RequestHeader)["hello"])
		assert.Equal(t, ast.TextBlock("some script\n"), gf.Sections[0].(ast.SectionDefaults).Request.Blocks[1].(ast.RequestScript))
	})

	t.Run("in-between", func(t *testing.T) {
		const raw = `
### Tests

GET https://exmaple.com

### Defaults

// Some headers
[Header]
foo: bar
hello: world

[Script]
some script

### Setup

GET https://exmaple.com
`

		p := stringParser(raw)
		gf, err := p.Parse()
		assert.Nil(t, err, err)
		assert.Equal(t, []string{"bar"}, gf.Sections[1].(ast.SectionDefaults).Request.Blocks[0].(ast.RequestHeader)["foo"])
		assert.Equal(t, []string{"world"}, gf.Sections[1].(ast.SectionDefaults).Request.Blocks[0].(ast.RequestHeader)["hello"])
		assert.Equal(t, ast.TextBlock("some script\n"), gf.Sections[1].(ast.SectionDefaults).Request.Blocks[1].(ast.RequestScript))
	})

	t.Run("multiple", func(t *testing.T) {
		const raw = `
### Defaults

// Some headers
[Header]
foo: bar
hello: world

[Script]
some script

### Defaults

// Some headers
[Header]
foo: bar
hello: moon

[Script]
some script 2
`

		p := stringParser(raw)
		gf, err := p.Parse()
		assert.Nil(t, err, err)
		assert.Equal(t, []string{"bar"}, gf.Sections[0].(ast.SectionDefaults).Request.Blocks[0].(ast.RequestHeader)["foo"])
		assert.Equal(t, []string{"moon"}, gf.Sections[1].(ast.SectionDefaults).Request.Blocks[0].(ast.RequestHeader)["hello"])
		assert.Equal(t, ast.TextBlock("some script\n"), gf.Sections[0].(ast.SectionDefaults).Request.Blocks[1].(ast.RequestScript))
		assert.Equal(t, ast.TextBlock("some script 2\n"), gf.Sections[1].(ast.SectionDefaults).Request.Blocks[1].(ast.RequestScript))
	})
}

func TestExecute(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		const raw = `
execute ../pathTo/someGoatfile
`

		p := stringParser(raw)
		gf, err := p.Parse()
		assert.Nil(t, err, err)

		assert.Equal(t, "../pathTo/someGoatfile", gf.Actions[0].(*ast.Execute).Path)
	})

	t.Run("multiple", func(t *testing.T) {
		const raw = `
execute ../pathTo/someGoatfile1
---
execute ../pathTo/someGoatfile2

execute ../pathTo/someGoatfile3
execute ../pathTo/someGoatfile4
`

		p := stringParser(raw)
		gf, err := p.Parse()
		assert.Nil(t, err, err)

		assert.Equal(t, 4, len(gf.Actions))
	})

	t.Run("multiple-in-section", func(t *testing.T) {
		const raw = `
### Tests
execute ../pathTo/someGoatfile1
---
execute ../pathTo/someGoatfile2

execute ../pathTo/someGoatfile3
execute ../pathTo/someGoatfile4
`

		p := stringParser(raw)
		gf, err := p.Parse()
		assert.Nil(t, err, err)

		assert.Equal(t, 4, len(gf.Sections[0].(ast.SectionTests).Actions))
	})

	t.Run("params", func(t *testing.T) {
		const raw = `
execute ../pathTo/someGoatfile (foo=1)

---

execute ../pathTo/someGoatfile (foo=1 bar="hello")

---

execute ../pathTo/someGoatfile (
	foo =    "hello"
	bar=2
	bazz= {{.someParam}}
)
`

		p := stringParser(raw)
		gf, err := p.Parse()
		assert.Nil(t, err, err)

		assert.Equal(t, "../pathTo/someGoatfile", gf.Actions[0].(*ast.Execute).Path)
		assert.Equal(t, ast.KV{
			"foo": int64(1),
		}, gf.Actions[0].(*ast.Execute).Parameters)

		assert.Equal(t, "../pathTo/someGoatfile", gf.Actions[1].(*ast.Execute).Path)
		assert.Equal(t, ast.KV{
			"foo": int64(1),
			"bar": "hello",
		}, gf.Actions[1].(*ast.Execute).Parameters)

		assert.Equal(t, "../pathTo/someGoatfile", gf.Actions[2].(*ast.Execute).Path)
		assert.Equal(t, ast.KV{
			"foo":  "hello",
			"bar":  int64(2),
			"bazz": ParameterValue(".someParam"),
		}, gf.Actions[2].(*ast.Execute).Parameters)
	})

	t.Run("return", func(t *testing.T) {
		const raw = `
execute ../pathTo/someGoatfile (foo=1) return (foo as bar)

---

execute ../pathTo/someGoatfile (foo=1 bar="hello") return (foo as bar bar as bazz)

---

execute ../pathTo/someGoatfile (
	foo =    "hello"
	bar=2
	bazz= {{.someParam}}
) return (
	foo as   bar
	bar    as bazz
)
`

		p := stringParser(raw)
		gf, err := p.Parse()
		assert.Nil(t, err, err)

		assert.Equal(t, "../pathTo/someGoatfile", gf.Actions[0].(*ast.Execute).Path)
		assert.Equal(t, ast.KV{
			"foo": int64(1),
		}, gf.Actions[0].(*ast.Execute).Parameters)
		assert.Equal(t, ast.Assignments{
			"foo": "bar",
		}, gf.Actions[0].(*ast.Execute).Returns)

		assert.Equal(t, "../pathTo/someGoatfile", gf.Actions[1].(*ast.Execute).Path)
		assert.Equal(t, ast.KV{
			"foo": int64(1),
			"bar": "hello",
		}, gf.Actions[1].(*ast.Execute).Parameters)
		assert.Equal(t, ast.Assignments{
			"foo": "bar",
			"bar": "bazz",
		}, gf.Actions[1].(*ast.Execute).Returns)

		assert.Equal(t, "../pathTo/someGoatfile", gf.Actions[2].(*ast.Execute).Path)
		assert.Equal(t, ast.KV{
			"foo":  "hello",
			"bar":  int64(2),
			"bazz": ParameterValue(".someParam"),
		}, gf.Actions[2].(*ast.Execute).Parameters)
		assert.Equal(t, ast.Assignments{
			"foo": "bar",
			"bar": "bazz",
		}, gf.Actions[2].(*ast.Execute).Returns)
	})
}

// --- Helpers --------------------------------------------

func stringParser(raw string) *Parser {
	return NewParser(strings.NewReader(raw), ".")
}

func swapTicks(v string) string {
	return strings.ReplaceAll(v, "´", "`")
}

func importsToPaths(imp []ast.Import) []string {
	res := make([]string, 0, len(imp))
	for _, i := range imp {
		res = append(res, i.Path)
	}
	return res
}
