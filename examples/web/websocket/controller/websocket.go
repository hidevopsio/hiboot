package controller

import (
	"github.com/hidevopsio/hiboot/examples/web/websocket/service"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/starter/websocket"
)

type websocketController struct {
	at.RestController
	connectionFunc           websocket.ConnectionFunc
	countHandlerConstructor  service.CountHandlerConstructor
	statusHandlerConstructor service.StatusHandlerConstructor
}

func newWebsocketController(connectionFunc websocket.ConnectionFunc,
	countHandlerConstructor service.CountHandlerConstructor,
	statusHandlerConstructor service.StatusHandlerConstructor) *websocketController {
	c := &websocketController{
		connectionFunc:           connectionFunc,
		countHandlerConstructor:  countHandlerConstructor,
		statusHandlerConstructor: statusHandlerConstructor,
	}
	return c
}

func init() {
	app.Register(newWebsocketController)
}

func (c *websocketController) Get(ctx *web.Context) {
	c.connectionFunc(ctx, c.countHandlerConstructor.(func(websocket.Connection) websocket.Handler))
}

func (c *websocketController) GetStatus(ctx *web.Context) {
	c.connectionFunc(ctx, c.statusHandlerConstructor.(func(websocket.Connection) websocket.Handler))
}
