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
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	_ "github.com/hidevopsio/hiboot/pkg/starter/actuator"
	"github.com/hidevopsio/hiboot/pkg/starter/grpc"
	"github.com/hidevopsio/hiboot/examples/grpc/helloworld/protobuf"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"time"
	"golang.org/x/net/context"
)

// controller
type greeterController struct {
	web.Controller
	greeterClient protobuf.GreeterClient
}

func (c *greeterController) Init(greeterClient protobuf.GreeterClient, clientContext grpc.ClientContext)  {
	c.greeterClient = greeterClient
}

func (c *greeterController) GetByName(name string) string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := c.greeterClient.SayHello(ctx, &protobuf.HelloRequest{Name: name})

	if err == nil {
		return response.Message
	}

	return "no response from grpc server"
}

func init() {
	// for running test
	io.EnsureWorkDir("examples/grpc/helloworld/greeter-client")
	// register grpc client
	grpc.RegisterClient("greeter-client", protobuf.NewGreeterClient)
	// register greeterController
	web.Add(new(greeterController))
}

func main() {
	web.NewApplication().Run()
}
