package swagger

import (
	"github.com/go-openapi/spec"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
)

const (
	Profile = "swagger"
)

type OpenAPIDefinition struct {
	at.ConfigurationProperties `value:"swagger"`
	spec.Swagger
}

func init() {
	app.Register(new(OpenAPIDefinition))
}