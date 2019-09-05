package actuator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfiguration(t *testing.T) {
	c := newConfiguration(&properties{})
	assert.NotEqual(t, nil, c)
}
