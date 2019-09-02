package controller

import (
	"github.com/stretchr/testify/assert"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/app/web/context"
	"hidevops.io/hiboot/pkg/starter/websocket"
	"net/http"
	"testing"
)

func TestWebSocketController(t *testing.T) {
	mockController := newWebsocketController(func(handler websocket.Handler, conn *websocket.Connection) {
		// For controller's unit testing, do nothing
		ctx := conn.GetValue("context").(context.Context)
		ctx.StatusCode(http.StatusOK)
	})

	testApp := web.NewTestApp(mockController).Run(t)
	assert.NotEqual(t, nil, testApp)

	testApp.Get("/websocket").Expect().Status(http.StatusOK)
	testApp.Get("/websocket/status").Expect().Status(http.StatusOK)
}
