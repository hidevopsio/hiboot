//go:generate statik -src=./dist

package main

import (
	_ "hidevops.io/hiboot/examples/web/swagger/greeting-server/controller"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/starter/actuator"
	"hidevops.io/hiboot/pkg/starter/logging"
	"hidevops.io/hiboot/pkg/starter/swagger"
)

func init() {
	app.Register(swagger.OpenAPIDefinitionBuilder().
		Version("1.0.0").
		Title("HiBoot Swagger Demo Application - Greeting Server").
		Description("Greeting Server is an application that demonstrate the usage of Swagger Annotations").
		Schemes("http", "https").
		Host("localhost:8080").
		BasePath("/api/v1/greeting-server"),
	)
}

//run http://localhost:8080/api/v1/greeting-server/swagger-ui to open swagger ui
func main() {
	web.NewApplication().
		SetProperty(app.ProfilesInclude, actuator.Profile, logging.Profile, swagger.Profile).
		SetProperty("server.context_path", "/api/v1/greeting-server").
		Run()
}
