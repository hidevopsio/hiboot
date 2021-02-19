package config

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/at"
)

const Profile string = "config"

type configuration struct {
	at.AutoConfiguration

	properties *properties
}

func newConfiguration(properties *properties) *configuration  {
	return &configuration{properties: properties}
}

func init() {
	app.Register(newConfiguration)
}

type Foo struct {
	Name string `json:"name"`
}

type Bar struct {
	Name string `json:"name"`
}

func (c *configuration) Foo(bar *Bar) *Foo {
	return &Foo{Name: bar.Name}
}

func (c *configuration) Bar() *Bar {
	return &Bar{Name: c.properties.Name}
}