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
		assert.Equal(t, ast.RequestHeader{ast.HeaderEntries{ast.KVList[string]{
			ast.KV[string]{Key: "Key-1", Value: "value 1", Pos: pos(38, 5, 0)},
			ast.KV[string]{Key: "key-2", Value: "value 2", Pos: pos(53, 6, 0)},
		}}}, res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestHeader))
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
		assert.Equal(t, ast.RequestHeader{ast.HeaderEntries{ast.KVList[string]{
			ast.KV[string]{Key: "Key-1", Value: "value 1", Pos: pos(38, 5, 0)},
			ast.KV[string]{Key: "key-2", Value: "value 2", Pos: pos(53, 6, 0)},
		}}}, res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestHeader))
		assert.Equal(t, ast.TextBlock{Content: "some\nbody\n"}, res.Actions[0].(*ast.Request).Blocks[1].(ast.RequestBody).DataContent)
		assert.Equal(t, ast.RequestQueryParams{ast.KVList[any]{
			ast.KV[any]{Key: "keyInt", Value: int64(2), Pos: pos(101, 13, 0)},
			ast.KV[any]{Key: "keyString", Value: "some string", Pos: pos(112, 14, 0)},
		}}, res.Actions[0].(*ast.Request).Blocks[2].(ast.RequestQueryParams))
		assert.Equal(t, ast.RequestAuth{ast.KVList[any]{
			ast.KV[any]{Key: "username", Value: "foo", Pos: pos(146, 17, 0)},
			ast.KV[any]{Key: "password", Value: "{{.creds.password}}", Pos: pos(163, 18, 0)},
		}}, res.Actions[0].(*ast.Request).Blocks[3].(ast.RequestAuth))
	})

	t.Run("single-file-value", func(t *testing.T) {
		const raw = `
		
GET https://example.com

[FormData]
file = @file.txt
fileWithContentType = @../../file/with/content/type.txt:text/csv
spaced = @"my file.png":"image / png"
windows = @"C:\test\file.png":image/png // yes, people actually still use this
		`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t, 1, len(res.Actions))
		assert.Equal(t, "GET", res.Actions[0].(*ast.Request).Head.Method)
		assert.Equal(t, "https://example.com", res.Actions[0].(*ast.Request).Head.Url)

		formData := res.Actions[0].(*ast.Request).Blocks[0].(ast.FormData).KVList.ToMap()
		assert.Equal(t, ast.FileDescriptor{Path: "file.txt"}, formData["file"])
		assert.Equal(t, ast.FileDescriptor{
			Path:        "../../file/with/content/type.txt",
			ContentType: "text/csv",
		}, formData["fileWithContentType"])
		assert.Equal(t, ast.FileDescriptor{
			Path:        "my file.png",
			ContentType: "image / png",
		}, formData["spaced"])
		assert.Equal(t, ast.FileDescriptor{
			Path:        "C:\\test\\file.png",
			ContentType: "image/png",
		}, formData["windows"])
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
		assert.Equal(t, ast.RequestHeader{ast.HeaderEntries{}}, res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestHeader))
	})

	t.Run("values", func(t *testing.T) {
		const raw = `

GET https://example.com

[Header]
key: value
// Hey, A comment!
key-2:  value 2
Some-Key-3: 		some value 3
SOME_KEY_4: 		§$%&/()=!§

multiple-1: value 1
multiple-1: value 2

		`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t, ast.RequestHeader{ast.HeaderEntries{ast.KVList[string]{
			ast.KV[string]{Key: "key", Value: "value", Pos: pos(36, 5, 0)},
			ast.KV[string]{Key: "key-2", Value: "value 2", Pos: pos(66, 7, 0)},
			ast.KV[string]{Key: "Some-Key-3", Value: "some value 3", Pos: pos(82, 8, 0)},
			ast.KV[string]{Key: "SOME_KEY_4", Value: "§$%&/()=!§", Pos: pos(109, 9, 0)},
			ast.KV[string]{Key: "multiple-1", Value: "value 1", Pos: pos(135, 11, 0)},
			ast.KV[string]{Key: "multiple-1", Value: "value 2", Pos: pos(155, 12, 0)},
		}}}, res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestHeader))
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
		assert.Equal(t, ast.NoContent{}, res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestBody).DataContent)
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
			ast.TextBlock{Content: "some body content\nsome more content\n"},
			res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestBody).DataContent)
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
			ast.TextBlock{Content: "some body content\nsome more content\n"},
			(res.Actions[0].(*ast.Request).Blocks[0]).(ast.RequestBody).DataContent)
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
			ast.TextBlock{Content: "some body content\nsome more content\n"},
			res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestBody).DataContent)
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
			ast.TextBlock{Content: "some body content\nsome more content"},
			res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestBody).DataContent)
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
			ast.TextBlock{Content: "some body content\nsome more content"},
			res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestBody).DataContent)
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
			ast.TextBlock{Content: "some body content\nsome more content"},
			res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestBody).DataContent)
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
		assert.Equal(t, ast.NoContent{}, res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestBody).DataContent)
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
			ast.TextBlock{Content: "some body content\nsome more content\n"},
			res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestBody).DataContent)
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
			ast.TextBlock{Content: "some body content\n\n[QueryParams]\nsome more content\n"},
			res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestBody).DataContent)
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
			ast.TextBlock{Content: "some body content\n\n---\n\nsome more content\n"},
			res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestBody).DataContent)
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
			ast.TextBlock{Content: "some body content\n\n### setup\n\nsome more content\n"},
			res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestBody).DataContent)
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
			ast.TextBlock{Content: `assert(response.StatusCode == 200, "invalid status code");` +
				"\nvar id = response.Body.id;\n"},
			res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestScript).DataContent)
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
		assert.Equal(t, ast.RequestQueryParams{ast.KVList[any]{}}, res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestQueryParams))
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
		assert.Equal(t, ast.RequestQueryParams{ast.KVList[any]{
			ast.KV[any]{Key: "string1", Value: "some string 1", Pos: pos(41, 5, 0)},
			ast.KV[any]{Key: "string2", Value: "some string 2", Pos: pos(67, 6, 0)},
			ast.KV[any]{Key: "string3", Value: "some string 3", Pos: pos(97, 7, 0)},
		}}, res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestQueryParams))
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
		assert.Equal(t, ast.RequestQueryParams{ast.KVList[any]{
			ast.KV[any]{Key: "int1", Value: int64(1), Pos: pos(41, 5, 0)},
			ast.KV[any]{Key: "int2", Value: int64(1000), Pos: pos(50, 6, 0)},
			ast.KV[any]{Key: "int3", Value: int64(-123), Pos: pos(63, 7, 0)},
		}}, res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestQueryParams))
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
		assert.Equal(t, ast.RequestQueryParams{ast.KVList[any]{
			ast.KV[any]{Key: "float1", Value: float64(1.234), Pos: pos(41, 5, 0)},
			ast.KV[any]{Key: "float2", Value: float64(1000.234), Pos: pos(56, 6, 0)},
			ast.KV[any]{Key: "float3", Value: float64(0.12), Pos: pos(75, 7, 0)},
			ast.KV[any]{Key: "float4", Value: float64(-12.34), Pos: pos(89, 8, 0)},
		}}, res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestQueryParams))
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
		assert.Equal(t, ast.RequestQueryParams{ast.KVList[any]{
			ast.KV[any]{Key: "bool1", Value: true, Pos: pos(41, 5, 0)},
			ast.KV[any]{Key: "bool2", Value: false, Pos: pos(54, 6, 0)},
		}}, res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestQueryParams))
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
		assert.Equal(t, ast.RequestQueryParams{ast.KVList[any]{
			ast.KV[any]{Key: "arrayEmpty1", Value: []any(nil), Pos: pos(41, 5, 0)},
			ast.KV[any]{Key: "arrayEmpty2", Value: []any(nil), Pos: pos(58, 6, 0)},
			ast.KV[any]{Key: "arrayString1", Value: []any{"some string"}, Pos: pos(79, 8, 0)},
			ast.KV[any]{Key: "arrayString2", Value: []any{"some string", "another string", "and another one"}, Pos: pos(110, 9, 0)},
			ast.KV[any]{Key: "arrayInt1", Value: []any{int64(1)}, Pos: pos(178, 11, 0)},
			ast.KV[any]{Key: "arrayInt2", Value: []any{int64(1), int64(2), int64(-3), int64(4_000)}, Pos: pos(194, 12, 0)},
			ast.KV[any]{Key: "arrayFloat1", Value: []any{1.23}, Pos: pos(224, 14, 0)},
			ast.KV[any]{Key: "arrayFloat2", Value: []any{1.0, -1.1, 1.234}, Pos: pos(245, 15, 0)},
			ast.KV[any]{Key: "arrayMixed", Value: []any{"a string", int64(2), 3.456, true}, Pos: pos(278, 17, 0)},
			ast.KV[any]{Key: "arrayNested", Value: []any{[]any{int64(1), int64(2)}, []any{[]any{true, false}, "foo"}}, Pos: pos(321, 19, 0)},
			ast.KV[any]{Key: "arrayMultiline", Value: []any{"foo", "bar"}, Pos: pos(368, 21, 0)},
			ast.KV[any]{Key: "arrayLeadingComma1", Value: []any{true, false}, Pos: pos(405, 26, 0)},
			ast.KV[any]{Key: "arrayLeadingComma2", Value: []any{true, false}, Pos: pos(441, 27, 0)},
		}}, res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestQueryParams))
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

		assert.Equal(t, 5, len(res.Comments))

		assert.Equal(t, "Some comment", res.Comments[0].Content)
		assert.Equal(t, 1, res.Comments[0].Pos.Line)
		assert.Equal(t, 0, res.Comments[0].Pos.LinePos)

		assert.Equal(t, "Some comment", res.Comments[1].Content)
		assert.Equal(t, 2, res.Comments[1].Pos.Line)
		assert.Equal(t, 3, res.Comments[1].Pos.LinePos)

		assert.Equal(t, "another comment", res.Comments[2].Content)
		assert.Equal(t, 3, res.Comments[2].Pos.Line)
		assert.Equal(t, 24, res.Comments[2].Pos.LinePos)

		assert.Equal(t, "comment", res.Comments[3].Content)
		assert.Equal(t, 4, res.Comments[3].Pos.Line)
		assert.Equal(t, 0, res.Comments[3].Pos.LinePos)

		assert.Equal(t, "heyo", res.Comments[4].Content)
		assert.Equal(t, 6, res.Comments[4].Pos.Line)
		assert.Equal(t, 0, res.Comments[4].Pos.LinePos)
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
		assert.Equal(t, ast.RequestQueryParams{ast.KVList[any]{
			ast.KV[any]{Key: "key1", Value: "value", Pos: pos(79, 5, 0)},
			ast.KV[any]{Key: "key2", Value: 1.23, Pos: pos(113, 6, 0)},
			ast.KV[any]{Key: "arr", Value: []any{int64(1), int64(2)}, Pos: pos(151, 8, 0)},
		}}, res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestQueryParams))

		assert.Equal(t, 10, len(res.Comments))
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

func TestParse_Delimiters(t *testing.T) {
	t.Run("general", func(t *testing.T) {
		const raw = `
GET https://foo.com

---

GET https://bar.com

[Header]
A: B

[Body]
some stuff

[Script]
hello world

----

-----

GET https://baz.com

### Setup

------

GET https://bar.com

-------
		`

		p := stringParser(raw)
		res, err := p.Parse()

		assert.Nil(t, err, err)
		assert.Equal(t, 5, len(res.Delimiters))

		assert.Equal(t, 0, res.Delimiters[0].ExtraLen)
		assert.Equal(t, 1, res.Delimiters[1].ExtraLen)
		assert.Equal(t, 2, res.Delimiters[2].ExtraLen)
		assert.Equal(t, 3, res.Delimiters[3].ExtraLen)
		assert.Equal(t, 4, res.Delimiters[4].ExtraLen)
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
		assert.Equal(t, ParameterValue(".param"), res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestOptions).GetUnchecked("someoption"))
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
		assert.Equal(t, ParameterValue(" .param "), res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestOptions).GetUnchecked("someoption"))
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
		assert.Equal(t, ParameterValue(" print {{.param1}} {{.param2}} "), res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestOptions).GetUnchecked("someoption1"))
		assert.Equal(t, ParameterValue(" print {{if .param1}}true{{else}}false{{end}} "), res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestOptions).GetUnchecked("someoption2"))
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
		assert.Equal(t, ParameterValue(` print "}}" `), res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestOptions).GetUnchecked("someoption"))
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
		assert.Equal(t, ParameterValue(" print `}}` "), res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestOptions).GetUnchecked("someoption"))
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
		assert.Equal(t, ParameterValue(` print {{ "}}" }} `), res.Actions[0].(*ast.Request).Blocks[0].(ast.RequestOptions).GetUnchecked("someoption"))
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
		assert.Equal(t, ast.TextBlock{Content: "script stuff 1\n"}, gf.Actions[1].(*ast.Request).Blocks[0].(ast.RequestScript).DataContent)
		assert.Equal(t, "Log section 2", gf.Actions[2].(ast.LogSection).Content)
		assert.Equal(t, ast.TextBlock{Content: "script stuff 2\n"}, gf.Actions[3].(*ast.Request).Blocks[0].(ast.RequestScript).DataContent)
		assert.Equal(t, "Log section 3", gf.Actions[4].(ast.LogSection).Content)
		assert.Equal(t, ast.TextBlock{Content: "script stuff 3"}, gf.Actions[5].(*ast.Request).Blocks[0].(ast.RequestScript).DataContent)
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
		assert.Equal(t, "bar", gf.Sections[0].(ast.SectionDefaults).Request.Blocks[0].(ast.RequestHeader).HeaderEntries.GetUnchecked("foo"))
		assert.Equal(t, "world", gf.Sections[0].(ast.SectionDefaults).Request.Blocks[0].(ast.RequestHeader).HeaderEntries.GetUnchecked("hello"))
		assert.Equal(t, int64(42), gf.Sections[0].(ast.SectionDefaults).Request.Blocks[1].(ast.RequestOptions).GetUnchecked("answer"))
		assert.Equal(t, "value", gf.Sections[0].(ast.SectionDefaults).Request.Blocks[1].(ast.RequestOptions).GetUnchecked("some"))
		assert.Equal(t, "foo", gf.Sections[0].(ast.SectionDefaults).Request.Blocks[2].(ast.RequestAuth).GetUnchecked("username"))
		assert.Equal(t, "bar", gf.Sections[0].(ast.SectionDefaults).Request.Blocks[2].(ast.RequestAuth).GetUnchecked("password"))
		assert.Equal(t, ast.TextBlock{Content: "hello\nworld\n"}, gf.Sections[0].(ast.SectionDefaults).Request.Blocks[3].(ast.RequestBody).DataContent)
		assert.Equal(t, ast.TextBlock{Content: "some pre script\n"}, gf.Sections[0].(ast.SectionDefaults).Request.Blocks[4].(ast.RequestPreScript).DataContent)
		assert.Equal(t, ast.TextBlock{Content: "some script\n"}, gf.Sections[0].(ast.SectionDefaults).Request.Blocks[5].(ast.RequestScript).DataContent)
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
		assert.Equal(t, "bar", gf.Sections[0].(ast.SectionDefaults).Request.Blocks[0].(ast.RequestHeader).GetUnchecked("foo"))
		assert.Equal(t, "world", gf.Sections[0].(ast.SectionDefaults).Request.Blocks[0].(ast.RequestHeader).GetUnchecked("hello"))
		assert.Equal(t, ast.TextBlock{Content: "some script\n"}, gf.Sections[0].(ast.SectionDefaults).Request.Blocks[1].(ast.RequestScript).DataContent)
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
		assert.Equal(t, "bar", gf.Sections[0].(ast.SectionDefaults).Request.Blocks[0].(ast.RequestHeader).HeaderEntries.GetUnchecked("foo"))
		assert.Equal(t, "world", gf.Sections[0].(ast.SectionDefaults).Request.Blocks[0].(ast.RequestHeader).HeaderEntries.GetUnchecked("hello"))
		assert.Equal(t, ast.TextBlock{Content: "some script\n"}, gf.Sections[0].(ast.SectionDefaults).Request.Blocks[1].(ast.RequestScript).DataContent)
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
		assert.Equal(t, "bar", gf.Sections[0].(ast.SectionDefaults).Request.Blocks[0].(ast.RequestHeader).HeaderEntries.GetUnchecked("foo"))
		assert.Equal(t, "world", gf.Sections[0].(ast.SectionDefaults).Request.Blocks[0].(ast.RequestHeader).HeaderEntries.GetUnchecked("hello"))
		assert.Equal(t, ast.TextBlock{Content: "some script\n"}, gf.Sections[0].(ast.SectionDefaults).Request.Blocks[1].(ast.RequestScript).DataContent)
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
		assert.Equal(t, "bar", gf.Sections[1].(ast.SectionDefaults).Request.Blocks[0].(ast.RequestHeader).HeaderEntries.GetUnchecked("foo"))
		assert.Equal(t, "world", gf.Sections[1].(ast.SectionDefaults).Request.Blocks[0].(ast.RequestHeader).HeaderEntries.GetUnchecked("hello"))
		assert.Equal(t, ast.TextBlock{Content: "some script\n"}, gf.Sections[1].(ast.SectionDefaults).Request.Blocks[1].(ast.RequestScript).DataContent)
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
		assert.Equal(t, "bar", gf.Sections[0].(ast.SectionDefaults).Request.Blocks[0].(ast.RequestHeader).HeaderEntries.GetUnchecked("foo"))
		assert.Equal(t, "moon", gf.Sections[1].(ast.SectionDefaults).Request.Blocks[0].(ast.RequestHeader).HeaderEntries.GetUnchecked("hello"))
		assert.Equal(t, ast.TextBlock{Content: "some script\n"}, gf.Sections[0].(ast.SectionDefaults).Request.Blocks[1].(ast.RequestScript).DataContent)
		assert.Equal(t, ast.TextBlock{Content: "some script 2\n"}, gf.Sections[1].(ast.SectionDefaults).Request.Blocks[1].(ast.RequestScript).DataContent)
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
		assert.Equal(t, ast.KVList[any]{
			ast.KV[any]{Key: "foo", Value: int64(1), Pos: pos(33, 1, 32)},
		}, gf.Actions[0].(*ast.Execute).Parameters)

		assert.Equal(t, "../pathTo/someGoatfile", gf.Actions[1].(*ast.Execute).Path)
		assert.Equal(t, ast.KVList[any]{
			ast.KV[any]{Key: "foo", Value: int64(1), Pos: pos(77, 5, 32)},
			ast.KV[any]{Key: "bar", Value: "hello", Pos: pos(82, 5, 37)},
		}, gf.Actions[1].(*ast.Execute).Parameters)

		assert.Equal(t, "../pathTo/someGoatfile", gf.Actions[2].(*ast.Execute).Path)
		assert.Equal(t, ast.KVList[any]{
			ast.KV[any]{Key: "foo", Value: "hello", Pos: pos(134, 10, 0)},
			ast.KV[any]{Key: "bar", Value: int64(2), Pos: pos(152, 11, 0)},
			ast.KV[any]{Key: "bazz", Value: ParameterValue(".someParam"), Pos: pos(159, 12, 0)},
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
		assert.Equal(t, ast.KVList[any]{
			ast.KV[any]{Key: "foo", Value: int64(1), Pos: pos(33, 1, 32)},
		}, gf.Actions[0].(*ast.Execute).Parameters)
		assert.Equal(t, ast.Assignments{ast.KVList[string]{
			ast.KV[string]{Key: "foo", Value: "bar", Pos: pos(48, 1, 47)},
		}}, gf.Actions[0].(*ast.Execute).Returns)

		assert.Equal(t, "../pathTo/someGoatfile", gf.Actions[1].(*ast.Execute).Path)
		assert.Equal(t, ast.KVList[any]{
			ast.KV[any]{Key: "foo", Value: int64(1), Pos: pos(97, 5, 32)},
			ast.KV[any]{Key: "bar", Value: "hello", Pos: pos(102, 5, 37)},
		}, gf.Actions[1].(*ast.Execute).Parameters)
		assert.Equal(t, ast.Assignments{ast.KVList[string]{
			ast.KV[string]{Key: "foo", Value: "bar", Pos: pos(124, 5, 59)},
			ast.KV[string]{Key: "bar", Value: "bazz", Pos: pos(134, 5, 69)},
		}}, gf.Actions[1].(*ast.Execute).Returns)

		assert.Equal(t, "../pathTo/someGoatfile", gf.Actions[2].(*ast.Execute).Path)
		assert.Equal(t, ast.KVList[any]{
			ast.KV[any]{Key: "foo", Value: "hello", Pos: pos(186, 10, 0)},
			ast.KV[any]{Key: "bar", Value: int64(2), Pos: pos(204, 11, 0)},
			ast.KV[any]{Key: "bazz", Value: ParameterValue(".someParam"), Pos: pos(211, 12, 0)},
		}, gf.Actions[2].(*ast.Execute).Parameters)
		assert.Equal(t, ast.Assignments{ast.KVList[string]{
			ast.KV[string]{Key: "foo", Value: "bar", Pos: pos(244, 14, 0)},
			ast.KV[string]{Key: "bar", Value: "bazz", Pos: pos(258, 15, 0)},
		}}, gf.Actions[2].(*ast.Execute).Returns)
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

func pos(pos int, line int, linePos int) ast.Pos {
	return ast.Pos{Pos: pos, Line: line, LinePos: linePos}
}
