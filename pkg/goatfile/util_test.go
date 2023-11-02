package goatfile

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestApplyTemplate(t *testing.T) {
	t.Run("no-templates", func(t *testing.T) {
		const raw = `This string has no templates at all!`

		res, err := ApplyTemplate(raw, nil)
		assert.Nil(t, err)
		assert.Equal(t, raw, res)
	})

	t.Run("template-apply", func(t *testing.T) {
		const raw = `Hello, I am {{.name}}!`

		res, err := ApplyTemplate(raw, map[string]any{"name": "Jhon"})
		assert.Nil(t, err)
		assert.Equal(t, "Hello, I am Jhon!", res)
	})

	t.Run("template-apply-withescape", func(t *testing.T) {
		const raw = `Hello, I am {{.name}} and I am \{\{.age\}\} years old!`

		res, err := ApplyTemplate(raw, map[string]any{"name": "Jhon"})
		assert.Nil(t, err)
		assert.Equal(t, "Hello, I am Jhon and I am {{.age}} years old!", res)
	})

	t.Run("template-apply-escessiveparams", func(t *testing.T) {
		const raw = `Hello, I am {{.name}}!`

		res, err := ApplyTemplate(raw, map[string]any{
			"name":   "Jhon",
			"age":    123,
			"height": 1.75,
		})
		assert.Nil(t, err)
		assert.Equal(t, "Hello, I am Jhon!", res)
	})

	t.Run("template-apply-missingparams", func(t *testing.T) {
		const raw = `Hello, I am {{.name}} and I am {{.age}} years old!`

		_, err := ApplyTemplate(raw, map[string]any{
			"name": "Jhon",
		})
		assert.Error(t, err)
	})

	t.Run("template-list", func(t *testing.T) {
		const raw = `{{ .param }}`

		res, err := ApplyTemplate(raw, map[string]any{
			"param": []any{"foo", "bar", "bazz", 1, 2, 3.1415},
		})
		assert.Nil(t, err)
		assert.Equal(t, `["foo","bar","bazz",1,2,3.1415]`, res)
	})

	t.Run("template-map", func(t *testing.T) {
		const raw = `{{ .param }}`

		res, err := ApplyTemplate(raw, map[string]any{
			"param": map[string]any{
				"bar": 123,
				"foo": "bar",
			},
		})
		assert.Nil(t, err)
		assert.Equal(t, `{"bar":123,"foo":"bar"}`, res)
	})
}

func TestApplyTemplateToArray(t *testing.T) {
	t.Run("onedimensional-strings", func(t *testing.T) {
		arr := []any{"{{.foo}}", "- {{ .bar }} -", "bazz"}
		params := map[string]any{
			"foo": "some foo",
			"bar": "some bar",
		}

		err := ApplyTemplateToArray(arr, params)
		assert.Nil(t, err)
		assert.Equal(t, []any{
			"some foo",
			"- some bar -",
			"bazz",
		}, arr)
	})

	t.Run("onedimensional-mixed", func(t *testing.T) {
		arr := []any{"{{.foo}}", 123, true}
		params := map[string]any{
			"foo": "some foo",
			"bar": "some bar",
		}

		err := ApplyTemplateToArray(arr, params)
		assert.Nil(t, err)
		assert.Equal(t, []any{
			"some foo",
			123,
			true,
		}, arr)
	})

	t.Run("onedimensional-error", func(t *testing.T) {
		arr := []any{"{{.bazz}}"}
		params := map[string]any{
			"foo": "some foo",
			"bar": "some bar",
		}

		err := ApplyTemplateToArray(arr, params)
		assert.Error(t, err)
	})

	t.Run("multidimensional-mixed", func(t *testing.T) {
		arr := []any{"{{.foo}}", []any{"bar", 123, "{{.bazz}}", []any{"{{.fuzz}}", true}}}
		params := map[string]any{
			"foo":  "some foo",
			"bazz": "some bazz",
			"fuzz": "some fuzz",
		}

		err := ApplyTemplateToArray(arr, params)
		assert.Nil(t, err)
		assert.Equal(t, []any{
			"some foo",
			[]any{
				"bar",
				123,
				"some bazz",
				[]any{
					"some fuzz",
					true,
				},
			},
		}, arr)
	})

	t.Run("multidimensional-error", func(t *testing.T) {
		arr := []any{"{{.foo}}", []any{"bar", 123, "{{.bazz}}", []any{"{{.fuzz}}", true}}}
		params := map[string]any{
			"foo":  "some foo",
			"bazz": "some bazz",
		}

		err := ApplyTemplateToArray(arr, params)
		assert.Error(t, err)
	})
}

func TestApplyTemplateToMap(t *testing.T) {
	t.Run("general", func(t *testing.T) {
		m := map[string]any{
			"a": ParameterValue(".foo"),
			"b": "{{.bar}}",
			"c": []any{"{{.bazz}}", "test", 123},
			"d": map[string]any{
				"a": "{{.fuzz}}",
				"b": ParameterValue(".foo"),
				"c": []any{"{{.bazz}}", "test", 123},
			},
			"e": 123,
			"f": true,
			"g": nil,
		}
		params := map[string]any{
			"foo":  123,
			"bar":  "some bar",
			"bazz": "some bazz",
			"fuzz": "some fuzz",
		}

		err := ApplyTemplateToMap(m, params)
		assert.Nil(t, err)
		assert.Equal(t, map[string]any{
			"a": int64(123),
			"b": "some bar",
			"c": []any{"some bazz", "test", 123},
			"d": map[string]any{
				"a": "some fuzz",
				"b": int64(123),
				"c": []any{"some bazz", "test", 123},
			},
			"e": 123,
			"f": true,
			"g": nil,
		}, m)
	})

	t.Run("error", func(t *testing.T) {
		m := map[string]any{
			"a": ParameterValue(".foo"),
			"b": "{{.bar}}",
			"c": []any{"{{.bazz}}", "test", 123},
			"d": map[string]any{
				"a": "{{.fuzz}}",
				"b": ParameterValue(".foo"),
				"c": []any{"{{.bazz}}", "test", 123},
			},
			"e": 123,
			"f": true,
			"g": nil,
		}
		params := map[string]any{
			"fuzz": "some fuzz",
		}

		err := ApplyTemplateToMap(m, params)
		assert.Error(t, err)
	})
}

func TestApplyTemplate_Builtins(t *testing.T) {
	t.Run("template-builtin-base64", func(t *testing.T) {
		const raw = `{{base64 "hello world"}}`

		res, err := ApplyTemplate(raw, nil)
		assert.Nil(t, err)
		assert.Equal(t, "aGVsbG8gd29ybGQ", res)
	})

	t.Run("template-builtin-base64url", func(t *testing.T) {
		const raw = `{{base64url "hello world"}}`

		res, err := ApplyTemplate(raw, nil)
		assert.Nil(t, err)
		assert.Equal(t, "aGVsbG8gd29ybGQ", res)
	})

	t.Run("template-builtin-md5", func(t *testing.T) {
		const raw = `{{md5 "hello world"}}`

		res, err := ApplyTemplate(raw, nil)
		assert.Nil(t, err)
		assert.Equal(t, "5eb63bbbe01eeed093cb22bb8f5acdc3", res)
	})

	t.Run("template-builtin-sha1", func(t *testing.T) {
		const raw = `{{sha1 "hello world"}}`

		res, err := ApplyTemplate(raw, nil)
		assert.Nil(t, err)
		assert.Equal(t, "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed", res)
	})

	t.Run("template-builtin-sha256", func(t *testing.T) {
		const raw = `{{sha256 "hello world"}}`

		res, err := ApplyTemplate(raw, nil)
		assert.Nil(t, err)
		assert.Equal(t, "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9", res)
	})

	t.Run("template-builtin-sha512", func(t *testing.T) {
		const raw = `{{sha512 "hello world"}}`

		res, err := ApplyTemplate(raw, nil)
		assert.Nil(t, err)
		assert.Equal(t, "309ecc489c12d6eb4cc40f50c902f2b4d0ed77ee511a7c7a9bcd3ca86d4cd86f989dd35bc5ff499670da34255b45b0cfd830e81f605dcf7dc5542e93ae9cd76f", res)
	})

	t.Run("template-builtin-randomString", func(t *testing.T) {
		const raw = `{{randomString}}`

		res, err := ApplyTemplate(raw, nil)
		assert.Nil(t, err)
		assert.Len(t, res, 8)
	})

	t.Run("template-builtin-randomString-parameterized", func(t *testing.T) {
		const raw = `{{randomString 20}}`

		res, err := ApplyTemplate(raw, nil)
		assert.Nil(t, err)
		assert.Len(t, res, 20)
	})

	t.Run("template-builtin-randomInt", func(t *testing.T) {
		const raw = `{{randomString}}`

		_, err := ApplyTemplate(raw, nil)
		assert.Nil(t, err)
	})

	t.Run("template-builtin-randomInt-parameterized", func(t *testing.T) {
		const raw = `{{randomString 20}}`

		_, err := ApplyTemplate(raw, nil)
		assert.Nil(t, err)
	})

	t.Run("template-builtin-timestamp", func(t *testing.T) {
		const raw = `{{timestamp}}`

		now := time.Now().Unix()
		res, err := ApplyTemplate(raw, nil)
		assert.Nil(t, err)

		resInt, err := strconv.Atoi(res)
		assert.Nil(t, err)
		assert.InDelta(t, now, resInt, 1)
	})

	t.Run("template-builtin-timestamp-parameterized", func(t *testing.T) {
		const raw = `{{timestamp "2006-01-02T15:04:05Z07:00"}}`

		now := time.Now().Unix()
		res, err := ApplyTemplate(raw, nil)
		assert.Nil(t, err)

		resTime, err := time.Parse(time.RFC3339, res)
		assert.Nil(t, err)

		assert.InDelta(t, now, resTime.Unix(), 1)
	})

	t.Run("template-builtin-isset", func(t *testing.T) {
		const raw = `{{isset . "somekey"}}`

		res, err := ApplyTemplate(raw, map[string]any{"anotherkey": "anotherval"})
		assert.Nil(t, err)
		assert.Equal(t, res, "false")

		res, err = ApplyTemplate(raw, map[string]any{"somekey": "someval"})
		assert.Nil(t, err)
		assert.Equal(t, res, "true")
	})

	t.Run("template-builtin-json", func(t *testing.T) {
		res, err := ApplyTemplate(`{{json .foo}}`, map[string]any{"foo": map[string]any{"bar": 1, "bazz": 2}})
		assert.Nil(t, err)
		assert.Equal(t, res, `{"bar":1,"bazz":2}`)

		res, err = ApplyTemplate(`{{json .foo "  "}}`, map[string]any{"foo": map[string]any{"bar": 1, "bazz": 2}})
		assert.Nil(t, err)
		assert.Equal(t, res, "{\n  \"bar\": 1,\n  \"bazz\": 2\n}")

		res, err = ApplyTemplate(`{{json 1}}`, nil)
		assert.Nil(t, err)
		assert.Equal(t, res, "1")
	})
}

func TestExtend(t *testing.T) {
	assert.Equal(t, "hello.txt", Extend("hello", "txt"))
	assert.Equal(t, "hello.png", Extend("hello.png", "txt"))
	assert.Equal(t, "a/b.c/hello.txt", Extend("a/b.c/hello", "txt"))
	assert.Equal(t, "a/b.c/hello.png", Extend("a/b.c/hello.png", "txt"))
}

func TestCrlf2lf(t *testing.T) {
	assert.Equal(t, "hello\nworld\n", crlf2lf("hello\r\nworld\r\n"))
}

func TestUnescapeTemplateDelims(t *testing.T) {
	assert.Equal(t, "hello {{.world}}", unescapeTemplateDelims("hello {{.world}}"))
	assert.Equal(t, "hello {{.world}}", unescapeTemplateDelims("hello \\{\\{.world\\}\\}"))
	assert.Equal(t, "{{}}", unescapeTemplateDelims("\\{{}\\}"))
	assert.Equal(t, "{{\\n}}", unescapeTemplateDelims("\\{{\\n}\\}"))
}
