package main

import (
	"github.com/hidevopsio/hiboot/pkg/starter/web"
)


type FooController struct{
}

type Controllers struct{
	Foo *FooController
}

func (c *FooController) GetHello(ctx *web.Context)  {

	ctx.Response("Success", "Hello, World")
}

func main()  {

	controllers := &Controllers{}
	app, err := web.NewApplication(controllers)
	if err == nil {
		app.Run()
	}
}