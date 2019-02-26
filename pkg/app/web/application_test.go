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
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/app/web/context"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/model"
	_ "hidevops.io/hiboot/pkg/starter/actuator"
	"hidevops.io/hiboot/pkg/starter/jwt"
	_ "hidevops.io/hiboot/pkg/starter/locale"
	_ "hidevops.io/hiboot/pkg/starter/logging"
	"hidevops.io/hiboot/pkg/utils/io"
	"hidevops.io/hiboot/pkg/utils/reflector"
	"net/http"
	"testing"
	"time"
)

type UserRequest struct {
	model.RequestBody
	Username string `validate:"required"`
	Password string `validate:"required"`
}

type FooRequest struct {
	at.RequestBody
	Name string
}

type BarRequest struct {
	at.RequestParams
	Name string
}

type FoobarRequestForm struct {
	at.RequestForm
	Name string
}

type FoobarRequestParams struct {
	at.RequestParams
	Name string
}

type Bar struct {
	Name     string
	Greeting string
}

func newBar() *Bar {
	return &Bar{}
}

// PATH /foo
type FooController struct {
	at.RestController
	token jwt.Token
}

func newFooController(token jwt.Token) *FooController {
	return &FooController{
		token: token,
	}
}

type ExampleController struct {
	at.RestController
}

type InvalidController struct{}

type FooBar struct {
	Name string
}

type FooBarService struct {
	fooBar *FooBar
}

func newFooBarService(fooBar *FooBar) *FooBarService {
	return &FooBarService{
		fooBar: fooBar,
	}
}

func (s *FooBarService) FooBar() *FooBar {
	return s.fooBar
}

func init() {
	log.SetLevel(log.DebugLevel)
	app.Register(&FooBar{Name: "fooBar"})
	app.Register(newFooBarService)
	app.Register(newHelloContextAware)
}

func (c *FooController) Before(ctx context.Context) {
	log.Debug("FooController.Before")

	ctx.Next()
}

func (c *FooController) PostLogin(request *UserRequest) (response model.Response, err error) {
	log.Debug("FooController.Login")

	// you make validate username and password first
	token, err := c.token.Generate(jwt.Map{
		"username": request.Username,
		"password": request.Password,
	}, 10, time.Minute)
	response = new(model.BaseResponse)
	response.SetData(token)

	return
}

// POST /
func (c *FooController) Post(request *FooRequest) (response model.Response, err error) {
	log.Debug("FooController.Post")

	response = new(model.BaseResponse)
	if request.Name != "John" {
		return response, errors.New("only John is illegal name in this test")
	}
	response.SetData("Hello, " + request.Name)
	return
}

// GET /options/{options}
func (c *FooController) GetByOptions(options []string) (response model.Response) {
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
func (c *FooController) GetByName(name string) (response model.Response) {
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
func (c *FooController) GetById(id int) string {
	log.Debugf("FooController.Get by id: %v", id)
	return "hello"
}

// GET /hello
func (c *FooController) GetHello(ctx context.Context) string {
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

func (c *FooController) GetBar(bar string) {
	log.Debugf("FooController.GetBar: %v", bar)
}

func (c *FooController) Options() {
	log.Debug("FooController.Options")
}

// GET /invalid
func (c *FooController) GetInvalid() (err error) {
	return fmt.Errorf("response invalid")
}

// GET /integer
func (c *FooController) GetInteger() (i int) {
	return 123
}

// GET /intPointer
func (c *FooController) GetIntPointer() (ip *int) {
	iv := 123
	ip = &iv
	return
}

// GET /intNullPointer
func (c *FooController) GetIntNilPointer() (ip *int) {
	return
}

func (c *FooController) GetError() error {
	return fmt.Errorf("buildconfigs.mio.io \"hello-world\" not found")
}

type FooRequestBody struct {
	at.RequestBody
	Name string `json:"name"`
}

func (c *FooController) GetRequestBody(request *FooRequestBody) string {
	return request.Name
}

type FooRequestForm struct {
	at.RequestForm
	Name string `json:"name"`
}

func (c *FooController) GetRequestForm(request *FooRequestForm) string {
	return request.Name
}

type FooRequestParams struct {
	at.RequestParams
	Name string `json:"name"`
}

func (c *FooController) GetRequestParams(request *FooRequestParams) string {
	return request.Name
}

type HelloContextAware struct {
	at.ContextAware

	context context.Context
}

func newHelloContextAware(context context.Context) *HelloContextAware {
	return &HelloContextAware{context: context}
}

func (c *FooController) GetContext(hca *HelloContextAware) string {
	return "testing context aware dependency injection"
}

func (c *FooController) GetErr() float32 {
	return 0.01
}

func (c *FooController) After() {
	log.Debug("FooController.After")
}

// BarController
type BarController struct {
	at.JwtRestController
}

func newBarController() *BarController {
	return &BarController{}
}

func (c *BarController) Get(request *BarRequest) (response model.Response, err error) {
	log.Debug("BarController.Get")
	response = new(model.BaseResponse)
	response.SetData("Hello, " + request.Name)

	return response, nil
}

type FoobarController struct {
	at.RestController
}

func newFoobarController() *FoobarController {
	return &FoobarController{}
}

func (c *FoobarController) Post(request *FoobarRequestForm) (response model.Response, err error) {
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
type HelloController struct {
	at.RestController
	fooBar *FooBar
}

func newHelloController(fooBar *FooBar) *HelloController {
	return &HelloController{
		fooBar: fooBar,
	}
}

// Get hello
// The first word of method is the http method GET, the rest is the context mapping hello
// in this method, the name Get means that the method context mapping is '/'
func (c *HelloController) Get() string {
	return "hello"
}

func (c *HelloController) GetWorld(ctx context.Context) {
	ctx.HTML("<h1>Hello World</h1>")
}

// Get /all
func (c *HelloController) GetAll(ctx context.Context) {

	data := []struct {
		Name string
		Age  int
	}{
		{
			Name: "John Doe",
			Age:  18,
		},
		{
			Name: "Zhang San",
			Age:  25,
		},
	}

	ctx.ResponseBody("success", data)
}

// GetMap
func (c *HelloController) GetMap() map[string]interface{} {
	response := make(map[string]interface{})
	response["message"] = "Hello from Map"
	return response
}

// Define our controller, start with the name Foo, the first word of the Camelcase FooController is the controller name
// the lower cased foo will be the context mapping of the controller
// context path can be overwritten by the tag value of annotation at.ContextPath
type HelloViewController struct {
	at.RestController
	at.ContextPath `value:"/"`
	fooBar         *FooBar
}

func newHelloViewController(fooBar *FooBar) *HelloViewController {
	return &HelloViewController{
		fooBar: fooBar,
	}
}

// Get hello
// The first word of method is the http method GET, the rest is the context mapping hello
// in this method, the name Get means that the method context mapping is '/'
func (c *HelloViewController) Get(ctx context.Context) {
	ctx.View("index.html")
}

// AnyTest
func (c *HelloViewController) AnyTest() string {
	return "Greeting from any method"
}

func TestWebViewApplicationWithProperties(t *testing.T) {
	testApp := web.NewTestApp(newHelloViewController).
		SetProperty("web.view.enabled", true).
		Run(t)
	t.Run("should response 200 when GET /", func(t *testing.T) {
		testApp.
			Get("/").
			Expect().Status(http.StatusOK)
	})
}

func TestWebViewApplicationWithArgs(t *testing.T) {
	testApp := web.NewTestApp(newHelloViewController).
		SetProperty("server.port", 8080).
		SetProperty("web.view.enabled", true).
		Run(t)
	t.Run("should response 200 when GET /", func(t *testing.T) {
		testApp.
			Get("/").
			Expect().Status(http.StatusOK)
	})

	t.Run("should response 200 when GET /test", func(t *testing.T) {
		testApp.
			Get("/test").
			Expect().Status(http.StatusOK)
	})
}

func TestWebApplication(t *testing.T) {
	testApp := web.RunTestApplication(t, newHelloController, newFooController, newBarController, newFoobarController)

	t.Run("should response 200 when GET /hello/all", func(t *testing.T) {
		testApp.
			Request(http.MethodGet, "/hello/all").
			Expect().Status(http.StatusOK).
			Body().Contains("Success").Contains("John Doe").Contains("Zhang San")
	})

	t.Run("should response 200 when GET /hello", func(t *testing.T) {
		testApp.
			Get("/hello/").
			Expect().Status(http.StatusOK)
	})

	t.Run("should response 200 when GET /map", func(t *testing.T) {
		testApp.
			Get("/hello/map").
			Expect().Status(http.StatusOK)
	})

	t.Run("should response 200 when GET /hello/world", func(t *testing.T) {
		testApp.
			Get("/hello/world").
			Expect().Status(http.StatusOK)
	})

	t.Run("should response 你好, 世界 with language zh-CN", func(t *testing.T) {
		// test cn-ZH
		testApp.Get("/hello").
			WithHeader("Accept-Language", "zh-CN").
			Expect().Status(http.StatusOK).
			Body().Contains("你好, 世界")
	})

	t.Run("should response Success with language en-US", func(t *testing.T) {
		// test en-US
		testApp.Get("/hello").
			WithHeader("Accept-Language", "en-US").
			Expect().
			Status(http.StatusOK).
			Body().Contains("Hello, World")
	})

	//t.Run("should pass health check", func(t *testing.T) {
	//	app.Get("/health").Expect().Status(http.StatusOK)
	//})

	t.Run("should login with username and password", func(t *testing.T) {
		testApp.Post("/foo/login").
			WithJSON(&UserRequest{Username: "johndoe", Password: "iHop91#15"}).
			Expect().Status(http.StatusOK).
			Body().NotEqual("")
	})

	t.Run("should failed to pass validation", func(t *testing.T) {
		testApp.Post("/foo/login").
			WithJSON(&UserRequest{Username: "johndoe"}).
			Expect().Status(http.StatusBadRequest)
	})

	t.Run("should return success after POST /foo", func(t *testing.T) {
		testApp.Post("/foo").
			WithJSON(&FooRequest{Name: "John"}).
			Expect().Status(http.StatusOK).
			Body().Contains("Hello")
	})

	t.Run("should return StatusInternalServerError after POST /foo with illegal name", func(t *testing.T) {
		testApp.Post("/foo").
			WithJSON(&FooRequest{Name: "Illegal name"}).
			Expect().Status(http.StatusInternalServerError)
	})

	t.Run("should return success after GET /foo/hello", func(t *testing.T) {
		testApp.Get("/foo/hello").
			WithJSON(&FooRequest{Name: "John"}).
			Expect().Status(http.StatusOK).
			Body().Equal("Hello, World")
	})

	t.Run("should parse request body GET /foo/requestBody", func(t *testing.T) {
		testApp.Get("/foo/requestBody").
			WithJSON(&FooRequestBody{Name: "foo"}).
			Expect().Status(http.StatusOK).
			Body().Equal("foo")
	})

	t.Run("should parse request body GET /foo/context", func(t *testing.T) {
		testApp.Get("/foo/context").
			Expect().Status(http.StatusOK)
	})

	t.Run("should parse request body GET /foo/err", func(t *testing.T) {
		testApp.Get("/foo/err").
			Expect().Status(http.StatusInternalServerError)
	})

	//t.Run("should parse request body GET /foo/requestForm", func(t *testing.T) {
	//	testApp.Get("/foo/requestForm").
	//		WithFormField("name", "foo").
	//		Expect().Status(http.StatusOK).
	//		Body().Equal("foo")
	//})

	t.Run("should parse request body GET /foo/requestParams", func(t *testing.T) {
		testApp.Get("/foo/requestParams").
			WithQuery("name", "foo").
			Expect().Status(http.StatusOK).
			Body().Equal("foo")
	})

	t.Run("should return http.StatusUnauthorized after GET /bar", func(t *testing.T) {
		testApp.Get("/bar").
			Expect().Status(http.StatusUnauthorized)
	})

	t.Run("should return http.StatusInternalServerError when input form field validation failed", func(t *testing.T) {
		// test request form
		testApp.Post("/foobar").
			WithFormField("name", "John Doe").
			Expect().Status(http.StatusInternalServerError)
	})

	t.Run("should return (http.StatusOK on /foobar", func(t *testing.T) {
		//  test request query
		testApp.Get("/foobar").
			WithQuery("name", "John Doe").
			Expect().Status(http.StatusOK)
	})

	t.Run("should return http.StatusOK on /bar with jwt token", func(t *testing.T) {
		log.Println(io.GetWorkDir())
		jwtToken := jwt.NewJwtToken(&jwt.Properties{
			PrivateKeyPath: "config/ssl/app.rsa",
			PublicKeyPath:  "config/ssl/app.rsa.pub",
		})
		// test jwt
		token, err := jwtToken.Generate(jwt.Map{
			"username": "johndoe",
			"password": "PA$$W0RD",
		}, 100, time.Millisecond)
		if err == nil {

			t := fmt.Sprintf("Bearer %v", token)

			testApp.Get("/bar").
				WithHeader("Authorization", t).
				Expect().Status(http.StatusOK)

			time.Sleep(2 * time.Second)

			testApp.Get("/bar").
				WithHeader("Authorization", t).
				Expect().Status(http.StatusUnauthorized)
		}
	})

	t.Run("should return http.StatusOK on /foo with PUT, PATCH, DELETE methods", func(t *testing.T) {
		testApp.Put("/foo/id/{id}/name/{name}/age/{age}").
			WithPath("id", 123456).
			WithPath("name", "Mike").
			WithPath("age", 18).
			Expect().Status(http.StatusOK)
	})

	t.Run("should return http.StatusOK on /foo with PUT, PATCH, DELETE methods", func(t *testing.T) {
		testApp.Put("/foo/id/{id}/name/{name}/age/{age}").
			WithPath("id", 0).
			WithPath("name", "").
			WithPath("age", 0).
			Expect().Status(http.StatusOK)
	})

	t.Run("should return http.StatusOK on /foo with PUT, PATCH, DELETE methods", func(t *testing.T) {
		testApp.Put("/foo/id/{id}/name/{name}/age/{age}").
			WithPath("id", " ").
			WithPath("name", " ").
			WithPath("age", " ").
			Expect().Status(http.StatusOK)
	})

	t.Run("should return http.StatusOK on /foo with PUT, PATCH, DELETE methods", func(t *testing.T) {
		testApp.Put("/foo/id/{id}/name/{name}/age/{age}").
			WithPath("id", "").
			WithPath("name", "").
			WithPath("age", " ").
			Expect().Status(http.StatusOK)
	})

	t.Run("should return http.StatusOK on /foo with PUT, PATCH, DELETE methods", func(t *testing.T) {
		testApp.Put("/foo/id/{id}/name/{name}/age/{age}").
			WithPath("id", "").
			WithPath("name", "").
			WithPath("age", "").
			Expect().Status(http.StatusTemporaryRedirect)
	})

	t.Run("should Get foo by id", func(t *testing.T) {
		testApp.Get("/foo/id/{id}").
			WithPath("id", 123).
			Expect().Status(http.StatusOK)
	})

	t.Run("should Get by options", func(t *testing.T) {
		testApp.Get("/foo/options/{options}").
			WithPath("options", "mars,earth,mercury,jupiter").
			Expect().Status(http.StatusOK).
			Body().Contains("mercury")
	})

	t.Run("should Get foo by name", func(t *testing.T) {
		testApp.Get("/foo/name/{name}").
			WithPath("name", "Peter Phil").
			Expect().Status(http.StatusOK).
			Body().Contains("Peter Phil")
	})

	t.Run("should Get foo by name", func(t *testing.T) {
		testApp.Get("/foo/name/{name}").
			WithPath("name", "张三").
			Expect().Status(http.StatusOK).
			Body().Contains("张三")
	})

	t.Run("should Patch foo by id", func(t *testing.T) {
		testApp.Patch("/foo/id/{id}").
			WithPath("id", 456).
			Expect().Status(http.StatusOK)
	})

	t.Run("should Delete foo by id", func(t *testing.T) {
		testApp.Delete("/foo/id/{id}").
			WithPath("id", 789).
			Expect().Status(http.StatusOK)
	})

	t.Run("should call Options method", func(t *testing.T) {
		testApp.Options("/foo").
			Expect().Status(http.StatusOK)
	})

	t.Run("should return failure for unsupported args", func(t *testing.T) {
		testApp.Get("/foo/bar").
			Expect().Status(http.StatusInternalServerError)
	})

	t.Run("should return invalid response", func(t *testing.T) {
		testApp.Get("/foo/invalid").
			Expect().Status(http.StatusInternalServerError)
	})

	t.Run("should return integer", func(t *testing.T) {
		testApp.Get("/foo/integer").
			Expect().Status(http.StatusInternalServerError)
	})

	t.Run("should return integer pointer", func(t *testing.T) {
		testApp.Get("/foo/intPointer").
			Expect().Status(http.StatusInternalServerError)
	})

	t.Run("should return integer nil pointer", func(t *testing.T) {
		testApp.Get("/foo/intNilPointer").
			Expect().Status(http.StatusInternalServerError)
	})

	t.Run("should return error message", func(t *testing.T) {
		testApp.Get("/foo/error").
			Expect().Status(http.StatusInternalServerError)
	})
}

//func TestInvalidController(t *testing.T)  {
//	ta := new(testApplication)
//	err := ta.Init(new(InvalidController))
//	err, ok := err.(*system.InvalidControllerError)
//	assert.Equal(t, ok, true)
//}

func TestNewApplication(t *testing.T) {

	app.Register(new(ExampleController))

	testApp := web.NewApplication().
		SetAddCommandLineProperties(true).
		SetProperty("server.port", 8080).
		SetProperty(app.BannerDisabled, true)
	t.Run("should init web application", func(t *testing.T) {
		assert.NotEqual(t, nil, testApp)
	})

	t.Run("should get interface name", func(t *testing.T) {
		type myInterface interface {
		}

		mi := new(myInterface)

		typ, err := reflector.GetType(mi)
		assert.Equal(t, nil, err)
		assert.Equal(t, "myInterface", typ.Name())
	})

	go testApp.Run()
	time.Sleep(time.Second)
}

func TestAnonymousController(t *testing.T) {
	t.Run("should failed to register anonymous controller", func(t *testing.T) {
		testApp := web.RunTestApplication(t, (*Bar)(nil))
		assert.NotEqual(t, nil, testApp)
	})

	t.Run("should failed to register anonymous controller", func(t *testing.T) {
		testApp := web.RunTestApplication(t, newBar)
		assert.NotEqual(t, nil, testApp)
	})
}
