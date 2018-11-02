package mocks

import (
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/starter/websocket"
)

func HandlerFunc(ctx *web.Context, constructor websocket.HandlerConstructor) websocket.Connection {
	conn := new(Connection)

	//websocket.HandleConnection(constructor, conn)

	return conn
}
