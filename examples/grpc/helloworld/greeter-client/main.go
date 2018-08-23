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
	// embedded web.Controller
	web.Controller
	// declare greeterClient
	greeterClient protobuf.GreeterClient
}

// Init inject greeterClient
func (c *greeterController) Init(greeterClient protobuf.GreeterClient)  {
	c.greeterClient = greeterClient
}

// GET /greeter/{name}
func (c *greeterController) GetByName(name string) string {

	// set 2 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
	defer cancel()

	// call grpc server method
	response, err := c.greeterClient.SayHello(ctx, &protobuf.HelloRequest{Name: name})

	// got response
	if err == nil {
		return response.Message
	}

	// response with err
	return err.Error()
}

func init() {
	// optional: for running test
	io.EnsureWorkDir("examples/grpc/helloworld/greeter-client")

	// must: register grpc client
	grpc.RegisterClient("greeter-client", protobuf.NewGreeterClient)

	// must: register greeterController
	web.Add(new(greeterController))
}

func main() {
	web.NewApplication().Run()
}
