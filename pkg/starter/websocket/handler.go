package websocket

type Handler interface {
	OnMessage(data []byte)
	OnDisconnect()
}

// RegisterHandler is the handler register
type Register func(handler Handler, conn *Connection)

// HandleConnection is the handler for websocket
func registerHandler(handler Handler, conn *Connection) {
	conn.OnMessage(handler.OnMessage)
	conn.OnDisconnect(handler.OnDisconnect)
	conn.Wait()
}
