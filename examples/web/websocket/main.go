package main

import (
	_ "hidevops.io/hiboot/examples/web/websocket/controller"
	_ "hidevops.io/hiboot/examples/web/websocket/service"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/starter/websocket"
)

func main() {
	web.NewApplication().
		SetProperty(app.ProfilesInclude, websocket.Profile).
		Run()
}
