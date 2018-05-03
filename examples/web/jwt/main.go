package main

import (
	"time"
	"net/http"
	"github.com/hidevopsio/hiboot/pkg/starter/web/jwt"
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"github.com/hidevopsio/hiboot/pkg/log"
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
}

type BarController struct{
}

func (c *FooController) PostLogin(ctx *web.Context)  {
	log.Debug("FooController.SayHello")

	userRequest := &UserRequest{}
	if ctx.RequestBody(userRequest) == nil {
		jwtToken, err := ctx.GenerateJwtToken(jwt.Map{
			"username": userRequest.Username,
			"password": userRequest.Password,
		}, 10, time.Minute)

		log.Debugf("token: %v", jwtToken)

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

func (c *BarController) GetSayHello(ctx *web.Context)  {
	log.Debug("BarController.SayHello")

	ctx.Response("Success", &Bar{Greeting: "hello bar"})

}

type Controllers struct{

	Foo *FooController `auth:"anon"`
	Bar *BarController `controller:"bar" auth:"anon"`
}

func main()  {

	controllers := &Controllers{}
	app, err := web.NewApplication(controllers)
	if err == nil {
		app.Run()
	}
}