package main

import (
	"github.com/gobuffalo/packr"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/app/web/context"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/starter/actuator"
	"hidevops.io/hiboot/pkg/starter/logging"
)

type staticController struct {
	at.RestController
	at.RequestMapping `value:"/static"`
}

func init() {
	app.Register(newStaticController)
}

func newStaticController() *staticController {
	return &staticController{}
}

// ui
func (c *staticController) UI(at struct{ at.GetMapping `value:"/ui"` }, ctx context.Context) {

	box := packr.NewBox("./static")
	ctx.StaticResource(box)

	return
}

// static resource annotation
func (c *staticController) SimpleUI(at struct {
	at.GetMapping     `value:"/simple/ui"`
	at.StaticResource `value:"./static"`
}) {
}

func main() {
	web.NewApplication(newStaticController).
		SetProperty(app.ProfilesInclude, actuator.Profile, logging.Profile).
		Run()
}
