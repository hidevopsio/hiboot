package reflector

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

type Foo struct{}


func TestNewReflectType(t *testing.T) {
	foo := NewReflectType(Foo{})
	assert.NotEqual(t, nil, foo)
}
