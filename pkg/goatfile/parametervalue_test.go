package goatfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParameterValue_ApplyTemplate_Primitive(t *testing.T) {
	t.Run("result-integer", func(t *testing.T) {
		const raw = ParameterValue(`.param`)

		res, err := raw.ApplyTemplate(map[string]any{
			"param": 123,
		})
		assert.Nil(t, err)
		assert.Equal(t, int64(123), res)
	})

	t.Run("result-float", func(t *testing.T) {
		const raw = ParameterValue(`.param`)

		res, err := raw.ApplyTemplate(map[string]any{
			"param": 0.123,
		})
		assert.Nil(t, err)
		assert.Equal(t, 0.123, res)
	})

	t.Run("result-bool", func(t *testing.T) {
		const raw = ParameterValue(`.param`)

		res, err := raw.ApplyTemplate(map[string]any{
			"param": true,
		})
		assert.Nil(t, err)
		assert.Equal(t, true, res)
	})

	t.Run("result-string", func(t *testing.T) {
		const raw = ParameterValue(`.param`)

		res, err := raw.ApplyTemplate(map[string]any{
			"param": `"some string"`,
		})
		assert.Nil(t, err)
		assert.Equal(t, "some string", res)
	})
}

func TestParameterValue_ApplyTemplate_Complex(t *testing.T) {
	t.Run("print-1", func(t *testing.T) {
		const raw = ParameterValue(`print "123"`)

		res, err := raw.ApplyTemplate(nil)
		assert.Nil(t, err)
		assert.Equal(t, int64(123), res)
	})

	t.Run("print-2", func(t *testing.T) {
		const raw = ParameterValue(` print .param1 .param2 `)

		res, err := raw.ApplyTemplate(map[string]any{
			"param1": "123",
			"param2": ".456",
		})
		assert.Nil(t, err)
		assert.Equal(t, 123.456, res)
	})
}
