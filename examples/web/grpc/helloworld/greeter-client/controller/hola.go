package controller

import (
	"golang.org/x/net/context"
	protobuf2 "hidevops.io/hiboot/examples/web/grpc/helloworld/protobuf"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/starter/grpc"
)

// controller
type holaController struct {
	// embedded at.RestController
	at.RestController
	// declare HolaServiceClient
	holaServiceClient protobuf2.HolaServiceClient
}

// Init inject holaServiceClient
func newHolaController(holaServiceClient protobuf2.HolaServiceClient) *holaController {
	return &holaController{
		holaServiceClient: holaServiceClient,
	}
}

// GET /greeter/name/{name}
func (c *holaController) GetByName(name string) (response string) {

	// call grpc server method
	// pass context.Background() for the sake of simplicity
	result, err := c.holaServiceClient.SayHola(context.Background(), &protobuf2.HolaRequest{Name: name})

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
		protobuf2.NewHolaServiceClient)

	// must: register Controller
	app.Register(newHolaController)
}
