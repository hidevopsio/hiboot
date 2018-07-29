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
	"github.com/hidevopsio/hiboot/pkg/model"
)

type UserRequest struct {
	model.RequestBody
	Username string	`validate:"required"`
	Password string	`validate:"required"`
}

type FooRequest struct {
	model.RequestBody
	Name string
}

type BarRequest struct {
	model.RequestParams
	Name string
}

type FoobarRequestForm struct {
	model.RequestForm
	Name string
}

type FoobarRequestParams struct {
	model.RequestParams
	Name string
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

func (c *FooController) PostLogin(request *UserRequest) (response model.Response, err error)  {
	log.Debug("FooController.Login")

	// you make validate username and password first

	jwtToken, err := GenerateJwtToken(JwtMap{
		"username": request.Username,
		"password": request.Password,
	}, 10, time.Minute)

	response.Data = jwtToken

	return response, err
}

func (c *FooController) Post(request *FooRequest) (response model.Response, err error)  {
	log.Debug("FooController.Post")

	response.Data = "Hello, " + request.Name

	return response, nil
}

func (c *FooController) Get() string  {
	log.Debug("FooController.Post")
	return "hello"
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

func (c *BarController) Get(request *BarRequest) (response model.Response, err error)  {
	log.Debug("BarController.Get")

	response.Data = "Hello, " + request.Name

	return response, nil
}

type FoobarController struct {
	Controller
}

func (c *FoobarController) Post(request *FoobarRequestForm) (response model.Response, err error)  {

	response.Data = "Hello, " + request.Name

	return
}

func (c *FoobarController) Get(request *FoobarRequestParams) (response model.Response, err error) {

	response.Data = "Hello, " + request.Name

	return
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
func (c *HelloController) Get() string {
	return "hello"
}

func TestWebApplication(t *testing.T)  {
	app := NewTestApplication(t, new(HelloController), new(FooController), new(BarController), new(FoobarController))

	t.Run("should response 200 when GET /", func(t *testing.T) {
		app.
			Get("/").
			Expect().Status(http.StatusOK)
	})

	t.Run("should response 你好，世界 with language zh-CN", func(t *testing.T) {
		// test cn-ZH
		app.Get("/").
			WithHeader("Accept-Language", "zh-CN").
			Expect().Status(http.StatusOK)
	})

	t.Run("should response Success with language en-US", func(t *testing.T) {
		// test en-US
		app.Request("GET", "/").
			WithHeader("Accept-Language", "en-US").
			Expect().
			Status(http.StatusOK).
			Body().Contains("Hello, World")
	})

	t.Run("should pass health check", func(t *testing.T) {
		app.Get("/health").Expect().Status(http.StatusOK)
	})

	t.Run("should login with username and password", func(t *testing.T) {
		app.Post("/foo/login").
			WithJSON(&UserRequest{Username: "johndoe", Password: "iHop91#15"}).
			Expect().Status(http.StatusOK)
	})

	t.Run("should failed to pass validation", func(t *testing.T) {
		app.Post("/foo/login").
			WithJSON(&UserRequest{Username: "johndoe"}).
			Expect().Status(http.StatusBadRequest)
	})

	app.Post("/foo").
		WithJSON(&FooRequest{Name: "John"}).
		Expect().Status(http.StatusOK)

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
