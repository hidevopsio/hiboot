package websocket

import (
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/starter/websocket/at"
)

// Controller is the base web controller that use WebSocket
type Controller struct {
	at.WebSocketRestController
	web.Controller
}
