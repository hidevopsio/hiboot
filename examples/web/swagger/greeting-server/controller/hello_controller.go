package controller

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
)

type helloController struct {
	at.RestController
	at.RequestMapping `value:"/hello"`
}

func init() {
	app.Register(newHelloController)
}

func newHelloController() *helloController {
	return &helloController{}
}

type HelloQueryParam struct {
	at.RequestParams
	at.ApiModel `value:""`
	Name        string `api:"defaults to World if not given"`
}

// Hello
func (c *helloController) Hello(at struct {
	at.GetMapping   `value:"/"`
	at.ApiOperation `operationId:"getGreeting" description:"This is the Greeting api for demo"`
	at.ApiParam     `value:"Path variable employee ID" required:"true"`
	Code200 struct {
		at.ApiResponse `code:"200" description:"returns a greeting"`
		at.ApiResponseSchema `code:"200" type:"string" description:"contains the actual greeting as plain text"`
	}
	Code404 struct {
		at.ApiResponse `code:"404" description:"greeter is not available"`
		at.ApiResponseSchema `code:"404" type:"string" description:"Report 'not found' error message"`
	}
}, request *HelloQueryParam) (response string) {
	if request.Name == "" {
		request.Name = "world"
	}
	response = "Hello, " + request.Name
	return
}
