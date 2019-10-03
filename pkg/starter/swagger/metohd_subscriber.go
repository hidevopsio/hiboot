package swagger

import (
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/inject/annotation"
)

type HttpMethod interface {
	GetMethod() string
	GetPath() string
}

type httpMethodSubscriber struct {
	at.HttpMethodSubscriber `value:"swagger"`
	apiPathsBuilder         *apiPathsBuilder
}

func newHttpMethodSubscriber(builder *apiPathsBuilder) *httpMethodSubscriber {
	return &httpMethodSubscriber{apiPathsBuilder: builder}
}

// TODO: use data instead of atController
func (s *httpMethodSubscriber) Subscribe(atController *annotation.Annotations, atMethod *annotation.Annotations) {
	s.apiPathsBuilder.Build(atController, atMethod)
}
