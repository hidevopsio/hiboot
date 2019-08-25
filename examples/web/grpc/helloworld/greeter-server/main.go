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
// find protoc install instruction on http://google.github.io/proto-lens/installing-protoc.html
// go get -u -v github.com/golang/protobuf/{proto,protoc-gen-go}
//go:generate protoc -I ../protobuf --go_out=plugins=grpc:../protobuf ../protobuf/helloworld.proto

package main

import (
	"golang.org/x/net/context"
	protobuf2 "hidevops.io/hiboot/examples/web/grpc/helloworld/protobuf"
	"hidevops.io/hiboot/pkg/app/web"
	_ "hidevops.io/hiboot/pkg/starter/actuator"
	"hidevops.io/hiboot/pkg/starter/grpc"
)

// server is used to implement protobuf.GreeterServer.
type helloServiceServerImpl struct {
}

func newHelloServiceServer() protobuf2.HelloServiceServer {
	return &helloServiceServerImpl{}
}

// SayHello implements helloworld.GreeterServer
func (s *helloServiceServerImpl) SayHello(ctx context.Context, request *protobuf2.HelloRequest) (*protobuf2.HelloReply, error) {
	// response to client
	return &protobuf2.HelloReply{Message: "Hello " + request.Name}, nil
}

// server is used to implement protobuf.GreeterServer.
type holaServiceServerImpl struct {
}

func newHolaServiceServer() protobuf2.HolaServiceServer {
	return &holaServiceServerImpl{}
}

// SayHello implements helloworld.GreeterServer
func (s *holaServiceServerImpl) SayHola(ctx context.Context, request *protobuf2.HolaRequest) (*protobuf2.HolaReply, error) {
	// response to client
	return &protobuf2.HolaReply{Message: "Hola " + request.Name}, nil
}

func init() {
	// must: register grpc server
	// please note that holaServiceServerImpl must implement protobuf.HelloServiceServer, or it won't be registered.
	grpc.Server(protobuf2.RegisterHelloServiceServer, newHelloServiceServer)
	grpc.Server(protobuf2.RegisterHolaServiceServer, newHolaServiceServer)
}

func main() {
	// create new web application and run it
	web.NewApplication().Run()
}
