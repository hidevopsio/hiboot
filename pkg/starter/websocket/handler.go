package websocket

import "github.com/kataras/iris/websocket"

type Handler interface {
	OnMessage(data []byte)
	OnDisconnect()
}

func HandleConnection(constructor HandlerConstructor, conn websocket.Connection) {
	handler := constructor(conn)
	conn.OnMessage(handler.OnMessage)
	conn.OnDisconnect(handler.OnDisconnect)
	conn.Wait()
}
