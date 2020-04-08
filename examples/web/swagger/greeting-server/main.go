//go:generate statik -src=./dist

package main

import (
	_ "hidevops.io/hiboot/examples/web/swagger/greeting-server/controller"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/app/web/server"
	"hidevops.io/hiboot/pkg/starter/actuator"
	"hidevops.io/hiboot/pkg/starter/logging"
	"hidevops.io/hiboot/pkg/starter/swagger"
)

const (
	version  = "1.0.3"
	basePath = "/api/v1/greeting-server"
)

func init() {
	app.Register(swagger.ApiInfoBuilder().
		ContactName("John Deng").
		ContactEmail("john.deng@outlook.com").
		ContactURL("https://hidevops.io").
		Title("HiBoot Swagger Demo Application - Greeting Server").
		Description("Greeting Server is an application that demonstrate the usage of Swagger Annotations"),
	// alternatively, you can set below properties by using SetProperty() in main, config/application.yml or program arguments to take advantage of HiBoot DI
	//Version(version).
	//Schemes("http", "https").
	//Host("localhost:8080").
	//BasePath(basePath),
	)
}

//run http://localhost:8080/api/v1/greeting-server/swagger-ui to open swagger ui
func main() {
	web.NewApplication().
		SetProperty(app.ProfilesInclude, actuator.Profile, logging.Profile, swagger.Profile).
		SetProperty(app.Version, version).
		SetProperty(server.Host, "localhost:8080").
		SetProperty(server.ContextPath, basePath).
		Run()
}
