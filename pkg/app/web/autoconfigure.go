package web

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/log"
)

const Profile = "web"

type configuration struct {
	at.AutoConfiguration

	Properties *properties
}

func newWebConfiguration() *configuration {
	return &configuration{}
}

func init() {
	app.IncludeProfiles(Profile)
	app.Register(newWebConfiguration)
}

// Context is the instance of context.Context
func (c *configuration) Context(app *webApp) context.Context {
	return NewContext(app)
}

// DefaultView set the default view
func (c *configuration) DefaultView(app *webApp) {

	if c.Properties.View.Enabled {
		v := c.Properties.View
		app.RegisterView(iris.HTML(v.ResourcePath, v.Extension))

		route := app.Get(v.ContextPath, Handler(func(ctx context.Context) {
			ctx.View(v.DefaultPage)
		}))
		route.MainHandlerName = fmt.Sprintf("%s%s ", v.ContextPath, v.DefaultPage)
		log.Infof("Mapped \"%v\" onto %v", v.ContextPath, v.DefaultPage)
	}
}
