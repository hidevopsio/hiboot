package foo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDummy(t *testing.T) {
	c := NewConfiguration()
	assert.NotEqual(t, nil, c)
}
