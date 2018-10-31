package mocks

import (
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/starter/websocket"
)

func HandlerFunc(ctx *web.Context, constructor websocket.HandlerConstructor) websocket.Connection {
	conn := new(Connection)

	//websocket.HandleConnection(constructor, conn)

	return conn
}
