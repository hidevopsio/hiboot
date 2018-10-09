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

// Package grpc provides the hiboot starter for injectable grpc client and server dependency
package grpc

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"google.golang.org/grpc"
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
func RegisterClient(name string, cbs ...interface{}) {
	for _, cb := range cbs {
		svr := &grpcService{
			name: name,
			cb:   cb,
		}
		grpcClients = append(grpcClients, svr)
	}
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

// ClientConnector is the interface that connect to grpc client
// it can be injected to struct at runtime
func (c *configuration) GrpcClientConnector() ClientConnector {
	return newClientConnector(c.instantiateFactory)
}

// RunGrpcServers create gRPC Clients that registered by application
func (c *configuration) GrpcClientFactory(cc ClientConnector) ClientFactory {
	return newClientFactory(c.Properties, cc)
}

// GrpcServer create new gRpc Server
func (c *configuration) GrpcServer() (grpcServer *grpc.Server) {
	// just return if grpc server is not enabled
	if c.Properties.Server.Enabled {
		grpcServer = grpc.NewServer()
	}
	return
}

// RunGrpcServers create gRPC servers that registered by application
func (c *configuration) GrpcServerFactory(grpcServer *grpc.Server) ServerFactory {
	return newServerFactory(c.Properties, grpcServer)
}
