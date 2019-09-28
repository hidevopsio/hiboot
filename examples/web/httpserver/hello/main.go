package hello

// import web starter from hiboot
import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/starter/actuator"
	"hidevops.io/hiboot/pkg/starter/logging"
)

// main function
func main() {
	web.NewApplication().
		SetProperty(app.ProfilesInclude, logging.Profile, actuator.Profile).
		SetProperty(web.ViewEnabled, true).
		Run()
}
