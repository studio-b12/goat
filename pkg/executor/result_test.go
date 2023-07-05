package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	a := Result{}
	b := Result{
		Tests: ResultSection{
			all:    3,
			failed: 2,
		},
	}

	a.Merge(b)

	assert.Equal(t, 3, a.Tests.all)
	assert.Equal(t, 2, a.Tests.failed)
}
