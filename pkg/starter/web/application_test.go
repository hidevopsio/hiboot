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
	"fmt"
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

type ExampleController struct{
	Controller
}

type InvalidController struct {}

func init() {
	log.SetLevel(log.DebugLevel)
	utils.ChangeWorkDir("../../../")
}

func (c *FooController) Before()  {
	log.Debug("FooController.Before")

	c.Ctx.Next()
}

func (c *FooController) PostLogin()  {
	log.Debug("FooController.Login")
	userRequest := &UserRequest{}
	if c.Ctx.RequestBody(userRequest) == nil {
		jwtToken, err := GenerateJwtToken(JwtMap{
			"username": userRequest.Username,
			"password": userRequest.Password,
		}, 10, time.Minute)

		log.Debugf("token: %v", jwtToken)

		if err == nil {
			c.Ctx.ResponseBody("Success", jwtToken)
		} else {
			c.Ctx.ResponseError(err.Error(), http.StatusInternalServerError)
		}
	}
}

func (c *FooController) Post()  {
	log.Debug("FooController.Post")

	foo := &FooRequest{}
	if c.Ctx.RequestBody(foo) == nil {
		c.Ctx.ResponseBody("Success", &FooResponse{Greeting: "hello, " + foo.Name})
	}
}


func (c *FooController) Put()  {
	log.Debug("FooController.Put")
}

func (c *FooController) Patch()  {
	log.Debug("FooController.Patch")
}

func (c *FooController) Delete()  {
	log.Debug("FooController.Delete")
}

func (c *FooController) After()  {
	log.Debug("FooController.After")

}

// BarController
type BarController struct{
	JwtController
}

func (c *BarController) Get()  {
	log.Debug("BarController.Get")

	c.Ctx.ResponseBody("Success", &Bar{Greeting: "hello bar"})

}

type FoobarController struct {
	Controller
}

func (c *FoobarController) Post()  {
	foo := &FooRequest{}
	err := c.Ctx.RequestForm(foo)
	if err == nil {
		c.Ctx.ResponseBody("Success", &FooResponse{Greeting: "hello, " + foo.Name})
	} else {
		c.Ctx.ResponseError(err.Error(), http.StatusInternalServerError)
	}
}

func (c *FoobarController) Get()  {
	foo := &FooRequest{}
	if c.Ctx.RequestParams(foo) == nil {
		c.Ctx.ResponseBody("Success", &FooResponse{Greeting: "hello, " + foo.Name})
	}
}

// Define our controller, start with the name Foo, the first word of the Camelcase FooController is the controller name
// the lower cased foo will be the context mapping of the controller
// context mapping can be overwritten by FooController.ContextMapping
type HelloController struct{
	Controller
	ContextMapping string `value:"/"`
}

// Get hello
// The first word of method is the http method GET, the rest is the context mapping hello
// in this method, the name Get means that the method context mapping is '/'
func (c *HelloController) Get()  {

	c.Ctx.ResponseBody("success", "Hello, World")
}

func TestWebApplication(t *testing.T)  {
	app := NewTestApplication(t, new(HelloController), new(FooController), new(BarController), new(FoobarController))

	t.Run("should response 200 when GET /", func(t *testing.T) {
		app.
			Get("/").
			Expect().Status(http.StatusOK).
			Body().Contains("Success")
	})

	t.Run("should response 成功 with language zh-CN", func(t *testing.T) {
		// test cn-ZH
		app.Get("/").
			WithHeader("Accept-Language", "zh-CN").
			Expect().Status(http.StatusOK).
			Body().Contains("成功")
	})

	t.Run("should response Success with language en-US", func(t *testing.T) {
		// test en-US
		app.Request("GET", "/").
			WithHeader("Accept-Language", "en-US").
			Expect().
			Status(http.StatusOK).
			Body().Contains("Success")
	})

	app.Get("/health").Expect().Status(http.StatusOK)

	app.Post("/foo/login").
		WithJSON(&UserRequest{Username: "johndoe", Password: "iHop91#15"}).
		Expect().Status(http.StatusOK).
		Body().Contains("Success")

	app.Post("/foo").
		WithJSON(&FooRequest{Name: "John"}).
		Expect().Status(http.StatusOK).
		Body().Contains("Success")

	app.Get("/bar").
		Expect().Status(http.StatusUnauthorized)

	// test request form
	app.Post("/foobar").
		WithFormField("name", "John Doe").
		Expect().Status(http.StatusInternalServerError)

	//  test request query
	app.Get("/foobar").
		WithQuery("name", "John Doe").
		Expect().Status(http.StatusOK)

	// test jwt
	pt, err := GenerateJwtToken(JwtMap{
		"username": "johndoe",
		"password": "PA$$W0RD",
	}, 100, time.Millisecond)
	if err == nil {

		t := fmt.Sprintf("Bearer %v", string(*pt))

		app.Get("/bar").
			WithHeader("Authorization", t).
			Expect().Status(http.StatusOK)

		time.Sleep(2 * time.Second)

		app.Get("/bar").
			WithHeader("Authorization", t).
			Expect().Status(http.StatusUnauthorized)
	}

	app.Put("/foo").Expect().Status(http.StatusOK)
	app.Patch("/foo").Expect().Status(http.StatusOK)
	app.Delete("/foo").Expect().Status(http.StatusOK)
}

func TestInvalidController(t *testing.T)  {
	ta := new(TestApplication)
	err := ta.Init(new(InvalidController))
	err, ok := err.(*system.InvalidControllerError)
	assert.Equal(t, ok, true)
}

func TestNewApplication(t *testing.T) {
	var app *Application

	//t.Run("should return no controller error", func(t *testing.T) {
	//	NewApplication()
	//})

	Add(new(ExampleController))

	wa := NewApplication()
	t.Run("should init web application", func(t *testing.T) {
		log.Debugf("### web application: %v", app)
		assert.NotEqual(t, nil, wa.app)
	})

	t.Run("should get application config", func(t *testing.T) {
		config := wa.Config()
		assert.NotEqual(t, nil, config)
		assert.NotEqual(t, "", config.App.Name)
	})


}
