package controllers

import (
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"github.com/hidevopsio/hiboot/pkg/log"
)

type UserRequest struct {
	Username string
	Password string
}


type Bar struct {
	Greeting string
}


type BarController struct{
	web.Controller
}


func init()  {
	web.Add(&BarController{
		web.Controller{
			ContextMapping: "/bars",
			AuthType:       web.AuthTypeJwt,
		},
	})
}

func (c *BarController) GetSayHello(ctx *web.Context)  {
	log.Debug("BarController.SayHello")

	ctx.Response("Success", &Bar{Greeting: "hello bar"})

}
