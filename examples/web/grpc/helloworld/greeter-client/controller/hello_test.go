package controller

import (
	"github.com/golang/mock/gomock"
	mock2 "hidevops.io/hiboot/examples/web/grpc/helloworld/mock"
	protobuf2 "hidevops.io/hiboot/examples/web/grpc/helloworld/protobuf"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/starter/grpc/mockgrpc"
	"hidevops.io/hiboot/pkg/starter/logging"
	"net/http"
	"testing"
)

func TestHelloClient(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockHelloClient := mock2.NewMockHelloServiceClient(ctrl)
	app.Register("protobuf.helloServiceClient", mockHelloClient)

	testApp := web.NewTestApp(t, newHelloController).SetProperty(logging.Level, logging.LevelDebug).Run(t)

	req := &protobuf2.HelloRequest{Name: "Steve"}

	mockHelloClient.EXPECT().SayHello(
		gomock.Any(),
		&mockgrpc.RPCMsg{Message: req},
	).Return(&protobuf2.HelloReply{Message: "Hello " + req.Name}, nil)

	testApp.Get("/hello/name/{name}").
		WithPath("name", req.Name).
		Expect().Status(http.StatusOK).
		Body().Contains(req.Name)
}
