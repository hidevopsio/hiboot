package grpc

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/inject"
)

type postProcessor struct {
}

func init() {
	// register postProcessor
	app.RegisterPostProcessor(newPostProcessor)
}

func newPostProcessor() *postProcessor {
	return &postProcessor{}
}

func (p *postProcessor) AfterInitialization(factory interface{}) {
	//log.Debug("[grpc] AfterInitialization")
	// TODO should call factory.Register()
	for _, srv := range grpcServers {
		if srv.svc != nil {
			inject.IntoObject(srv.svc)
		}
	}

	for _, cli := range grpcClients {
		if cli.svc != nil {
			inject.IntoObject(cli.svc)
		}
	}
}
