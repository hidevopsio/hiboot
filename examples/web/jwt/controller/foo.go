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

package controllers

import (
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/log"
	"time"
)

type UserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type FooRequest struct {
	Name string `json:"name" validate:"required"`
	Age  int `json:"age"`
}

type FooResponse struct {
	Greeting string `json:"greeting"`
	Age  int `json:"age"`
}

type FooController struct {
	web.Controller
}

// init - add &FooController{} to web application
func init() {
	web.Add(&FooController{})
}

func (c *FooController) Before(ctx *web.Context) {
	log.Debug("FooController.Before")
	ctx.Next()
}

// Post login
// The first word of method is the http method POST, the rest is the context mapping
func (c *FooController) PostLogin(ctx *web.Context) {
	log.Debug("FooController.Login")

	userRequest := &UserRequest{}
	if ctx.RequestBody(userRequest) == nil {
		jwtToken, _ := web.GenerateJwtToken(web.JwtMap{
			"username": userRequest.Username,
			"password": userRequest.Password,
		}, 10, time.Minute)

		//log.Debugf("token: %v", *jwtToken)

		ctx.ResponseBody("success", jwtToken)
	}
}

func (c *FooController) Post(ctx *web.Context) {
	log.Debug("FooController.Post")

	foo := &FooRequest{}
	if ctx.RequestBody(foo) == nil {
		ctx.ResponseBody("success", &FooResponse{Greeting: "Hello, " + foo.Name})
	}

}

func (c *FooController) Get(ctx *web.Context) {
	log.Debug("FooController.Get")

	foo := &FooRequest{}

	if ctx.RequestParams(foo) == nil {
		ctx.ResponseBody("success", &FooResponse{
			Greeting: "Hello, " + foo.Name,
			Age: foo.Age })
	}
}

func (c *FooController) After(ctx *web.Context) {
	log.Debug("FooController.After")
}