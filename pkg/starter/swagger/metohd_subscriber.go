package swagger

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/inject/annotation"
)

type HttpMethod interface {
	GetMethod() string
	GetPath() string
}

type httpMethodSubscriber struct {
	at.HttpMethodSubscriber `value:"swagger"`
	pathsBuilder            *pathsBuilder
}

func newHttpMethodSubscriber(openAPIDefinitionBuilder *pathsBuilder) *httpMethodSubscriber {
	return &httpMethodSubscriber{pathsBuilder: openAPIDefinitionBuilder}
}

// TODO: use data instead of atController
func (s *httpMethodSubscriber) Subscribe(atController *annotation.Annotations, atMethod *annotation.Annotations) {
	s.pathsBuilder.Build(atController, atMethod)
}

func init() {
	app.Register(newHttpMethodSubscriber)
}