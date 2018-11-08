package service

import (
	"fmt"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/starter/websocket"
	"time"
)

type countHandler struct {
	connection websocket.Connection
}

type CountHandlerConstructor interface{}

func NewCountHandlerConstructor() CountHandlerConstructor {
	return func(connection websocket.Connection) websocket.Handler {
		return &countHandler{connection: connection}
	}
}

func init() {
	app.Register(NewCountHandlerConstructor)
}

func (h *countHandler) OnMessage(data []byte) {
	message := string(data)
	log.Debugf("client: %v", message)
	var i int
	for {
		i++
		h.connection.EmitMessage([]byte(fmt.Sprintf("=== %d ===", i)))
		time.Sleep(time.Second)
	}
}

func (h *countHandler) OnDisconnect() {
	log.Debugf("Connection with ID: %v has been disconnected!", h.connection.ID())
}
