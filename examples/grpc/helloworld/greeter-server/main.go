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

// if protoc report command not found error, should install proto and protc-gen-go
// go get -u -v github.com/golang/protobuf/{proto,protoc-gen-go}
//go:generate protoc -I ../protobuf --go_out=plugins=grpc:../protobuf ../protobuf/helloworld.proto

package main

import (
	"github.com/hidevopsio/hiboot/examples/grpc/helloworld/protobuf"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	_ "github.com/hidevopsio/hiboot/pkg/starter/actuator"
	"github.com/hidevopsio/hiboot/pkg/starter/grpc"
	"golang.org/x/net/context"
)

// server is used to implement protobuf.GreeterServer.
type greeterServerImpl struct {
}

// SayHello implements helloworld.GreeterServer
func (s *greeterServerImpl) SayHello(ctx context.Context, request *protobuf.HelloRequest) (*protobuf.HelloReply, error) {
	// response to client
	return &protobuf.HelloReply{Message: "Hello " + request.Name}, nil
}

func init() {
	// must: register grpc server
	// please note that greeterService must implement protobuf.GreeterServer, or it won't be registered.
	grpc.RegisterServer(protobuf.RegisterGreeterServer, new(greeterServerImpl))
}

func main() {
	// create new web application and run it
	web.NewApplication().Run()
}
