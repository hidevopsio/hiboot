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

package main

import (
	"github.com/hidevopsio/hiboot/pkg/starter/web"
)

// Define our controller, start with the name Foo, the first word of the Camelcase FooController is the controller name
// the lower cased foo will be the context mapping of the controller
// context mapping can be overwritten by FooController.ContextMapping
// if the controller name is a single word Controller, then the context mapping will be '/'
type Controller struct{
	web.Controller
}

// Get hello
// the first word of method is the http method GET, the rest is the context mapping hello
// if the method name is a single word Get, the the context mapping will be '/'
func (c *Controller) Get(ctx *web.Context)  {

	ctx.ResponseBody("Success", "Hello, World")
}

func main()  {

	// create new web application
	app, err := web.NewApplication(&Controller{})

	// run the application
	if err == nil {
		app.Run()
	}
}