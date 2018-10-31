package controller

import (
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestWebSocketController(t *testing.T) {
	testApp := web.NewTestApp(newWebsocketController).
		SetProperty("web.view.enabled", true).
		SetProperty("server.port", 12768).
		Run(t)
	assert.NotEqual(t, nil, testApp)

	testApp.Get("/websocket").Expect().Status(http.StatusOK)
}
