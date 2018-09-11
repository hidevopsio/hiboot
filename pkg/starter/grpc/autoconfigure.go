// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package grpc

import (
	"fmt"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/mapstruct"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type configuration struct {
	app.Configuration
	Properties properties `mapstructure:"grpc"`

	instantiateFactory factory.InstantiateFactory
}

type grpcService struct {
	name string
	cb   interface{}
	svc  interface{}
}

var (
	grpcServers []*grpcService
	grpcClients []*grpcService
)

// RegisterClient register server from application
func RegisterServer(cb interface{}, s interface{}) {
	svrName, _ := reflector.GetName(s)
	svr := &grpcService{
		name: svrName,
		cb:   cb,
		svc:  s,
	}
	grpcServers = append(grpcServers, svr)
}

// Server alias to RegisterServer
var Server = RegisterServer

// RegisterClient register client from application
func RegisterClient(name string, cb interface{}, s ...interface{}) {
	var svc interface{}
	if s != nil && len(s) != 0 {
		svc = s[0]
	}
	svr := &grpcService{
		name: name,
		cb:   cb,
		svc:  svc,
	}
	grpcClients = append(grpcClients, svr)
}

// Client register client from application, it is a alias to RegisterClient
var Client = RegisterClient

func init() {
	app.AutoConfiguration(newConfiguration)
}

func newConfiguration(instantiateFactory factory.InstantiateFactory) *configuration {
	return &configuration{
		instantiateFactory: instantiateFactory,
	}
}

// RunGrpcServers create gRPC Clients that registered by application
func (c *configuration) BuildGrpcClients() {
	clientProps := c.Properties.Client
	for _, cli := range grpcClients {
		prop := new(client)
		if err := mapstruct.Decode(prop, clientProps[cli.name]); err != nil {
			log.Error(err)
			break
		}
		host := prop.Host
		if host == "" {
			host = cli.name
		}
		address := host + ":" + prop.Port
		conn := c.instantiateFactory.GetInstance(cli.name)
		var err error
		if conn == nil {
			// connect to grpc server
			conn, err = grpc.Dial(address, grpc.WithInsecure())
			c.instantiateFactory.SetInstance(cli.name, conn)
			if err != nil {
				log.Errorf("failed to connect to grpc server %v", address)
				break
			}
			log.Infof("gRPC client connected to: %v", address)
		}
		if cli.cb != nil {
			// get return type for register instance name
			gRpcCli, err := reflector.CallFunc(cli.cb, conn)
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

		// register client service
		if cli.svc != nil {
			svcName, err := reflector.GetName(cli.svc)
			if err == nil {
				c.instantiateFactory.SetInstance(svcName, cli.svc)
			}
		}

	}
}

func (c *configuration) GrpcServer() *grpc.Server {
	// just return if grpc server is not enabled
	if !c.Properties.Server.Enabled {
		return nil
	}
	return grpc.NewServer()
}

// RunGrpcServers create gRPC servers that registered by application
func (c *configuration) RunGrpcServers(grpcServer *grpc.Server) {
	// just return if grpc server is not enabled
	if !c.Properties.Server.Enabled || grpcServer == nil {
		return
	}

	address := c.Properties.Server.Host + ":" + c.Properties.Server.Port
	lis, err := net.Listen(c.Properties.Server.Network, address)
	if err != nil {
		log.Fatalf("failed to listen: %v, %v", address, err)
	}

	// register server
	// Register reflection service on gRPC server.
	chn := make(chan bool)
	go func() {
		for _, srv := range grpcServers {
			reflector.CallFunc(srv.cb, grpcServer, srv.svc)
			svcName, err := reflector.GetName(srv.svc)
			if err == nil {
				log.Infof("Registered %v on gRPC server", svcName)
			}
		}
		reflection.Register(grpcServer)
		chn <- true
		if err := grpcServer.Serve(lis); err != nil {
			fmt.Printf("failed to serve: %v", err)
		}
		fmt.Printf("gRPC server exit\n")
	}()
	<-chn

	log.Infof("gRPC server listening on: localhost%v", address)
}
