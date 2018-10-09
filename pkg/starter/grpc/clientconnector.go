package grpc

import (
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"google.golang.org/grpc"
)

// ClientConnector interface is response for creating grpc client connection
type ClientConnector interface {
	Connect(name string, cb interface{}, prop *ClientProperties) (err error)
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
func (c *clientConnector) Connect(name string, cb interface{}, prop *ClientProperties) (err error) {
	host := prop.Host
	if host == "" {
		host = name
	}
	address := host + ":" + prop.Port
	conn := c.instantiateFactory.GetInstance(name)
	if conn == nil {
		// connect to grpc server
		conn, err = grpc.Dial(address, grpc.WithInsecure())
		c.instantiateFactory.SetInstance(name, conn)
		if err == nil {
			log.Infof("gRPC client connected to: %v", address)
		}
	}
	if err == nil && cb != nil {
		// get return type for register instance name
		gRpcCli, err := reflector.CallFunc(cb, conn)
		if err == nil {
			clientInstanceName, err := reflector.GetName(gRpcCli)
			if err == nil {
				// register grpc client
				c.instantiateFactory.SetInstance(clientInstanceName, gRpcCli)
				// register clientConn
				c.instantiateFactory.SetInstance(clientInstanceName+"Conn", conn)

				log.Infof("Registered gRPC client %v", clientInstanceName)
			}
		}
	}
	return nil
}
