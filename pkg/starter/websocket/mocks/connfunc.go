package mocks

import (
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/starter/websocket"
	ws "github.com/kataras/iris/websocket"
)

func ConnectionFunc(ctx *web.Context, constructor websocket.HandlerConstructor) ws.Connection {
	conn := new(Connection)
	//handler := constructor(conn)
	//conn.OnMessage(handler.OnMessage)
	//conn.OnDisconnect(handler.OnDisconnect)
	//conn.Wait()
	return conn
}
