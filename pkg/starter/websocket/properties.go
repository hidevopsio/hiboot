package websocket

type properties struct {
	ReadBufferSize  int `default:"1024"`
	WriteBufferSize int `default:"1024"`
}
