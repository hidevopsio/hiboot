package service

import (
	"fmt"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/starter/websocket"
	"time"
)

// CountHandler is the websocket handler
type CountHandler struct {
	at.ContextAware
	connection *websocket.Connection
}

func newCountHandler(connection *websocket.Connection) *CountHandler {
	h := &CountHandler{connection: connection}
	return h
}

func init() {
	app.Register(newCountHandler)
}

// OnMessage is the websocket message handler
func (h *CountHandler) OnMessage(data []byte) {
	message := string(data)
	log.Debugf("client: %v", message)
	var i int
	go func() {
		for {
			i++
			h.connection.EmitMessage([]byte(fmt.Sprintf("=== %v %d ===", message, i)))
			time.Sleep(time.Second)
		}
	}()
}

// OnDisconnect is the websocket disconnection handler
func (h *CountHandler) OnDisconnect() {
	log.Debugf("Connection with ID: %v has been disconnected!", h.connection.ID())
}
