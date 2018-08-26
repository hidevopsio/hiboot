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
	"net"
	"google.golang.org/grpc"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
	"google.golang.org/grpc/reflection"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"fmt"
	"github.com/hidevopsio/hiboot/pkg/utils/mapstruct"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/factory/instance"
)

type grpcConfiguration struct {
	app.Configuration
	Properties properties `mapstructure:"grpc"`

	instanceFactory *instance.InstanceFactory
}

type grpcService struct {
	name string
	cb  interface{}
	svc interface{}
}

var (
	grpcServers        []*grpcService
	grpcClients		   []*grpcService
)

func RegisterServer(cb interface{}, s interface{})  {
	svr := &grpcService{
		cb:  cb,
		svc: s,
	}
	grpcServers = append(grpcServers, svr)
}

func RegisterClient(name string, cb interface{}, s ...interface{})  {
	var svc interface{}
	if s != nil && len(s) != 0 {
		svc = s[0]
	}
	svr := &grpcService{
		name: name,
		cb:  cb,
		svc: svc,
	}
	grpcClients = append(grpcClients, svr)
}

func init() {
	app.AutoConfiguration(new(grpcConfiguration))
}

// inject instanceFactory
func (c *grpcConfiguration) Init(instanceFactory *instance.InstanceFactory) {
	c.instanceFactory = instanceFactory
}

func (c *grpcConfiguration) BuildGrpcClients() {
	clientProps := c.Properties.Client
	for _, cli := range grpcClients {
		prop := new(client)
		if err := mapstruct.Decode(prop, clientProps[cli.name]); err != nil {
			break
		}
		host := prop.Host
		if host == "" {
			host = cli.name
		}
		address := host + ":" + prop.Port
		// connect to grpc server
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			break
		}
		log.Infof("grpc client connected to: %v", address)
		clientInstanceName := str.ToLowerCamel(cli.name)
		if cli.cb != nil {
			gRpcCli, err := reflector.CallFunc(cli.cb, conn)
			if err == nil {
				// register grpc client
				c.instanceFactory.SetInstance(clientInstanceName, gRpcCli)
			}
		}
		// register clientConn
		c.instanceFactory.SetInstance(clientInstanceName + "Conn", conn)

		// register client service
		if cli.svc != nil {
			svcName, err := reflector.GetName(cli.svc)
			if err == nil {
				c.instanceFactory.SetInstance(svcName, cli.svc)
			}
		}
	}
}

func (c *grpcConfiguration) RunGrpcServers() {
	// just return if grpc server is not enabled
	if !c.Properties.Server.Enabled {
		return
	}

	address := c.Properties.Server.Host + ":" + c.Properties.Server.Port
	lis, err := net.Listen(c.Properties.Server.Network, address)
	if err != nil {
		log.Fatalf("failed to listen: %v, %v", address, err)
	}
	grpcServer := grpc.NewServer()

	// register server
	// Register reflection service on gRPC server.
	for _, srv := range grpcServers {
		reflector.CallFunc(srv.cb, grpcServer, srv.svc)
		reflection.Register(grpcServer)
		c := make(chan bool)
		go func() {
			c <- true
			if err := grpcServer.Serve(lis); err != nil {
				fmt.Errorf("failed to serve: %v", err)
			}
			fmt.Printf("grpc server exit\n")
		}()
		<- c
		log.Infof("grpc server listening at: %v", address)
	}
}