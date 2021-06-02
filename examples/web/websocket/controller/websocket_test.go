package controller

import (
	"net/http"
	"testing"

	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/starter/logging"
	"github.com/hidevopsio/hiboot/pkg/starter/websocket"
	"github.com/stretchr/testify/assert"
)

func TestWebSocketController(t *testing.T) {
	mockController := newWebsocketController(func(handler websocket.Handler, conn *websocket.Connection) {
		// For controller's unit testing, do nothing
		ctx := conn.GetValue("context").(context.Context)
		ctx.StatusCode(http.StatusOK)
	})

	testApp := web.NewTestApp(mockController).SetProperty(app.ProfilesInclude, websocket.Profile, logging.Profile).Run(t)
	assert.NotEqual(t, nil, testApp)

	testApp.Get("/websocket").Expect().Status(http.StatusOK)
	testApp.Get("/websocket/status").Expect().Status(http.StatusOK)
}
