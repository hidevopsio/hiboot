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

```

### Let's run it

```bash
go run main.go
```

### testing

```bash
curl http://localhost:8080
```

### Happy coding in go!