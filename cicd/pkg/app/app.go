package app

import (
	"github.com/hidevopsio/hi/boot/pkg/application"
	"github.com/hidevopsio/hi/cicd/pkg/web"
)


func init() {

	app := application.Instance()

	deploymentController := web.Controller{}

	deploymentRoutes := app.Party("/cicd", deploymentController.Before)
	{
		// Method POST: http://localhost:8080/deployment/deploy
		deploymentRoutes.Post("/deploy", deploymentController.Deploy)
	}
}

func Run()  {
	application.Run()
}


