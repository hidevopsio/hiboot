package grpc

import (
	"google.golang.org/grpc"
	"hidevops.io/hiboot/pkg/factory"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/utils/reflector"
)

// ClientConnector interface is response for creating grpc client connection
type ClientConnector interface {
	// Connect connect the gRPC client
	Connect(name string, cb interface{}, prop *ClientProperties) (gRPCCli interface{}, err error)
}

type clientConnector struct {
	instantiateFactory factory.InstantiateFactory
}

func newClientConnector(instantiateFactory factory.InstantiateFactory) ClientConnector {
	cc := &clientConnector{
		instantiateFactory: instantiateFactory,
	}
	return cc
}

// Connect connect to grpc server from client
// name: client name
// clientConstructor: client constructor
// properties: properties for configuring
func (c *clientConnector) Connect(name string, clientConstructor interface{}, properties *ClientProperties) (gRPCCli interface{}, err error) {
	host := properties.Host
	if host == "" {
		host = name
	}
	address := host + ":" + properties.Port
	conn := c.instantiateFactory.GetInstance(name)
	if conn == nil {
		// connect to grpc server
		conn, err = grpc.Dial(address, grpc.WithInsecure())
		c.instantiateFactory.SetInstance(name, conn)
		if err == nil {
			log.Infof("gRPC client connected to: %v", address)
		}
	}
	if err == nil && clientConstructor != nil {
		// get return type for register instance name
		gRPCCli, err = reflector.CallFunc(clientConstructor, conn)
	}
	return
}
