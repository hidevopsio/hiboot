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

package grpc_test

import (
	"testing"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"golang.org/x/net/context"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/starter/grpc"
)

func init() {

}

// gRpc server
// server is used to implement helloworld.GreeterServer.
type greeterService struct{}

// SayHello implements helloworld.GreeterServer
func (s *greeterService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

// gRpc client
type greeterClientService struct {
	greeterClient pb.GreeterClient
}

// Init inject greeterClient and clientContext
func (s *greeterClientService) Init(greeterClient pb.GreeterClient)  {
	s.greeterClient = greeterClient
}

func (s *greeterClientService) SayHello(name string) (*pb.HelloReply, error) {
	response, err := s.greeterClient.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	return response, err
}

func TestGrpcServerAndClient(t *testing.T) {

	greeterClientSvc := new(greeterClientService)
	grpc.RegisterServer(pb.RegisterGreeterServer, new(greeterService))
	grpc.RegisterClient("greeter-client", pb.NewGreeterClient, greeterClientSvc)

	app := web.NewTestApplication(t)
	assert.NotEqual(t, nil, app)

	//name := "Steve"
	//response, err := greeterClientSvc.SayHello(name)
	//assert.Equal(t, nil, err)
	//assert.Equal(t, "Hello " + name, response.Message)
}

func TestInjectIntoObject(t *testing.T) {
	//InjectIntoObject()
}