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
	"reflect"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	pb "google.golang.org/grpc/health/grpc_health_v1"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/factory"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/starter/jaeger"
	"hidevops.io/hiboot/pkg/utils/cmap"
	"hidevops.io/hiboot/pkg/utils/reflector"
)

const (
	// Profile is the profile of grpc, it should be as same as the package name
	Profile = "grpc"
)

type configuration struct {
	app.Configuration
	Properties *properties

	instantiateFactory factory.InstantiateFactory
}

func newConfiguration(properties *properties, instantiateFactory factory.InstantiateFactory) *configuration {
	c := &configuration{Properties: properties, instantiateFactory: instantiateFactory}
	var dep []string
	for _, srv := range grpcServers {
		if srv.svc != nil {
			dep = append(dep, srv.name)
		}
	}
	c.RuntimeDeps.Set(c.ServerFactory, dep)
	return c
}

type grpcService struct {
	name string
	cb   interface{}
	svc  interface{}
}

var (
	grpcServers []*grpcService
	grpcClients []*grpcService

	clientMap cmap.ConcurrentMap

	registerHealthCheckService = false
)

// RegisterServer register server from application
func RegisterServer(register interface{}, server interface{}) {
	svrName := reflector.GetLowerCamelFullName(server)
	svr := &grpcService{
		name: svrName,
		cb:   register,
		svc:  server,
	}
	app.Register(server)
	grpcServers = append(grpcServers, svr)
}

// Server alias to RegisterServer
var Server = RegisterServer

// registerClient register client from application
func registerClient(name string, clientConstructors ...interface{}) {
	for _, clientConstructor := range clientConstructors {
		svr := &grpcService{
			name: name,
			cb:   clientConstructor,
		}
		grpcClients = append(grpcClients, svr)

		// pre-allocate client in order to pass dependency check
		typ, ok := reflector.GetObjectType(clientConstructor)
		if ok {
			// NOTE: it's very important !!!
			// To register grpc client and grpc.ClientConn in advance.
			// client should depends on grpc.clientFactory
			metaData := &factory.MetaData{
				MetaObject: reflect.New(typ).Interface(),
				DepNames:   []string{"grpc.clientFactory"},
			}
			app.Register(metaData)
		}
	}
	// Just register grpc.ClientConn in order to pass the dependency check
	app.Register(new(grpc.ClientConn))
}

// RegisterClient register client from application
func RegisterClient(name string, clientConstructors ...interface{}) {
	// register newHealthCheckService if grpc client is enabled
	if !registerHealthCheckService {
		registerHealthCheckService = true
		app.Register(NewHealthCheckService)
	}

	_, ok := clientMap.Get(name)
	if !ok {
		clientMap.Set(name, true)
		clientConstructors = append(clientConstructors, pb.NewHealthClient)
	}
	registerClient(name, clientConstructors...)
}

// Client register client from application, it is a alias to RegisterClient
var Client = RegisterClient

func init() {
	clientMap = cmap.New()
	Server(pb.RegisterHealthServer, health.NewServer)
	app.IncludeProfiles(jaeger.Profile)
	app.Register(newConfiguration, new(properties))
}


// ClientConnector is the interface that connect to grpc client
// it can be injected to struct at runtime
func (c *configuration) ClientConnector(tracer jaeger.Tracer) ClientConnector {
	return newClientConnector(c.instantiateFactory, tracer)
}

// GrpcClientFactory create gRPC Clients that registered by application
func (c *configuration) ClientFactory(cc ClientConnector) ClientFactory {
	return newClientFactory(c.instantiateFactory, c.Properties, cc)
}

// GrpcServer create new gRpc Server
func (c *configuration) Server(tracer jaeger.Tracer) (grpcServer *grpc.Server) {
	log.Debugf(">>> create grpc server")
	// just return if grpc server is not enabled
	if c.Properties.Server.Enabled {
		grpcServer = grpc.NewServer(
			grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
				// add opentracing stream interceptor to chain
				grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(tracer)),
			)),
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
				// add opentracing unary interceptor to chain
				grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(tracer)),
			)),
		)
	}
	return
}

// GrpcServerFactory create gRPC servers that registered by application
// go:depends
func (c *configuration) ServerFactory(grpcServer *grpc.Server) ServerFactory {
	return newServerFactory(c.instantiateFactory, c.Properties, grpcServer)
}
