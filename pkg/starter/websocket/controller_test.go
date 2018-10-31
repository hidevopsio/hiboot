package websocket_test

import (
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/starter/websocket"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type websocketController struct {
	at.RestController

	server  *websocket.Server
	context *web.Context
}

func newWebsocketController(server *websocket.Server, context *web.Context) *websocketController {
	c := &websocketController{
		server:  server,
		context: context,
	}

	return c
}

func (c *websocketController) Get(ctx *web.Context) {
	ctx.View("index.html")
}

func (c *websocketController) OnChat(msg string) {
	log.Infof("%s sent: %s\n", c.context.RemoteAddr(), msg)
	// Write message back to the client message owner with:
	// Write message to all except this client with:
	//connection.To(websocket.Broadcast).Emit("chat", msg)
}

func TestWebSocketController(t *testing.T) {
	testApp := web.NewTestApp(newWebsocketController).
		SetProperty("web.view.enabled", true).
		SetProperty("server.port", 12768).
		Run(t)
	assert.NotEqual(t, nil, testApp)

	testApp.Get("/websocket").Expect().Status(http.StatusOK)
}
