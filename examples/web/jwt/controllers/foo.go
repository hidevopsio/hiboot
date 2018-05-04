package controllers

import (
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/starter/web/jwt"
	"time"
	"net/http"
)

type FooRequest struct {
	Name string
}

type FooResponse struct {
	Greeting string
}

type FooController struct{
	web.Controller
}

// init - add &FooController{} to web application
func init()  {
	web.Add(&FooController{})
}

func (c *FooController) Before(ctx *web.Context)  {
	log.Debug("FooController.Before")
	ctx.Next()
}

// Post login
// The first word of method is the http method POST, the rest is the context mapping
func (c *FooController) PostLogin(ctx *web.Context)  {
	log.Debug("FooController.Login")

	userRequest := &UserRequest{}
	if ctx.RequestBody(userRequest) == nil {
		jwtToken, err := jwt.GenerateToken(jwt.Map{
			"username": userRequest.Username,
			"password": userRequest.Password,
		}, 10, time.Minute)

		//log.Debugf("token: %v", *jwtToken)

		if err == nil {
			ctx.Response("Success", jwtToken)
		} else {
			ctx.ResponseError(err.Error(), http.StatusInternalServerError)
		}
	}
}

func (c *FooController) PostSayHello(ctx *web.Context)  {
	log.Debug("FooController.SayHello")

	foo := &FooRequest{}
	if ctx.RequestBody(foo) == nil {
		ctx.Response("Success", &FooResponse{Greeting: "hello, " + foo.Name})
	}
}