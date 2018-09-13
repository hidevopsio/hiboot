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
	"github.com/hidevopsio/hiboot/examples/grpc/helloworld/protobuf"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/log"
	_ "github.com/hidevopsio/hiboot/pkg/starter/actuator"
	"github.com/hidevopsio/hiboot/pkg/starter/grpc"
	_ "github.com/hidevopsio/hiboot/pkg/starter/logging"
	"golang.org/x/net/context"
)

// controller
type helloController struct {
	// embedded web.Controller
	web.Controller
	// declare HelloServiceClient
	helloServiceClient protobuf.HelloServiceClient
}

// Init inject helloServiceClient
func newHelloController(helloServiceClient protobuf.HelloServiceClient) *helloController {
	return &helloController{
		helloServiceClient: helloServiceClient,
	}
}

// GET /greeter/name/{name}
func (c *helloController) GetByName(name string) string {

	// call grpc server method
	// pass context.Background() for the sake of simplicity
	response, err := c.helloServiceClient.SayHello(context.Background(), &protobuf.HelloRequest{Name: name})

	// got response
	if err == nil {
		return response.Message
	}

	// response with err
	return err.Error()
}

// controller
type holaController struct {
	// embedded web.Controller
	web.Controller
	// declare HolaServiceClient
	holaServiceClient protobuf.HolaServiceClient
}

// Init inject holaServiceClient
func newHolaController(holaServiceClient protobuf.HolaServiceClient) *holaController {
	return &holaController{
		holaServiceClient: holaServiceClient,
	}
}

// GET /greeter/name/{name}
func (c *holaController) GetByName(name string) string {

	// call grpc server method
	// pass context.Background() for the sake of simplicity
	response, err := c.holaServiceClient.SayHola(context.Background(), &protobuf.HolaRequest{Name: name})

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
	// 	   greeter-services:   # client name
	//       host: localhost # server host
	//       port: 7575      # server port
	//
	grpc.Client("hello-world-service", protobuf.NewHelloServiceClient)
	grpc.Client("hello-world-service", protobuf.NewHolaServiceClient)

	// must: register greeterController
	web.RestController(newHelloController)
	web.RestController(newHolaController)
}

func main() {
	// create new web application and run it
	err := web.NewApplication().Run()
	log.Debug(err)
}
