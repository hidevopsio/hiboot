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
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/inject"
	"github.com/hidevopsio/hiboot/pkg/starter/grpc"
	hwmock "github.com/hidevopsio/hiboot/pkg/starter/grpc/mock_helloworld"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc/examples/helloworld/helloworld"
	"testing"
	"time"
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

// rpcMsg implements the gomock.Matcher interface
type rpcMsg struct {
	msg proto.Message
}

func (r *rpcMsg) Matches(msg interface{}) bool {
	m, ok := msg.(proto.Message)
	if !ok {
		return false
	}
	return proto.Equal(m, r.msg)
}

func (r *rpcMsg) String() string {
	return fmt.Sprintf("is %s", r.msg)
}

func TestGrpcServerAndClient(t *testing.T) {

	app.Component(newGreeterClientService)

	grpc.Server(helloworld.RegisterGreeterServer, newGreeterServerService)
	grpc.Client("greeter-service", helloworld.NewGreeterClient)

	testApp := web.NewTestApplication(t)
	assert.NotEqual(t, nil, testApp)

	applicationContext := testApp.(app.ApplicationContext)

	t.Run("should find gRpc client and call its services", func(t *testing.T) {
		cliSvc := applicationContext.GetInstance(greeterClientService{})
		assert.NotEqual(t, nil, cliSvc)
		if cliSvc != nil {
			greeterCliSvc := cliSvc.(*greeterClientService)
			assert.NotEqual(t, nil, greeterCliSvc.greeterClient)
		}
	})

	t.Run("should connect to gRpc service at runtime", func(t *testing.T) {
		cc := applicationContext.GetInstance("grpc.clientConnector").(grpc.ClientConnector)
		assert.NotEqual(t, nil, cc)
		prop := new(grpc.ClientProperties)
		inject.DefaultValue(prop)
		grpcCli, err := cc.Connect("", helloworld.NewGreeterClient, prop)
		assert.Equal(t, nil, err)
		assert.NotEqual(t, nil, grpcCli)
	})

	t.Run("should get message from mock gRpc server", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockGreeterClient := hwmock.NewMockGreeterClient(ctrl)
		req := &helloworld.HelloRequest{Name: "unit_test"}
		mockGreeterClient.EXPECT().SayHello(
			gomock.Any(),
			&rpcMsg{msg: req},
		).Return(&helloworld.HelloReply{Message: "Mocked Interface"}, nil)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := mockGreeterClient.SayHello(ctx, &helloworld.HelloRequest{Name: "unit_test"})
		assert.Equal(t, nil, err)
		assert.Equal(t, "Mocked Interface", r.Message)
	})
}
