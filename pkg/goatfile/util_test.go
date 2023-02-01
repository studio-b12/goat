package goatfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyTemplate(t *testing.T) {
	t.Run("no-templates", func(t *testing.T) {
		const raw = `This string has no templates at all!`

		res, err := applyTemplate(raw, nil)
		assert.Nil(t, err)
		assert.Equal(t, raw, res)
	})

	t.Run("template-apply", func(t *testing.T) {
		const raw = `Hello, I am {{.name}}!`

		res, err := applyTemplate(raw, map[string]any{"name": "Jhon"})
		assert.Nil(t, err)
		assert.Equal(t, "Hello, I am Jhon!", res)
	})

	t.Run("template-apply-withescape", func(t *testing.T) {
		const raw = `Hello, I am {{.name}} and I am \{\{.age\}\} years old!`

		res, err := applyTemplate(raw, map[string]any{"name": "Jhon"})
		assert.Nil(t, err)
		assert.Equal(t, "Hello, I am Jhon and I am {{.age}} years old!", res)
	})

	t.Run("template-apply-escessiveparams", func(t *testing.T) {
		const raw = `Hello, I am {{.name}}!`

		res, err := applyTemplate(raw, map[string]any{
			"name":   "Jhon",
			"age":    123,
			"height": 1.75,
		})
		assert.Nil(t, err)
		assert.Equal(t, "Hello, I am Jhon!", res)
	})

	t.Run("template-apply-missingparams", func(t *testing.T) {
		const raw = `Hello, I am {{.name}} and I am {{.age}} years old!`

		_, err := applyTemplate(raw, map[string]any{
			"name": "Jhon",
		})
		assert.Error(t, err)
	})
}

func TestApplyTemplateToArray(t *testing.T) {
	t.Run("onedimensional-strings", func(t *testing.T) {
		arr := []any{"{{.foo}}", "- {{ .bar }} -", "bazz"}
		params := map[string]any{
			"foo": "some foo",
			"bar": "some bar",
		}

		err := applyTemplateToArray(arr, params)
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

		err := applyTemplateToArray(arr, params)
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

		err := applyTemplateToArray(arr, params)
		assert.Error(t, err)
	})

	t.Run("multidimensional-mixed", func(t *testing.T) {
		arr := []any{"{{.foo}}", []any{"bar", 123, "{{.bazz}}", []any{"{{.fuzz}}", true}}}
		params := map[string]any{
			"foo":  "some foo",
			"bazz": "some bazz",
			"fuzz": "some fuzz",
		}

		err := applyTemplateToArray(arr, params)
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

		err := applyTemplateToArray(arr, params)
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

		err := applyTemplateToMap(m, params)
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

		err := applyTemplateToMap(m, params)
		assert.Error(t, err)
	})
}

func TestApplyTemplate_Builtins(t *testing.T) {
	t.Run("template-builtin-base64", func(t *testing.T) {
		const raw = `{{base64 "hello world"}}`

		res, err := applyTemplate(raw, nil)
		assert.Nil(t, err)
		assert.Equal(t, "aGVsbG8gd29ybGQ", res)
	})

	t.Run("template-builtin-base64url", func(t *testing.T) {
		const raw = `{{base64url "hello world"}}`

		res, err := applyTemplate(raw, nil)
		assert.Nil(t, err)
		assert.Equal(t, "aGVsbG8gd29ybGQ", res)
	})

	t.Run("template-builtin-md5", func(t *testing.T) {
		const raw = `{{md5 "hello world"}}`

		res, err := applyTemplate(raw, nil)
		assert.Nil(t, err)
		assert.Equal(t, "5eb63bbbe01eeed093cb22bb8f5acdc3", res)
	})

	t.Run("template-builtin-sha1", func(t *testing.T) {
		const raw = `{{sha1 "hello world"}}`

		res, err := applyTemplate(raw, nil)
		assert.Nil(t, err)
		assert.Equal(t, "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed", res)
	})

	t.Run("template-builtin-sha256", func(t *testing.T) {
		const raw = `{{sha256 "hello world"}}`

		res, err := applyTemplate(raw, nil)
		assert.Nil(t, err)
		assert.Equal(t, "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9", res)
	})

	t.Run("template-builtin-sha512", func(t *testing.T) {
		const raw = `{{sha512 "hello world"}}`

		res, err := applyTemplate(raw, nil)
		assert.Nil(t, err)
		assert.Equal(t, "309ecc489c12d6eb4cc40f50c902f2b4d0ed77ee511a7c7a9bcd3ca86d4cd86f989dd35bc5ff499670da34255b45b0cfd830e81f605dcf7dc5542e93ae9cd76f", res)
	})
}

func TestExtend(t *testing.T) {
	assert.Equal(t, "hello.txt", extend("hello", "txt"))
	assert.Equal(t, "hello.png", extend("hello.png", "txt"))
	assert.Equal(t, "a/b.c/hello.txt", extend("a/b.c/hello", "txt"))
	assert.Equal(t, "a/b.c/hello.png", extend("a/b.c/hello.png", "txt"))
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