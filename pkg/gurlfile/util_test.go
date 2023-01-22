package gurlfile

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

	t.Run("template-builtin-base64", func(t *testing.T) {
		const raw = `{{base64 "hello world"}}`

		res, err := applyTemplate(raw, nil)
		assert.Nil(t, err)
		assert.Equal(t, "aGVsbG8gd29ybGQ", res)
	})

	t.Run("template-builtin-base64url", func(t *testing.T) {
		const raw = `{{base64Url "hello world"}}`

		res, err := applyTemplate(raw, nil)
		assert.Nil(t, err)
		assert.Equal(t, "aGVsbG8gd29ybGQ", res)
	})
}
