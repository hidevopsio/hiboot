package websocket

type Handler interface {
	OnMessage(data []byte)
	OnDisconnect()
}
