package grpc

import (
	"github.com/hidevopsio/hiboot/pkg/inject"
	"github.com/hidevopsio/hiboot/pkg/app"
)

type postProcessor struct {

}

func init() {
	// register postProcessor
	app.RegisterPostProcessor(new(postProcessor))
}

func (p *postProcessor) BeforeInitialization(factory interface{})  {
	//log.Debug("[grpc] BeforeInitialization")
}

func (p *postProcessor) AfterInitialization(factory interface{})  {
	//log.Debug("[grpc] AfterInitialization")
	for _, srv := range grpcServers {
		inject.IntoObject(srv.svc)
	}

	for _, cli := range grpcClients {
		inject.IntoObject(cli.svc)
	}
}