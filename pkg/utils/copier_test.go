package utils

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestCopier(t *testing.T) {

	foo := &Foo{
		Name: "foo",
	}

	bar := &Bar{}

	err := Copy(bar, foo)

	assert.Equal(t, nil, err)
	assert.Equal(t, bar.Name, foo.Name)

}
