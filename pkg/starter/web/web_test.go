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

type BarController struct{
	Controller
}

type InvalidController struct {}

func init() {
	utils.ChangeWorkDir("../../../")
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

// Get hello
// The first word of method is the http method GET, the rest is the context mapping hello
// in this method, the name Get means that the method context mapping is '/'
func (c *HelloController) Get(ctx *Context)  {

	ctx.ResponseBody("success", "Hello, World")
}

func TestHelloWorld(t *testing.T)  {

	// create new test server
	e, err := NewTestServer(t, &HelloController{Controller{ContextMapping: "/"}})
	assert.Equal(t, nil, err)

	// run the application
	e.Request("GET", "/").
		Expect().Status(http.StatusOK).Body().Contains("Success")
}

func TestHelloWorldLocale(t *testing.T)  {
	// create new test server
	e, err := NewTestServer(t, &HelloController{Controller{ContextMapping: "/"}})
	assert.Equal(t, nil, err)

	// test cn-ZH
	e.Request("GET", "/").
		WithHeader("Accept-Language", "cn-ZH").
	Expect().
		Status(http.StatusOK).
		Body().Contains("成功")

	// test en-US
	e.Request("GET", "/").
		WithHeader("Accept-Language", "en-US").
		Expect().
		Status(http.StatusOK).
		Body().Contains("Success")
}

func TestWebApplication(t *testing.T)  {

	e, err := NewTestServer(t,
		&FooController{},
		&BarController{Controller{AuthType: AuthTypeJwt}},
	)
	assert.Equal(t, nil, err)

	e.Request("POST", "/foo/login").WithJSON(&UserRequest{Username: "johndoe", Password: "iHop91#15"}).
		Expect().Status(http.StatusOK).Body().Contains("Success")

	e.Request("POST", "/foo/sayHello").WithJSON(&FooRequest{Name: "John"}).
		Expect().Status(http.StatusOK).Body().Contains("Success")

	e.Request("GET", "/bar/sayHello").
		Expect().Status(http.StatusUnauthorized)
}


func TestInvalidController(t *testing.T)  {

	_, err := NewTestServer(t, &InvalidController{})
	err, ok := err.(*system.InvalidControllerError)
	assert.Equal(t, ok, true)
}
