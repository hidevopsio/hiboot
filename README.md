# Hiboot

[![Build Status](https://travis-ci.org/hidevopsio/hiboot.svg?branch=master)](https://travis-ci.org/hidevopsio/hiboot) 
[![codecov](https://codecov.io/gh/hidevopsio/hiboot/branch/master/graph/badge.svg)](https://codecov.io/gh/hidevopsio/hiboot)
[![Licensed under Apache License version 2.0](hiboot.svg)](https://www.apache.org/licenses/LICENSE-2.0)

## About

Hiboot is a cloud native web and cli application framework written in Go.

Hiboot is not trying to reinvent everything, it integrates popular libraries but make it simpler, easier to use.

With auto configuration, you can integrate any other libraries easily with dependency injection support.

## Overview

* Web MVC (Model-View-Controller)
* Auto Configuration, pre-create instance with properties configs for dependency injection
* Dependency injection with struct tag name **\`inject:""\`** or **Init** method
* Some useful utils, include enhanced reflection, struct copier, config file replacer, etc.


## Getting started

This section will show you how to create and run a simplest hiboot application. Let’s get started!

### Getting started with Hiboot web application

#### Get the source code

```bash
go get -u github.com/hidevopsio/hiboot

cd $GOPATH/src/github.com/hidevopsio/hiboot/examples/web/helloworld/


```

#### Sample code
 
Below is the simplest web application in Go.


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

#### Run web application

```bash
dep ensure

go run main.go
```

#### Testing the API by curl

```bash
curl http://localhost:8080/
```

```
Hello, world
```

### Getting started with Hiboot cli application

Writing Hiboot cli application is as simple as web application, you can take the advantage of dependency injection introduced by Hiboot.

e.g. flag tag dependency injection

```go

// declare main package
package main

// import cli starter and fmt
import "github.com/hidevopsio/hiboot/pkg/starter/cli"
import "fmt"

// define the command
type HelloCommand struct {
	// embedding cli.BaseCommand in each command
	cli.BaseCommand
	// inject (bind) flag to field 'To', so that it can be used on Run method, please note that the data type must be pointer
	To *string `flag:"name=to,shorthand=t,value=world,usage=e.g. --to=world or -t world"`
}

// Init constructor
func (c *HelloCommand) Init() {
	c.Use = "hello"
	c.Short = "hello command"
	c.Long = "run hello command for getting started"
}

// Run run the command
func (c *HelloCommand) Run(args []string) error {
	fmt.Printf("Hello, %v\n", *c.To)
	return nil
}

// main function
func main() {
	// create new cli application and run it
	cli.NewApplication(new(HelloCommand)).Run()
}

```

#### Run cli application

```bash
dep ensure

go run main.go
```

```bash
Hello, world
```

#### Build the cli application and run

```bash
go build
```

Let's get help

```bash
./hello --help
```

```bash
run hello command for getting started

Usage:
  hello [flags]

Flags:
  -h, --help        help for hello
  -t, --to string   e.g. --to=world or -t world (default "world")

```

Greeting to Hiboot

```bash
./hello --to Hiboot
```

```bash
Hello, Hiboot
```



