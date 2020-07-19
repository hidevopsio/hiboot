package main

import (
	_ "hidevops.io/hiboot/examples/web/middleware/controller"
	_ "hidevops.io/hiboot/examples/web/middleware/logging"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/starter/actuator"
	"hidevops.io/hiboot/pkg/starter/logging"
)

// HiBoot main function
func main() {
	// create new web application and run it
	web.NewApplication().
		SetProperty(app.ProfilesInclude, web.Profile, actuator.Profile, logging.Profile).
		Run()
}