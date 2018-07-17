// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package web

import (
	"testing"
	"github.com/hidevopsio/hiboot/pkg/log"
	"net/http"
	"time"
	"github.com/hidevopsio/hiboot/pkg/system"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/utils"
)

type UserRequest struct {
	Username string
	Password string
}

type FooRequest struct {
	Name string
}

type FooResponse struct {
	Greeting string
}

type Bar struct {
	Name string
	Greeting string
}

type FooController struct{
	Controller
}

type InvalidController struct {}

func init() {
	utils.ChangeWorkDir("../../../")
	Add(new(FooController))
	Add(new(BarController))
}

func (c *FooController) Before(ctx *Context)  {
	log.Debug("FooController.Before")

	ctx.Next()
}

func (c *FooController) PostLogin(ctx *Context)  {
	log.Debug("FooController.SayHello")

	userRequest := &UserRequest{}
	if ctx.RequestBody(userRequest) == nil {
		jwtToken, err := GenerateJwtToken(JwtMap{
			"username": userRequest.Username,
			"password": userRequest.Password,
		}, 10, time.Minute)

		log.Debugf("token: %v", jwtToken)

		if err == nil {
			ctx.ResponseBody("Success", jwtToken)
		} else {
			ctx.ResponseError(err.Error(), http.StatusInternalServerError)
		}
	}
}

func (c *FooController) PostSayHello(ctx *Context)  {
	log.Debug("FooController.SayHello")

	foo := &FooRequest{}
	if ctx.RequestBody(foo) == nil {
		ctx.ResponseBody("Success", &FooResponse{Greeting: "hello, " + foo.Name})
	}
}


func (c *FooController) Put(ctx *Context)  {
	log.Debug("FooController.Put")
}

func (c *FooController) Patch(ctx *Context)  {
	log.Debug("FooController.Patch")
}

func (c *FooController) Delete(ctx *Context)  {
	log.Debug("FooController.Delete")
}

func (c *FooController) After(ctx *Context)  {
	log.Debug("FooController.After")

}

// BarController
type BarController struct{
	JwtController
}

func (c *BarController) GetSayHello(ctx *Context)  {
	log.Debug("BarController.SayHello")

	ctx.ResponseBody("Success", &Bar{Greeting: "hello bar"})

}


// Define our controller, start with the name Foo, the first word of the Camelcase FooController is the controller name
// the lower cased foo will be the context mapping of the controller
// context mapping can be overwritten by FooController.ContextMapping
type HelloController struct{
	Controller
}

func (c *HelloController) Init()  {
	c.ContextMapping = "/"
}

// Get hello
// The first word of method is the http method GET, the rest is the context mapping hello
// in this method, the name Get means that the method context mapping is '/'
func (c *HelloController) Get(ctx *Context)  {

	ctx.ResponseBody("success", "Hello, World")
}

func TestHelloWorld(t *testing.T)  {

	// create new test server
	NewTestApplication(t, new(HelloController)).
		Get("/").
		Expect().Status(http.StatusOK).
		Body().Contains("Success")
}

func TestHelloWorldLocale(t *testing.T)  {
	// create new test server
	ta := NewTestApplication(t, new(HelloController))

	// test cn-ZH
	ta.Get("/").
		WithHeader("Accept-Language", "zh-CN").
		Expect().Status(http.StatusOK).
		Body().Contains("成功")

	// test en-US
	ta.Request("GET", "/").
		WithHeader("Accept-Language", "en-US").
		Expect().
		Status(http.StatusOK).
		Body().Contains("Success")
}

func TestWebApplication(t *testing.T)  {

	ta := NewTestApplication(t)

	ta.Get("/health").Expect().Status(http.StatusOK)

	ta.Post("/foo/login").
		WithJSON(&UserRequest{Username: "johndoe", Password: "iHop91#15"}).
		Expect().Status(http.StatusOK).
		Body().Contains("Success")

	ta.Post("/foo/sayHello").
		WithJSON(&FooRequest{Name: "John"}).
		Expect().Status(http.StatusOK).
		Body().Contains("Success")

	ta.Get("/bar/sayHello").
		Expect().Status(http.StatusUnauthorized)

	ta.Put("/foo").Expect().Status(http.StatusOK)
	ta.Patch("/foo").Expect().Status(http.StatusOK)
	ta.Delete("/foo").Expect().Status(http.StatusOK)
}

func TestInvalidController(t *testing.T)  {
	ta := new(TestApplication)
	err := ta.Init(new(InvalidController))
	err, ok := err.(*system.InvalidControllerError)
	assert.Equal(t, ok, true)
}

func TestNewApplication(t *testing.T) {
	var wa *Application
	t.Run("should init web application", func(t *testing.T) {
		wa = NewApplication(new(FooController))
		assert.NotEqual(t, nil, wa.app)
	})

	t.Run("should get application config", func(t *testing.T) {
		config := wa.Config()
		assert.NotEqual(t, "", config.App.Name)
	})

	//t.Run("should not init application without controller", func(t *testing.T) {
	//	wa := NewApplication()
	//	assert.NotEqual(t, nil, wa.app)
	//})
}
