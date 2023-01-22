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
