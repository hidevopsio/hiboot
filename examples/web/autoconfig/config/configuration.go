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

func newConfiguration(properties *properties) *configuration  {
	return &configuration{properties: properties}
}

func init() {
	app.Register(newConfiguration)
}

type Foo struct {
	Name string `json:"name" value:"foo"`
}

type Bar struct {
	at.ContextAware

	Name string `json:"name" value:"bar"`
}

type FooWithError struct {
	at.ContextAware
	Name string `json:"name" value:"foo"`
}

func (c *configuration) FooWithError() (foo *FooWithError, err error) {
	err = ErrFoo
	return
}

func (c *configuration) Foo() *Foo {
	return &Foo{}
}

func (c *configuration) Bar(ctx context.Context) *Bar {
	if ctx.GetHeader("Authorization") == "fake" {
		ctx.StatusCode(http.StatusUnauthorized)
		return nil
	}
	return &Bar{Name: c.properties.Name}
}