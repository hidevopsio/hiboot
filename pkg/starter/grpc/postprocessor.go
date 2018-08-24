package grpc

import (
	"reflect"
	"github.com/hidevopsio/hiboot/pkg/inject"
	"github.com/hidevopsio/hiboot/pkg/factory"
)

type postProcessor struct {

}

func init() {
// register postProcessor
	factory.AddPostProcessor(new(postProcessor))
}

func (p *postProcessor) BeforeInitialization()  {

}

func (p *postProcessor) AfterInitialization()  {
	for _, srv := range grpcServers {
		err := inject.IntoObject(reflect.ValueOf(srv.svc))
		if err != nil {
			break
		}
	}
}