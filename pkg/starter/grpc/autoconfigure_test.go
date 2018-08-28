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
	"testing"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"golang.org/x/net/context"
)

func init() {
	RegisterServer(pb.RegisterGreeterServer, new(greeterService))
	RegisterClient("greeter-client", pb.NewGreeterClient)
}

// gRpc server
// server is used to implement helloworld.GreeterServer.
type greeterService struct{}

// SayHello implements helloworld.GreeterServer
func (s *greeterService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

// gRpc client
type greeterClient struct {
	greeterClient pb.GreeterClient
}

// Init inject greeterClient and clientContext
func (s *greeterClient) Init(greeterClient pb.GreeterClient,)  {
	s.greeterClient = greeterClient
}

func (s *greeterClient) SayHello(name string) (*pb.HelloReply, error) {
	response, err := s.greeterClient.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	return response, err
}

func TestGrpcServerAndClient(t *testing.T) {
	//grpcConfig := configuration{
	//	Properties: properties{
	//		TimeoutSecond: 1,
	//		Server: server{
	//			Enabled: true,
	//			Network: "tcp",
	//			Port: "7575",
	//		},
	//		Client: map[string]interface{}{
	//			"greeter-client": client{
	//				Host: "localhost",
	//				Port: "7575",
	//			},
	//		},
	//	},
	//}

	//factory := starter.GetFactory()
	//factory.Instantiate(&grpcConfig)

	//greeterSvc := new(greeterClient)
	//
	//name := "Steve"
	//response, err := greeterSvc.SayHello(name)
	//assert.Equal(t, nil, err)
	//assert.Equal(t, "Hello " + name, response.Message)
}

func TestInjectIntoObject(t *testing.T) {
	//InjectIntoObject()
}