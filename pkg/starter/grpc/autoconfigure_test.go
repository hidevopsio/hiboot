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
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/starter/grpc"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc/examples/helloworld/helloworld"
	"testing"
)

// gRpc server
// server is used to implement helloworld.GreeterServer.
type greeterServerService struct{}

// newGreeterServerService is the constructor of greeterServerService
func newGreeterServerService() *greeterServerService {
	return &greeterServerService{}
}

// SayHello implements helloworld.GreeterServer
func (s *greeterServerService) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "Hello " + in.Name}, nil
}

// gRpc client
type greeterClientService struct {
	greeterClient helloworld.GreeterClient
}

// newGreeterClientService inject greeterClient
func newGreeterClientService(greeterClient helloworld.GreeterClient) *greeterClientService {
	return &greeterClientService{
		greeterClient: greeterClient,
	}
}

// SayHello gRpc client service that call remote server service
func (s *greeterClientService) SayHello(name string) (*helloworld.HelloReply, error) {
	response, err := s.greeterClient.SayHello(context.Background(), &helloworld.HelloRequest{Name: name})
	return response, err
}

func TestGrpcServerAndClient(t *testing.T) {

	app.Component(newGreeterClientService)

	grpc.Server(helloworld.RegisterGreeterServer, newGreeterServerService)
	grpc.Client("greeter-service", helloworld.NewGreeterClient)

	testApp := web.NewTestApplication(t)
	assert.NotEqual(t, nil, testApp)

	applicationContext := testApp.(app.ApplicationContext)

	cliSvc := applicationContext.FindInstance(greeterClientService{})
	assert.NotEqual(t, nil, cliSvc)
	if cliSvc != nil {
		greeterCliSvc := cliSvc.(*greeterClientService)
		assert.NotEqual(t, nil, greeterCliSvc.greeterClient)

		//name := "Steve"
		//response, err := greeterCliSvc.SayHello(name)
		//assert.Equal(t, nil, err)
		//assert.Equal(t, "Hello "+name, response.Message)
	}
}
