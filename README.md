# hiboot

'hiboot' is a cloud native web application framework.

### Get the code

```bash
go get -u github.com/hidevopsio/hiboot
```


### The simplest go web application in the world



```go

package main

import (
	"github.com/hidevopsio/hiboot/pkg/starter/web"
)

// Define our controller, start with the name Foo, the first word of the Camelcase FooController is the controller name
// the lower case foo will be the context mapping of the controller
// context mapping can also be overwritten by FooController.ContextMapping
// if the controller name is a single word Controller, then the context mapping will be '/'
type Controller struct{
	web.Controller
}

// Get hello
// the first word of method is the http method GET, the rest is the context mapping hello
// if the method name is a single word Get, the the context mapping will be '/'
func (c *Controller) Get(ctx *web.Context)  {

	ctx.Response("Success", "Hello, World")
}

func main()  {

	// create new web application
	app, err := web.NewApplication(&Controller{})

	// run the application
	if err == nil {
		app.Run()
	}
}

```

### Let's run it

```bash
go run main.go
```

### testing

```bash
curl http://localhost:8080/foo/hello
```

### Happy coding in go!