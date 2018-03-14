package app

import (
	"github.com/hi-devops-io/hi/boot/pkg/application"
	"github.com/hi-devops-io/hi/cicd/pkg/controllers"
)


func init() {

	app := application.Instance()

	deploymentController := controllers.DeploymentController{}

	deploymentRoutes := app.Party("/deployment", deploymentController.Before)
	{
		// Method POST: http://localhost:8080/deployment/deploy
		deploymentRoutes.Post("/deploy", deploymentController.Deploy)
	}
}

func Run()  {
	application.Run()
}


