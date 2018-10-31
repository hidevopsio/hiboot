package main

import (
	_ "github.com/hidevopsio/hiboot/examples/web/websocket/controller"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/starter/logging"
	"github.com/hidevopsio/hiboot/pkg/starter/websocket"
)

func main() {
	web.NewApplication().
		SetProperty(app.PropertyAppProfilesInclude, websocket.Profile, logging.Profile).
		SetProperty("web.view.enabled", true).
		SetProperty("server.port", 8080).
		Run()
}
