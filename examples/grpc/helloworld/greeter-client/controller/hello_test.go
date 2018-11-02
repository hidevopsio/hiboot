package controller

import (
	"github.com/golang/mock/gomock"
	"hidevops.io/hiboot/examples/grpc/helloworld/mock"
	"hidevops.io/hiboot/examples/grpc/helloworld/protobuf"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web"
	grpcmock "hidevops.io/hiboot/pkg/starter/grpc/mock"
	"net/http"
	"testing"
)

func TestHelloClient(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockHelloClient := mock.NewMockHelloServiceClient(ctrl)
	app.Register("protobuf.helloServiceClient", mockHelloClient)

	testApp := web.RunTestApplication(t, newHelloController)

	req := &protobuf.HelloRequest{Name: "Steve"}

	mockHelloClient.EXPECT().SayHello(
		gomock.Any(),
		&grpcmock.RPCMsg{Message: req},
	).Return(&protobuf.HelloReply{Message: "Hello " + req.Name}, nil)

	testApp.Get("/hello/name/{name}").
		WithPath("name", req.Name).
		Expect().Status(http.StatusOK).
		Body().Contains(req.Name)
}
