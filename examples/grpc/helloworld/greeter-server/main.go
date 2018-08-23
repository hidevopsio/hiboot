/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

//go:generate protoc -I ../protobuf --go_out=plugins=grpc:../protobuf ../protobuf/helloworld.proto

package main

import (
	"golang.org/x/net/context"
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	_ "github.com/hidevopsio/hiboot/pkg/starter/actuator"
	"github.com/hidevopsio/hiboot/pkg/starter/grpc"
	"github.com/hidevopsio/hiboot/examples/grpc/helloworld/protobuf"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
)

// server is used to implement helloworld.GreeterServer.
type greeterService struct{}

// SayHello implements helloworld.GreeterServer
func (s *greeterService) SayHello(ctx context.Context, in *protobuf.HelloRequest) (*protobuf.HelloReply, error) {
	return &protobuf.HelloReply{Message: "Hello " + in.Name}, nil
}

func init() {
	io.EnsureWorkDir("examples/grpc/helloworld/greeter-server")

	grpc.RegisterServer(protobuf.RegisterGreeterServer, new(greeterService))
}

func main() {
	web.NewApplication().Run()
}
