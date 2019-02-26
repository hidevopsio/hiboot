package controller

import (
	"golang.org/x/net/context"
	"hidevops.io/hiboot/examples/grpc/helloworld/protobuf"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/starter/grpc"
)

// controller
type holaController struct {
	// embedded at.RestController
	at.RestController
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
func (c *holaController) GetByName(name string) (response string) {

	// call grpc server method
	// pass context.Background() for the sake of simplicity
	result, err := c.holaServiceClient.SayHola(context.Background(), &protobuf.HolaRequest{Name: name})

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
		protobuf.NewHolaServiceClient)

	// must: register Controller
	app.Register(newHolaController)
}
