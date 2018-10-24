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

package controller

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/log"
)

type fooRequest struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age"`
}

type fooResponse struct {
	Greeting string `json:"greeting"`
	Age      int    `json:"age"`
}

type fooController struct {
	web.Controller
}

// init - add &FooController{} to web application
func init() {
	app.Register(newFooController)
}

func newFooController() *fooController {
	return &fooController{}
}

// Before intercept all requests that coming into this controller
func (c *fooController) Before(ctx *web.Context) {
	log.Debug("FooController.Before")
	ctx.Next()
}

// Post endpoint POST /foo
func (c *fooController) Post(ctx *web.Context) {
	log.Debug("FooController.Post")

	foo := &fooRequest{}
	if ctx.RequestBody(foo) == nil {
		ctx.ResponseBody("success", &fooResponse{Greeting: "Hello, " + foo.Name})
	}

}

// Get endpoint GET /foo
func (c *fooController) Get(ctx *web.Context) {
	log.Debug("FooController.Get")

	foo := &fooRequest{}

	if ctx.RequestParams(foo) == nil {
		ctx.ResponseBody("success", &fooResponse{
			Greeting: "Hello, " + foo.Name,
			Age:      foo.Age})
	}
}

// After interceptor of the controller
func (c *fooController) After(ctx *web.Context) {
	log.Debug("FooController.After")
}
