package grpc

import (
	"github.com/hidevopsio/hiboot/pkg/inject"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/app"
)

type postProcessor struct {

}

func init() {
	// register postProcessor
	app.RegisterPostProcessor(new(postProcessor))
}

func (p *postProcessor) BeforeInitialization()  {
	log.Debug("[grpc] BeforeInitialization")
}

func (p *postProcessor) AfterInitialization()  {
	log.Debug("[grpc] AfterInitialization")
	for _, srv := range grpcServers {
		err := inject.IntoObject(srv.svc)
		if err != nil {
			break
		}
	}
}