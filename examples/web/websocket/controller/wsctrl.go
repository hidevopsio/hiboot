package controller

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/starter/websocket"
	"sync/atomic"
)

type websocketController struct {
	at.RestController

	server     *websocket.Server
	connection websocket.Connection

	visits uint64
}

func newWebsocketController(server *websocket.Server) *websocketController {
	c := &websocketController{
		server: server,
	}

	return c
}

func init() {
	app.Register(newWebsocketController)
}

func (c *websocketController) onLeave(roomName string) {
	// visits--
	newCount := atomic.AddUint64(&c.visits, ^uint64(0))
	log.Debugf("[onLeave] %d online visitors", newCount)
	// This will call the "visit" event on all clients, except the current one,
	// (it can't because it's left but for any case use this type of design)
	c.connection.To(websocket.Broadcast).Emit("visit", newCount)
}

func (c *websocketController) update() {
	// visits++
	newCount := atomic.AddUint64(&c.visits, 1)
	log.Debugf("[update] %d online visitors", newCount)

	// This will call the "visit" event on all clients, including the current
	// with the 'newCount' variable.
	//
	// There are many ways that u can do it and faster, for example u can just send a new visitor
	// and client can increment itself, but here we are just "showcasing" the websocket controller.
	c.connection.To(websocket.All).Emit("visit", newCount)
}

func (c *websocketController) Get(ctx *web.Context) {
	log.Debug("GET /websocket")
	c.connection = c.server.Upgrade(ctx)
	c.connection.OnLeave(func(roomName string) {
		newCount := atomic.AddUint64(&c.visits, ^uint64(0))
		log.Debugf("[onLeave] %d online visitors", newCount)
		// This will call the "visit" event on all clients, except the current one,
		// (it can't because it's left but for any case use this type of design)
		c.connection.To(websocket.Broadcast).Emit("visit", newCount)
	})
	c.connection.On("visit", c.update)

	// call it after all event callbacks registration.
	c.connection.Wait()
}
