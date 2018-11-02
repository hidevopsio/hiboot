package grpc

import (
	"hidevops.io/hiboot/pkg/factory"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/utils/mapstruct"
	"hidevops.io/hiboot/pkg/utils/reflector"
)

// ClientFactory build grpc clients
type ClientFactory interface {
}

type clientFactory struct {
}

func newClientFactory(instantiateFactory factory.InstantiateFactory, properties properties, cc ClientConnector) ClientFactory {
	cf := &clientFactory{}

	clientProps := properties.Client
	var gRPCCli interface{}
	for _, cli := range grpcClients {
		prop := new(ClientProperties)
		err := mapstruct.Decode(prop, clientProps[cli.name])
		if err == nil {
			gRPCCli, err = cc.Connect(cli.name, cli.cb, prop)
			if err == nil {
				clientInstanceName := reflector.GetName(gRPCCli)
				// register grpc client
				instantiateFactory.SetInstance(clientInstanceName, gRPCCli)

				log.Infof("Registered gRPC client %v", clientInstanceName)
			}
		}
	}

	return cf
}
