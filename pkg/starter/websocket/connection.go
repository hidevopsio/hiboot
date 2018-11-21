package websocket

import (
	"github.com/kataras/iris/websocket"
	"hidevops.io/hiboot/pkg/at"
)

// Connection is the websocket connection
type Connection struct {
	at.ContextAware
	websocket.Connection
}

func newConnection(conn websocket.Connection) *Connection {
	return &Connection{Connection: conn}
}
