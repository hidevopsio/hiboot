package utils

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestLowerFirst(t *testing.T) {
	s := "Foo"

	ns := LowerFirst(s)

	assert.Equal(t, "foo", ns)
}

func TestUpperFirst(t *testing.T) {
	s := "foo"

	ns := UpperFirst(s)

	assert.Equal(t, "Foo", ns)
}


func TestStringInSlice(t *testing.T) {
	s := []string{
		"foo",
		"bar",
		"baz",
	}

	assert.Equal(t, true, StringInSlice("bar", s))
}
