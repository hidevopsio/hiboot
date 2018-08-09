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
// Line 1: main package
package main

// Line 2-3: import cli starter and fmt
import "github.com/hidevopsio/hiboot/pkg/starter/cli"
import "fmt"


type SampleCommand struct {
	cli.BaseCommand
	// inject flag to profile so that you can use it on Run method
	Name string `flag:"shorthand=n,value=world,usage=e.g. --name=world or -n world"`
}

func init() {
  cli.AddCommand("root", new(SampleCommand))
}

func (c *SampleCommand) Init()  {
	c.Use = "sample"
	c.Short = "sample command"
	c.Long = "run sample command for getting started"
}

func (c *SampleCommand) Run(args []string) error {
	fmt.Printf("Hello, %v\n", c.Name)
	return nil
}

func main() {
	// create new cli application and run it
	cli.NewApplication().Run()
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
  main [flags]

Flags:
  -h, --help        help for main
  -t, --to string   e.g. --to=world or -t world (default "world")

```

Greeting to Hiboot

```bash
./hello --to Hiboot
```

```bash
Hello, Hiboot
```



