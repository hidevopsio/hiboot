package service

import (
	"fmt"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/starter/websocket"
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
