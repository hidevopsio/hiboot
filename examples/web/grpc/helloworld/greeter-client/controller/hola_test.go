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

func TestHolaClient(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockHolaClient := mock2.NewMockHolaServiceClient(ctrl)
	app.Register("protobuf.holaServiceClient", mockHolaClient)

	testApp := web.NewTestApp(t, newHolaController).SetProperty(logging.Level, logging.LevelDebug).Run(t)

	req := &protobuf2.HolaRequest{Name: "Steve"}

	mockHolaClient.EXPECT().SayHola(
		gomock.Any(),
		&mockgrpc.RPCMsg{Message: req},
	).Return(&protobuf2.HolaReply{Message: "Hola " + req.Name}, nil)

	testApp.Get("/hola/name/{name}").
		WithPath("name", req.Name).
		Expect().Status(http.StatusOK).
		Body().Contains(req.Name)
}
