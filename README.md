# hiboot

[![Build Status](https://travis-ci.org/hidevopsio/hiboot.svg?branch=master)](https://travis-ci.org/hidevopsio/hiboot) 
[![codecov](https://codecov.io/gh/hidevopsio/hiboot/branch/master/graph/badge.svg)](https://codecov.io/gh/hidevopsio/hiboot)
[![Licensed under Apache License version 2.0](hiboot.svg)](https://www.apache.org/licenses/LICENSE-2.0)


'hiboot' is a web application framework written in Go. 

## Getting started

#### Get source code

```bash
go get -u github.com/hidevopsio/hiboot

cd $GOPATH/src/github.com/hidevopsio/hiboot/examples/web/helloworld/


```

#### The simplest web application in Go.


```go
// Line 1: main package
package main

// Line 2: import web starter from hiboot
import "github.com/hidevopsio/hiboot/pkg/starter/web"

// Line 3-5: RESTful Controller, derived from web.Controller. The context mapping of this controller is '/' by default
type Controller struct {
	web.Controller
}

// Line 6-8: Get method, the context mapping of this method is '/' by default
// the Method name Get means that the http request method is GET
func (c *Controller) Get() string {
	// response data
	return "Hello world"
}

// Line 9-11: main function
func main() {
	// create new web application and run it
	web.NewApplication(&Controller{}).Run()
}
```

### Let's run it

```bash
go run main.go
```

### Testing the API by curl

```bash
curl http://localhost:8080/
```

```
Hello, world
```

### Happy coding in go!


