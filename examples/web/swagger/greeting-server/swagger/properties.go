package swagger

import (
	"github.com/go-openapi/spec"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
)

const (
	Profile = "swagger"
)

type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

type License struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type OpenAPIDefinition struct {
	at.ConfigurationProperties `value:"swagger"`
	spec.Swagger
}

func init() {
	app.Register(new(OpenAPIDefinition))
}