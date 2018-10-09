package grpc

import (
	"github.com/hidevopsio/hiboot/pkg/utils/mapstruct"
)

// ClientFactory build grpc clients
type ClientFactory interface {
}

type clientFactory struct {
}

func newClientFactory(properties properties, cc ClientConnector) ClientFactory {
	cf := &clientFactory{}
	cf.buildClients(properties, cc)

	return cf
}

func (f *clientFactory) buildClients(properties properties, cc ClientConnector) {
	clientProps := properties.Client
	for _, cli := range grpcClients {
		prop := new(ClientProperties)
		err := mapstruct.Decode(prop, clientProps[cli.name])
		if err == nil {
			cc.Connect(cli.name, cli.cb, prop)
		}
	}
}
