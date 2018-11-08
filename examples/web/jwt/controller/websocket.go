package controller

import (
	"hidevops.io/hiboot/examples/web/jwt/service"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/starter/jwt"
	"hidevops.io/hiboot/pkg/starter/websocket"
)

type websocketController struct {
	jwt.Controller
	handlerFunc              websocket.HandlerFunc
	countHandlerConstructor  service.CountHandlerConstructor
	statusHandlerConstructor service.StatusHandlerConstructor
}

func newWebsocketController(handlerFunc websocket.HandlerFunc,
	countHandlerConstructor service.CountHandlerConstructor,
	statusHandlerConstructor service.StatusHandlerConstructor) *websocketController {
	c := &websocketController{
		handlerFunc:              handlerFunc,
		countHandlerConstructor:  countHandlerConstructor,
		statusHandlerConstructor: statusHandlerConstructor,
	}
	return c
}

func init() {
	app.Register(newWebsocketController)
}

func (c *websocketController) Get(ctx *web.Context) {
	c.handlerFunc(ctx, c.countHandlerConstructor.(func(websocket.Connection) websocket.Handler))
}

func (c *websocketController) GetStatus(ctx *web.Context) {
	c.handlerFunc(ctx, c.statusHandlerConstructor.(func(websocket.Connection) websocket.Handler))
}
