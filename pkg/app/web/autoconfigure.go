package web

import (
	"fmt"
	"github.com/kataras/iris"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web/context"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/log"
)

const Profile = "web"

type configuration struct {
	at.AutoConfiguration

	Properties properties `mapstructure:"web"`
}

func newWebConfiguration() *configuration {
	return &configuration{}
}

func init() {
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
