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
	"github.com/hidevopsio/hiboot/examples/grpc/helloworld/protobuf"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	_ "github.com/hidevopsio/hiboot/pkg/starter/actuator"
	"github.com/hidevopsio/hiboot/pkg/starter/grpc"
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
func (c *greeterController) Init(greeterClient protobuf.GreeterClient) {
	c.greeterClient = greeterClient
}

// GET /greeter/name/{name}
func (c *greeterController) GetByName(name string) string {

	// call grpc server method
	// pass context.Background() for the sake of simplicity
	response, err := c.greeterClient.SayHello(context.Background(), &protobuf.HelloRequest{Name: name})

	// got response
	if err == nil {
		return response.Message
	}

	// response with err
	return err.Error()
}

func init() {

	// must: register grpc client, the name greeter-client should configured in application.yml
	// see config/application-grpc.yml
	//
	// grpc:
	//   client:
	// 	   greeter-client:   # client name
	//       host: localhost # server host
	//       port: 7575      # server port
	//
	grpc.RegisterClient("greeter-client", protobuf.NewGreeterClient)

	// must: register greeterController
	web.RestController(new(greeterController))
}

func main() {
	// create new web application and run it
	web.NewApplication().Run()
}
