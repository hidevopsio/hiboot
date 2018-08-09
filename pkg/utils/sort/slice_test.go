package sort

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSortByLen(t *testing.T) {
	s := []string{"a", "bb", "ccc", "d", "ee", "fff"}
	sorted := []string{"a", "d", "bb", "ee", "ccc", "fff"}
	SortByLen(s)
	assert.Equal(t, s, sorted)
}
