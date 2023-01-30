package embedded

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/studio-b12/goat/pkg/goatfile"
)

func TestNewGoatfile(t *testing.T) {
	_, err := goatfile.Unmarshal(string(NewGoatfile), "")
	assert.Nil(t, err)
}
