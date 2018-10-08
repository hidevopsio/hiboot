package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig(t *testing.T) {

	c := newConfiguration()

	assert.NotEqual(t, nil, c.Foo())
	assert.NotEqual(t, nil, c.FooBar())

}
