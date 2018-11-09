package controller

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"hidevops.io/hiboot/examples/web/websocket/service"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/starter/jwt"
	"hidevops.io/hiboot/pkg/starter/websocket"
	"hidevops.io/hiboot/pkg/starter/websocket/mocks"
	"net/http"
	"testing"
	"time"
)

func TestWebSocketController(t *testing.T) {
	// create mock controller
	mockConn := new(mocks.Connection)
	countHandler := service.NewCountHandlerConstructor()
	statusHandler := service.NewStatusHandlerConstructor()
	handlerFunc := func(ctx *web.Context, constructor websocket.HandlerConstructor) websocket.Connection {
		handler := constructor(mockConn)
		mockConn.OnMessage(handler.OnMessage)
		mockConn.OnDisconnect(handler.OnDisconnect)
		mockConn.Wait()
		return mockConn
	}
	mockController := newWebsocketController(handlerFunc, countHandler, statusHandler)

	testApp := web.NewTestApp(mockController).Run(t)
	assert.NotEqual(t, nil, testApp)

	testWebsocket("/websocket", mockConn, testApp)
	testWebsocket("/websocket/status", mockConn, testApp)
}

func testWebsocket(path string,
	mockConn *mocks.Connection,
	testApp web.TestApplication) {

	mockConn.Mock.On("OnMessage", mock.AnythingOfType("websocket.NativeMessageFunc"))
	mockConn.Mock.On("OnDisconnect", mock.AnythingOfType("websocket.DisconnectFunc"))
	mockConn.Mock.On("Wait")

	jwtToken := jwt.NewJwtToken(&jwt.Properties{
		PrivateKeyPath: "config/ssl/app.rsa",
		PublicKeyPath:  "config/ssl/app.rsa.pub",
	})
	pt, _ := jwtToken.Generate(jwt.Map{
		"username": "johndoe",
		"password": "PA$$W0RD",
	}, 60, time.Second)
	token := fmt.Sprintf("Bearer %v", pt)
	testApp.Get(path).WithHeader("Authorization", token).Expect().Status(http.StatusOK)
}
