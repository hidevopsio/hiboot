package config

import (
	"errors"
	"net/http"

	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
)

const Profile string = "config"

var ErrFoo = errors.New("foo with error")

type configuration struct {
	at.AutoConfiguration

	properties *properties
}

func newConfiguration(properties *properties) *configuration {
	return &configuration{properties: properties}
}

func init() {
	app.Register(newConfiguration)
}

type Foo struct {
	Name string `json:"name" value:"foo"`
}

type Bar struct {
	at.Scope `value:"request"`

	Name       string `json:"name" value:"bar"`
	FooBarName string `json:"fooBarName" value:"foobar"`

	baz *Baz
}

type FooBar struct {
	at.ContextAware

	Name string `json:"name" value:"foobar"`
}

type Baz struct {
	at.Scope `value:"prototype"`

	Name string `json:"name" value:"baz"`
}

type FooWithError struct {
	at.Scope `value:"request"`
	Name     string `json:"name" value:"foo"`
}

func (c *configuration) FooWithError() (foo *FooWithError, err error) {
	err = ErrFoo
	return
}

func (c *configuration) Foo() *Foo {
	return &Foo{}
}

func (c *configuration) Bar(ctx context.Context, foobar *FooBar) *Bar {
	if ctx.GetHeader("Authorization") == "fake" {
		ctx.StatusCode(http.StatusUnauthorized)
		return nil
	}
	return &Bar{Name: c.properties.Name, FooBarName: foobar.Name}
}

func (c *configuration) FooBar() *FooBar {
	return &FooBar{}
}

// Baz is a prototype scoped instance
func (c *configuration) Baz() *Baz {
	return &Baz{}
}
