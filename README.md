# Hiboot - web/cli application framework 

<p align="center">
  <img src="https://github.com/hidevopsio/hiboot/blob/master/hiboot.png?raw=true" alt="hiboot">
</p>

<p align="center">
  <a href="https://travis-ci.org/hidevopsio/hiboot?branch=master">
    <img src="https://travis-ci.org/hidevopsio/hiboot.svg?branch=master" alt="Build Status"/>
  </a>
  <a href="https://codecov.io/gh/hidevopsio/hiboot">
    <img src="https://codecov.io/gh/hidevopsio/hiboot/branch/master/graph/badge.svg" />
  </a>
  <a href="https://opensource.org/licenses/Apache-2.0">
      <img src="https://img.shields.io/badge/License-Apache%202.0-green.svg" />
  </a>
  <a href="https://goreportcard.com/report/github.com/hidevopsio/hiboot">
      <img src="https://goreportcard.com/badge/github.com/hidevopsio/hiboot" />
  </a>
  <a href="https://godoc.org/github.com/hidevopsio/hiboot">
      <img src="https://godoc.org/github.com/golang/gddo?status.svg" />
  </a>
  <a href="https://gitter.im/hidevopsio/hiboot">
      <img src="https://img.shields.io/badge/GITTER-join%20chat-green.svg" />
  </a>
</p>

## About

Hiboot is a cloud native web and cli application framework written in Go.

Hiboot is not trying to reinvent everything, it integrates the popular libraries but make them simpler, easier to use. It borrowed some of the Spring features like dependency injection, aspect oriented programming, and auto configuration. You can integrate any other libraries easily by auto configuration with dependency injection support.

If you are a Java developer, you can start coding in Go without learning curve.

## Overview

* Web MVC (Model-View-Controller).
* Auto Configuration, pre-create instance with properties configs for dependency injection.
* Dependency injection with struct tag name **\`inject:""\`** or **Constructor** func.

## Introduction to Hiboot

One of the most significant feature of Hiboot is Dependency Injection. Hiboot implements JSR-330 standard.

Let's say that we have two implementations of AuthenticationService, below will explain how does Hiboot work.

```go
type AuthenticationService interface {
	Authenticate(credential Credential) error
}

type basicAuthenticationService struct {
}

func newBasicAuthenticationService() AuthenticationService {
	return &basicAuthenticationService{}
}

func (s *basicAuthenticationService) Authenticate(credential Credential) error {
	// business logic ...
	return nil
}

type oauth2AuthenticationService struct {
}

func newOauth2AuthenticationService() AuthenticationService {
	return &oauth2AuthenticationService{}
}

func (s *oauth2AuthenticationService) Authenticate(credential Credential) error {
	// business logic ...
	return nil
}

func init() {
	app.Register(newBasicAuthenticationService, newOauth2AuthenticationService)
}
```

### Field Injection

In Hiboot the injection into fields is triggered by **\`inject:""\`** struct tag. when inject tag is present
on a field, Hiboot tries to resolve the object to inject by the type of the field. If several implementations
of the same service interface are available, you have to disambiguate which implementation you want to be
injected. This can be done by naming the field to specific implementation.

```go
type userController struct {
	web.Controller

	BasicAuthenticationService AuthenticationService	`inject:""`
	Oauth2AuthenticationService AuthenticationService	`inject:""`
}

func newUserController() {
	return &userController{}
}

func init() {
	app.Register(newUserController)
}
```
### Constructor Injection

Although Field Injection is pretty convenient, but the Constructor Injection is the first-class citizen, we
usually advise people to use constructor injection as it has below advantages,

* It's testable, easy to implement unit test.
* Syntax validation, with syntax validation on most of the IDEs to avoid typo.
* No need to use a dedicated mechanism to ensure required properties are set.

```go
type userController struct {
	web.Controller

	basicAuthenticationService AuthenticationService
}

// Hiboot will inject the implementation of AuthenticationService
func newUserController(basicAuthenticationService AuthenticationService) {
	return &userController{
		basicAuthenticationService: basicAuthenticationService,
	}
}

func init() {
	app.Register(newUserController)
}
```

## Features

* **Apps**
    * cli - command line application
    * web - web application

* **Starters**
    * actuator - health check
    * locale - locale starter
    * logging - customized logging settings
    * jwt - jwt starter
    * grpc - grpc application starter

* **Tags** 
    * inject - inject generic instance into object
    * default - inject default value into struct object 
    * value - inject string value or references / variables into struct string field

* **Utils** 
    * cmap - concurrent map
    * copier - copy between struct
    * crypto - aes, base64, md5, and rsa encryption / decryption
    * gotest - go test util
    * idgen - twitter snowflake id generator
    * io - file io util
    * mapstruct - convert map to struct
    * replacer - replacing stuct field value with references or environment variables
    * sort - sort slice elements
    * str - string util enhancement util
    * validator - struct field validation 
       
and more features on the wey ...

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
import "github.com/hidevopsio/hiboot/pkg/app/web"

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

// import cli starter and fmt
import (
	"fmt"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/cli"
)

// define the command
type rootCommand struct {
	cli.RootCommand
	To string
}

func newRootCommand() *rootCommand {
	c := new(rootCommand)
	c.Use = "hello"
	c.Short = "hello command"
	c.Long = "run hello command for getting started"
	c.Example = `
hello -h : help
hello -t John : say hello to John
`
	c.PersistentFlags().StringVarP(&c.To, "to", "t", "world", "e.g. --to=world or -t world")
	return c
}

// Run run the command
func (c *rootCommand) Run(args []string) error {
	fmt.Printf("Hello, %v\n", c.To)
	return nil
}

// main function
func main() {
	// create new cli application and run it
	cli.NewApplication(newRootCommand).
		SetProperty(app.PropertyBannerDisabled, true).
		Run()
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

### Dependency injection in Go

Dependency injection is a concept valid for any programming language. The general concept behind dependency injection is called Inversion of Control. According to this concept a struct should not configure its dependencies statically but should be configured from the outside.

Dependency Injection design pattern allows us to remove the hard-coded dependencies and make our application loosely coupled, extendable and maintainable.

A Go struct has a dependency on another struct, if it uses an instance of this struct. We call this a struct dependency. For example, a struct which accesses a user controller has a dependency on user service struct.

Ideally Go struct should be as independent as possible from other Go struct. This increases the possibility of reusing these struct and to be able to test them independently from other struct.

The following example shows a struct which has no hard dependencies.

```go
package main

import (
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/model"
	"github.com/hidevopsio/hiboot/pkg/starter/jwt"
	"time"
)

// This example shows that token is injected through method Init,
// once you imported "github.com/hidevopsio/hiboot/pkg/starter/jwt",
// token jwt.Token will be injectable.
func Example() {
	// the web application entry
	web.NewApplication().Run()
}

// PATH: /login
type loginController struct {
	web.Controller

	token jwt.Token
}

type userRequest struct {
	// embedded field model.RequestBody mark that userRequest is request body
	model.RequestBody
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func init() {
	// Register Rest Controller through constructor newLoginController
	app.Register(newLoginController)
}

// newLoginController inject token through the argument token jwt.Token on constructor
// the dependency token is auto configured in jwt starter, see https://github.com/hidevopsio/hiboot/tree/master/pkg/starter/jwt
func newLoginController(token jwt.Token) *loginController {
	return &loginController{
		token: token,
	}
}

// Post /
// The first word of method is the http method POST, the rest is the context mapping
func (c *loginController) Post(request *userRequest) (response model.Response, err error) {
	jwtToken, _ := c.token.Generate(jwt.Map{
		"username": request.Username,
		"password": request.Password,
	}, 30, time.Minute)

	response = new(model.BaseResponse)
	response.SetData(jwtToken)

	return
}

```

## Community Contributions Guide

Thank you for considering contributing to the Hiboot framework, The contribution guide can be found [here](CONTRIBUTING.md).