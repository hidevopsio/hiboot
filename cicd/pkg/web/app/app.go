package app

import (
	"github.com/hidevopsio/hi/boot/pkg/application"
	"github.com/hidevopsio/hi/cicd/pkg/web/controllers"
)


func init() {

	app := application.Instance()

	controller := controllers.CicdController{}

	cicdRouters := app.Party("/cicd", controller.Before)
	{
		// Method POST: http://localhost:8080/deployment/deploy
		cicdRouters.Post("/run", controller.Run)
	}
}

func Run()  {
	application.Run()
}


