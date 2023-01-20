package gurlfile

import (
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

		assert.Nil(t, err)
		assert.Equal(t, 1, len(res.Tests))
		assert.Equal(t, "GET", res.Tests[0].Method)
		assert.Equal(t, "https://example.com", res.Tests[0].URI)
	})
}

func stringParser(raw string) *Parser {
	return NewParser(strings.NewReader(raw))
}
