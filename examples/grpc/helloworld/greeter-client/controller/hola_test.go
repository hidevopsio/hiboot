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

func TestHolaClient(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockHolaClient := mock.NewMockHolaServiceClient(ctrl)
	app.Register("protobuf.holaServiceClient", mockHolaClient)

	testApp := web.RunTestApplication(t, newHolaController)

	req := &protobuf.HolaRequest{Name: "Steve"}

	mockHolaClient.EXPECT().SayHola(
		gomock.Any(),
		&grpcmock.RPCMsg{Message: req},
	).Return(&protobuf.HolaReply{Message: "Hola " + req.Name}, nil)

	testApp.Get("/hola/name/{name}").
		WithPath("name", req.Name).
		Expect().Status(http.StatusOK).
		Body().Contains(req.Name)
}
