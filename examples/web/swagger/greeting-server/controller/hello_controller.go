package controller

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
)

type controller struct {
	at.RestController
	at.RequestMapping `value:"/"`
}

func init() {
	app.Register(newHelloController)
}

func newHelloController() *controller {
	return &controller{}
}

type HelloQueryParam struct {
	at.RequestParams
	at.ApiModel `value:""`
	Name        string `api:"defaults to World if not given"`
}

// Hello
func (c *controller) Hello(at struct {
	at.GetMapping `value:"/hello"`
	at.Operation  `operationId:"getGreeting" description:"This is the Greeting api for demo"`
	at.Produces   `values:"text/plain"`
	ParamName     struct {
		at.Parameter `type:"string" name:"name" in:"query" description:"defaults to World if not given" `
	}
	Code200 struct {
		at.Response       `code:"200" description:"returns a greeting"`
		at.ResponseSchema `code:"200" type:"string" description:"contains the actual greeting as plain text"`
	}
	Code404 struct {
		at.Response       `code:"404" description:"greeter is not available"`
		at.ResponseSchema `code:"404" type:"string" description:"Report 'not found' error message"`
	}
	// Response MyResponse
}, param *HelloQueryParam) (response string) {
	if param.Name == "" {
		param.Name = "world"
	}
	response = "Hello, " + param.Name
	return
}
