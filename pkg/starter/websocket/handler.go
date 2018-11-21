package websocket

// Handler is the interface the websocket handler
type Handler interface {
	OnMessage(data []byte)
	OnDisconnect()
}

// Register is the handler register
type Register func(handler Handler, conn *Connection)

// registerHandler is the handler for websocket
func registerHandler(handler Handler, conn *Connection) {
	conn.OnMessage(handler.OnMessage)
	conn.OnDisconnect(handler.OnDisconnect)
	conn.Wait()
}
