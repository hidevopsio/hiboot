package controller

import (
	"hidevops.io/hiboot/examples/web/websocket/service"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/starter/websocket"
)

type websocketController struct {
	at.RestController

	register websocket.Register
}

func newWebsocketController(interactRegister websocket.Register) *websocketController {
	return &websocketController{register: interactRegister}
}

func init() {
	app.Register(newWebsocketController)
}

// Get GET /websocket
func (c *websocketController) Get(handler *service.CountHandler, connection *websocket.Connection) {
	c.register(handler, connection)
}

// GetStatus GET /websocket/status
func (c *websocketController) GetStatus(handler *service.StatusHandler, connection *websocket.Connection) {
	c.register(handler, connection)
}
