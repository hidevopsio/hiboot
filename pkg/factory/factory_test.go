package factory

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type config struct {
	Deps Deps
}

func (c *config) Hello() {

}

func TestFactory(t *testing.T) {

}
func TestDependencies(t *testing.T) {

	c := new(config)
	deps := []string{"world"}
	c.Deps.Set("Hello", deps)
	assert.Equal(t, deps, c.Deps.Get("Hello"))

	c.Deps.Set(c.Hello, deps)
	assert.Equal(t, deps, c.Deps.Get("Hello"))

	c.Deps.Set(nil, deps)
}
