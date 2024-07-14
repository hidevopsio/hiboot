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

// Package helloworld provides the quick start web application example
// main package
package main

// import web starter from hiboot
import (
	barfoo "github.com/hidevopsio/hiboot/examples/web/depinj/bar/foo"
	"github.com/hidevopsio/hiboot/examples/web/depinj/foo"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/starter/logging"
)

// Controller Rest Controller with path /
// RESTful Controller, derived from at.RestController. The context mapping of this controller is '/' by default
type Controller struct {
	// at.RestController or at.RestController must be embedded here
	at.RestController
	at.RequestMapping `value:"/"`

	testSvc   *foo.TestService
	barFooSvc *barfoo.TestService
}

func init() {
	// register Controller
	app.Register(newController)
}

func newController(
	theFooSvc *foo.TestService,
	theBarFooSvc *barfoo.TestService,
) *Controller {
	return &Controller{
		testSvc:   theFooSvc,
		barFooSvc: theBarFooSvc,
	}
}

// Get /
func (c *Controller) Get() string {
	// response
	log.Infof("foo: %v, bar: %v, baz: %v", c.testSvc.Name, c.barFooSvc.Name)
	return "Hello " + c.testSvc.Name + " and " + c.barFooSvc.Name
}

// main function
func main() {
	// create new web application and run it
	web.NewApplication().
		SetProperty(app.ProfilesInclude, logging.Profile, web.Profile).
		Run()
}
