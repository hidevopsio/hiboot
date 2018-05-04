package main

import (
	"github.com/hidevopsio/hiboot/pkg/starter/web"
)

// Define our controller, start with the name Foo, the first word of the Camelcase FooController is the controller name
// the lower cased foo will be the context mapping of the controller
// context mapping can be overwritten by FooController.ContextMapping
type FooController struct{
	web.Controller
}

// Get hello
// The first word of method is the http method GET, the rest is the context mapping hello
func (c *FooController) GetHello(ctx *web.Context)  {

	ctx.Response("Success", "Hello, World")
}

func main()  {

	// create new web application
	app, err := web.NewApplication(&FooController{})

	// run the application
	if err == nil {
		app.Run()
	}
}