package embedded

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/studio-b12/gurl/pkg/gurlfile"
)

func TestNewGurlfile(t *testing.T) {
	_, err := gurlfile.Unmarshal(string(NewGurlfile), "")
	assert.Nil(t, err)
}
