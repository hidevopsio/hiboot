package swagger

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/log"
)

type httpMethodSubscriber struct {
	at.HttpMethodSubscriber `value:"swagger"`
}

func newHttpMethodSubscriber() *httpMethodSubscriber {
	return &httpMethodSubscriber{}
}

// TODO: use data instead of atController
func (s *httpMethodSubscriber) Subscribe(atController *web.Annotations, atMethod *web.Annotations) {
	log.Debug("==================================================")
	log.Debug(atController)
	log.Debug(atMethod)
	log.Debug("==================================================")
}

func init() {
	app.Register(newHttpMethodSubscriber)
}