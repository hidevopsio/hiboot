package controller

import (
	"golang.org/x/net/context"
	protobuf2 "hidevops.io/hiboot/examples/web/grpc/helloworld/protobuf"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/starter/grpc"
)

// controller
type helloController struct {
	// embedded at.RestController
	at.RestController
	// declare HelloServiceClient
	helloServiceClient protobuf2.HelloServiceClient
}

// Init inject helloServiceClient
func newHelloController(helloServiceClient protobuf2.HelloServiceClient) *helloController {
	return &helloController{
		helloServiceClient: helloServiceClient,
	}
}

// GET /greeter/name/{name}
func (c *helloController) GetByName(name string) (response string) {

	// call grpc server method
	// pass context.Background() for the sake of simplicity
	result, err := c.helloServiceClient.SayHello(context.Background(), &protobuf2.HelloRequest{Name: name})

	// got response
	if err == nil {
		response = result.Message
	}
	return
}

func init() {

	// must: register grpc client, the name greeter-client should configured in application.yml
	// see config/application-grpc.yml
	//
	// grpc:
	//   client:
	// 	   hello-world-service:   # client name
	//       host: localhost # server host
	//       port: 7575      # server port
	//
	grpc.Client("hello-world-service",
		protobuf2.NewHelloServiceClient)

	// must: register Rest Controller
	app.Register(newHelloController)
}
