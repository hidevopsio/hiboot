package controller

import (
	"github.com/hidevopsio/hiboot/examples/web/websocket/service"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/starter/websocket/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestWebSocketController(t *testing.T) {
	testController := newWebsocketController(mocks.ConnectionFunc,
		service.NewCountHandlerConstructor(),
		service.NewStatusHandlerConstructor())
	testApp := web.NewTestApp(testController).
		SetProperty("web.view.enabled", true).
		SetProperty("server.port", 12768).
		Run(t)
	assert.NotEqual(t, nil, testApp)

	testApp.Get("/").Expect().Status(http.StatusOK)
	testApp.Get("/websocket").Expect().Status(http.StatusOK)
	testApp.Get("/websocket/status").Expect().Status(http.StatusOK)

}
