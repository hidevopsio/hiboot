package websocket

import (
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/starter/websocket/ws"
)

// Connection is the websocket connection
type Connection struct {
	at.ContextAware
	websocket.Connection
}

func newConnection(conn websocket.Connection) *Connection {
	return &Connection{Connection: conn}
}
