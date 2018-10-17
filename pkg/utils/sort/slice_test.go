package sort

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSortByLen(t *testing.T) {
	s := []string{"a", "bb", "ccc", "d", "ee", "fff"}
	sorted := []string{"a", "d", "bb", "ee", "ccc", "fff"}
	ByLen(s)
	assert.Equal(t, s, sorted)
}
