# hiboot

'hiboot' is a cloud native web application framework.

### Get the code

```bash
go get -u github.com/hidevopsio/hiboot
```


### The simplest go web application

```go
package main

import (
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"github.com/hidevopsio/hiboot/pkg/log"
	"os"
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
	if err != nil {

		log.Error(err)
		os.Exit(1)
	}
	
	// run the application
	app.Run()
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