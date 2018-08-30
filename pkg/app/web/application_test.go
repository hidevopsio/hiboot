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

package web_test

import (
	"fmt"
	"time"
	"errors"
	"testing"
	"net/http"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/model"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/starter/jwt"
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

// PATH /foo
type FooController struct{
	web.Controller
	jwtToken jwt.Token
}

type ExampleController struct{
	web.Controller
}

type InvalidController struct {}

func init() {
	log.SetLevel(log.DebugLevel)
}

func (c *FooController) Init(jwtToken jwt.Token) {
	c.jwtToken = jwtToken
}

func (c *FooController) Before()  {
	log.Debug("FooController.Before")

	c.Ctx.Next()
}

func (c *FooController) PostLogin(request *UserRequest) (response model.Response, err error)  {
	log.Debug("FooController.Login")

	// you make validate username and password first
	token, err := c.jwtToken.Generate(jwt.Map{
		"username": request.Username,
		"password": request.Password,
	}, 10, time.Minute)
	response = new(model.BaseResponse)
	response.SetData(token)

	return
}

// POST /
func (c *FooController) Post(request *FooRequest) (response model.Response, err error)  {
	log.Debug("FooController.Post")

	response = new(model.BaseResponse)
	if request.Name != "John" {
		return response, errors.New("only John is illegal name in this test")
	}
	response.SetData("Hello, " + request.Name)
	return
}


// GET /options/{options}
func (c *FooController) GetByOptions(options []string) (response model.Response)  {
	type settings struct {
		Options []string
	}
	response = new(model.BaseResponse)
	response.SetCode(http.StatusOK)
	response.SetMessage("success")
	response.SetData(&settings{Options: options})
	return
}


// GET /name/{name}
func (c *FooController) GetByName(name string) (response model.Response)  {
	type user struct {
		Name string
	}
	response = new(model.BaseResponse)
	response.SetCode(http.StatusOK)
	response.SetMessage("success")
	response.SetData(&user{Name: name})
	return
}

// GET /id/{id}
func (c *FooController) GetById(id int) string  {
	log.Debugf("FooController.Get by id: %v", id)
	return "hello"
}

// GET /hello
func (c *FooController) GetHello(ctx *web.Context) string  {
	log.Debug("FooController.GetHello")
	return "hello"
}

// PUT /id/{id}/name/{name}/age/{age}
func (c *FooController) PutByIdNameAge(id int, name string, age int) error {
	log.Debugf("FooController.Put %v %v %v", id, name, age)
	return nil
}

// PATCH /id/{id}
func (c *FooController) PatchById(id int) error {
	log.Debug("FooController.Patch")
	return nil
}

// DELETE /id/{id}
func (c *FooController) DeleteById(id int) error {
	log.Debug("FooController.Delete ", id)
	return nil
}

func (c *FooController) After()  {
	log.Debug("FooController.After")
}

// BarController
type BarController struct{
	jwt.Controller
}

func (c *BarController) Get(request *BarRequest) (response model.Response, err error)  {
	log.Debug("BarController.Get")
	response = new(model.BaseResponse)
	response.SetData("Hello, " + request.Name)

	return response, nil
}

type FoobarController struct {
	web.Controller
}

func (c *FoobarController) Post(request *FoobarRequestForm) (response model.Response, err error)  {
	response = new(model.BaseResponse)
	response.SetData("Hello, " + request.Name)

	return
}

func (c *FoobarController) Get(request *FoobarRequestParams) (response model.Response, err error) {
	response = new(model.BaseResponse)
	response.SetData("Hello, " + request.Name)

	return
}

// Define our controller, start with the name Foo, the first word of the Camelcase FooController is the controller name
// the lower cased foo will be the context mapping of the controller
// context mapping can be overwritten by FooController.ContextMapping
type HelloController struct{
	web.Controller
	ContextMapping string `value:"/"`
}

// Get hello
// The first word of method is the http method GET, the rest is the context mapping hello
// in this method, the name Get means that the method context mapping is '/'
func (c *HelloController) Get() string {
	return "hello"
}

// Get /all
func (c *HelloController) GetAll() {

	data := []struct{
		Name string
		Age int
	}{
		{
			Name: "John Doe",
			Age: 18,
		},
		{
			Name: "Zhang San",
			Age: 25,
		},
	}

	c.Ctx.ResponseBody("success", data)
}

func TestWebApplication(t *testing.T)  {
	app := web.NewTestApplication(t, new(HelloController), new(FooController), new(BarController), new(FoobarController))

	t.Run("should response 200 when GET /all", func(t *testing.T) {
		app.
			Request(http.MethodGet, "/all").
			Expect().Status(http.StatusOK).
			Body().Contains("Success").Contains("John Doe").Contains("Zhang San")
	})

	t.Run("should response 200 when GET /", func(t *testing.T) {
		app.
			Get("/").
			Expect().Status(http.StatusOK)
	})

	t.Run("should response 你好, 世界 with language zh-CN", func(t *testing.T) {
		// test cn-ZH
		app.Get("/").
			WithHeader("Accept-Language", "zh-CN").
			Expect().Status(http.StatusOK).
			Body().Contains("你好, 世界")
	})

	t.Run("should response Success with language en-US", func(t *testing.T) {
		// test en-US
		app.Get("/").
			WithHeader("Accept-Language", "en-US").
			Expect().
			Status(http.StatusOK).
			Body().Contains("Hello, World")
	})

	//t.Run("should pass health check", func(t *testing.T) {
	//	app.Get("/health").Expect().Status(http.StatusOK)
	//})

	t.Run("should login with username and password", func(t *testing.T) {
		app.Post("/foo/login").
			WithJSON(&UserRequest{Username: "johndoe", Password: "iHop91#15"}).
			Expect().Status(http.StatusOK).
			Body().NotEqual("")
	})

	t.Run("should failed to pass validation", func(t *testing.T) {
		app.Post("/foo/login").
			WithJSON(&UserRequest{Username: "johndoe"}).
			Expect().Status(http.StatusBadRequest)
	})

	t.Run("should return success after POST /foo", func(t *testing.T) {
		app.Post("/foo").
			WithJSON(&FooRequest{Name: "John"}).
			Expect().Status(http.StatusOK).
			Body().Contains("Hello")
	})

	t.Run("should return StatusInternalServerError after POST /foo with illegal name", func(t *testing.T) {
		app.Post("/foo").
			WithJSON(&FooRequest{Name: "Illegal name"}).
			Expect().Status(http.StatusInternalServerError)
	})


	t.Run("should return success after GET /foo/hello", func(t *testing.T) {
		app.Get("/foo/hello").
			WithJSON(&FooRequest{Name: "John"}).
			Expect().Status(http.StatusOK).
			Body().Equal("Hello, World")
	})


	t.Run("should return http.StatusUnauthorized after GET /bar", func(t *testing.T) {
		app.Get("/bar").
			Expect().Status(http.StatusUnauthorized)
	})

	t.Run("should return http.StatusInternalServerError when input form field validation failed", func(t *testing.T) {
		// test request form
		app.Post("/foobar").
			WithFormField("name", "John Doe").
			Expect().Status(http.StatusInternalServerError)
	})

	t.Run("should return (http.StatusOK on /foobar", func(t *testing.T) {
		//  test request query
		app.Get("/foobar").
			WithQuery("name", "John Doe").
			Expect().Status(http.StatusOK)
	})

	t.Run("should return http.StatusOK on /bar with jwt token", func(t *testing.T) {
		log.Println(io.GetWorkDir())
		jwtToken := jwt.NewJwtToken(&jwt.Properties{
			PrivateKeyPath: "config/ssl/app.rsa",
			PublicKeyPath: "config/ssl/app.rsa.pub",
		})
		// test jwt
		token, err := jwtToken.Generate(jwt.Map{
			"username": "johndoe",
			"password": "PA$$W0RD",
		}, 100, time.Millisecond)
		if err == nil {

			t := fmt.Sprintf("Bearer %v", token)

			app.Get("/bar").
				WithHeader("Authorization", t).
				Expect().Status(http.StatusOK)

			time.Sleep(2 * time.Second)

			app.Get("/bar").
				WithHeader("Authorization", t).
				Expect().Status(http.StatusUnauthorized)
		}
	})


	t.Run("should return http.StatusOK on /foo with PUT, PATCH, DELETE methods", func(t *testing.T) {
		app.Put("/foo/id/{id}/name/{name}/age/{age}").
			WithPath("id", 123456).
			WithPath("name", "Mike").
			WithPath("age", 18).
			Expect().Status(http.StatusOK)
	})

	t.Run("should Get foo by id", func(t *testing.T) {
		app.Get("/foo/id/{id}").
			WithPath("id", 123).
			Expect().Status(http.StatusOK)
	})

	t.Run("should Get foo by name", func(t *testing.T) {
		app.Get("/foo/options/{options}").
			WithPath("options", "mars,earth,mercury,jupiter").
			Expect().Status(http.StatusOK).
			Body().Contains("mercury")
	})

	t.Run("should Get foo by name", func(t *testing.T) {
		app.Get("/foo/name/{name}").
			WithPath("name", "Peter Phil").
			Expect().Status(http.StatusOK).
			Body().Contains("Peter Phil")
	})

	t.Run("should Patch foo by id", func(t *testing.T) {
		app.Patch("/foo/id/{id}").
			WithPath("id", 456).
			Expect().Status(http.StatusOK)
	})
	t.Run("should Delete foo by id", func(t *testing.T) {
		app.Delete("/foo/id/{id}").
			WithPath("id", 789).
			Expect().Status(http.StatusOK)
	})
}

//func TestInvalidController(t *testing.T)  {
//	ta := new(testApplication)
//	err := ta.Init(new(InvalidController))
//	err, ok := err.(*system.InvalidControllerError)
//	assert.Equal(t, ok, true)
//}

func TestNewApplication(t *testing.T) {

	web.RestController(new(ExampleController))

	wa := web.NewApplication()
	t.Run("should init web application", func(t *testing.T) {
		assert.NotEqual(t, nil, wa)
	})

	t.Run("should get interface name", func(t *testing.T) {
		type myInterface interface {

		}

		mi := new(myInterface)

		typ, err := reflector.GetType(mi)
		assert.Equal(t, nil, err)
		assert.Equal(t, "myInterface", typ.Name())
	})

	go wa.Run()
}

