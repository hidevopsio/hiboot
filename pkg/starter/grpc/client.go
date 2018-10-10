package grpc

import (
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/mapstruct"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
)

// ClientFactory build grpc clients
type ClientFactory interface {
}

type clientFactory struct {
}

func newClientFactory(instantiateFactory factory.InstantiateFactory, properties properties, cc ClientConnector) ClientFactory {
	cf := &clientFactory{}

	clientProps := properties.Client
	var gRpcCli interface{}
	for _, cli := range grpcClients {
		prop := new(ClientProperties)
		err := mapstruct.Decode(prop, clientProps[cli.name])
		if err == nil {
			gRpcCli, err = cc.Connect(cli.name, cli.cb, prop)
			if err == nil {
				clientInstanceName, err := reflector.GetName(gRpcCli)
				if err == nil {
					// register grpc client
					instantiateFactory.SetInstance(clientInstanceName, gRpcCli)

					log.Infof("Registered gRPC client %v", clientInstanceName)
				}
			}
		}
	}

	return cf
}
