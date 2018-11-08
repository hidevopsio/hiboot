package service

import (
	"fmt"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/starter/websocket"
)

type statusHandler struct {
	connection websocket.Connection
}

type StatusHandlerConstructor interface{}

func NewStatusHandlerConstructor() StatusHandlerConstructor {
	return func(connection websocket.Connection) websocket.Handler {
		return &statusHandler{connection: connection}
	}
}

func init() {
	app.Register(NewStatusHandlerConstructor)
}

func (h *statusHandler) OnMessage(data []byte) {
	message := string(data)
	log.Debugf("client: %v", message)

	h.connection.EmitMessage([]byte(fmt.Sprintf("Status: Up")))

}

func (h *statusHandler) OnDisconnect() {
	log.Debugf("Connection with ID: %v has been disconnected!", h.connection.ID())
}
