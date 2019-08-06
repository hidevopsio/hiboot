package websocket

// Handler is the interface the websocket handler
type Handler interface {
	OnMessage(data []byte)
	OnDisconnect()
	OnPing()
	OnPong()
}

// Register is the handler register
type Register func(handler Handler, conn *Connection)
type InteractRegister func(handler Handler, conn *Connection)

// registerHandler is the handler for websocket
func registerHandler(handler Handler, conn *Connection) {
	conn.OnMessage(handler.OnMessage)
	conn.OnPong(handler.OnPong)
	conn.OnPing(handler.OnPing)
	conn.OnDisconnect(handler.OnDisconnect)
}

// registerHandler is the handler for websocket
func registerInteractHandler(handler Handler, conn *Connection) {
	registerHandler(handler, conn)
	conn.Wait()
}
