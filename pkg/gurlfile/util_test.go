package gurlfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveComments(t *testing.T) {
	var raw, exp, res string

	raw = `
some text // some comments
// another comment
text//close to comment
https://example.com

test
	`

	exp = `
some text

text//close to comment
https://example.com

test
	`

	res = removeComments(raw)
	assert.Equal(t, exp, res)
}
