package config

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/examples/cli/advanced/model"
)

type configuration struct {
	app.Configuration
}

func init() {
	app.AutoConfiguration(new(configuration))
}

func (c *configuration) Foo() *model.Foo {
	return new(model.Foo)
}

func (c *configuration) FooBar() *model.Foo {
	fb := new(model.Foo)
	fb.Name = "foobar"
	return fb
}
