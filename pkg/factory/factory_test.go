package factory

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFactory(t *testing.T) {

}

func TestDependencies(t *testing.T) {
	type config struct {
		Deps Deps
	}

	c := new(config)
	deps := []string{"world"}
	c.Deps.Set("hello", deps)
	assert.Equal(t, deps, c.Deps.Get("hello"))
}
