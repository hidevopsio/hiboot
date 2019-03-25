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

// Package jaeger provides the quick start jaeger application example
// main package
package main

// import web starter from hiboot
import (
	"context"
	"github.com/opentracing/opentracing-go/log"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/starter/jaeger"
	"hidevops.io/hiboot/pkg/starter/logging"
	"hidevops.io/hiboot/pkg/utils/httpcli"
	"net/http"
)

// Controller Rest Controller with path /
// RESTful Controller, derived from at.RestController. The context mapping of this controller is '/' by default
type Controller struct {
	// at.RestController or at.RestController must be embedded here
	at.RestController

	Formatter string `value:"${provider.formatter}"`
	Publisher string `value:"${provider.publisher}"`
}

// Get GET /name/{name}
func (c *Controller) GetByName(name string, span *jaeger.Span) string {
	span.SetTag("hello-to", name)

	helloStr := c.formatString(span, name)
	c.printHello(span, helloStr)

	span.Finish()
	// response
	return helloStr
}


func (c *Controller) formatString(span *jaeger.Span, helloTo string) string {

	finalUrl := c.Formatter + "/formatter/" + helloTo
	req, err := http.NewRequest("GET", finalUrl, nil)
	if err != nil {
		panic(err.Error())
	}

	// call formatter service
	newSpan := span.Inject(context.Background(), "GET", finalUrl, "formatString")
	defer newSpan.Finish()

	resp, err := httpcli.Do(req)
	if err != nil {
		panic(err.Error())
	}

	helloStr := string(resp)

	newSpan.LogFields(
		log.String("event", "string-format"),
		log.String("value", helloStr),
	)

	return helloStr
}

func (c *Controller) printHello(span *jaeger.Span, helloStr string) {

	finalUrl := c.Publisher + "/publisher/" + helloStr
	req, err := http.NewRequest("GET", finalUrl, nil)
	if err != nil {
		panic(err.Error())
	}
	// call publisher service
	newSpan := span.Inject(context.Background(), "GET", finalUrl, "printHello")
	defer newSpan.Finish()

	if _, err := httpcli.Do(req); err != nil {
		panic(err.Error())
	}
}


func newController() *Controller {
	return &Controller{}
}

func init() {
	app.Register(newController)
}

// main function
func main() {
	// create new web application and run it
	web.NewApplication().
		SetProperty(app.ProfilesInclude, logging.Profile, web.Profile, jaeger.Profile).
		Run()
}
