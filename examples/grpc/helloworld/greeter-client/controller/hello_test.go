package controller

import (
	"github.com/golang/mock/gomock"
	"github.com/hidevopsio/hiboot/examples/grpc/helloworld/mock_protobuf"
	"github.com/hidevopsio/hiboot/examples/grpc/helloworld/protobuf"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	mockproto "github.com/hidevopsio/hiboot/pkg/starter/grpc/mock"
	"net/http"
	"testing"
)

func TestHelloClient(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockHelloClient := mock_protobuf.NewMockHelloServiceClient(ctrl)
	app.Register("protobuf.helloServiceClient", mockHelloClient)

	testApp := web.NewTestApplication(t, newHelloController)

	req := &protobuf.HelloRequest{Name: "Steve"}

	mockHelloClient.EXPECT().SayHello(
		gomock.Any(),
		&mockproto.RPCMsg{Message: req},
	).Return(&protobuf.HelloReply{Message: "Hello " + req.Name}, nil)

	testApp.Get("/hello/name/{name}").
		WithPath("name", req.Name).
		Expect().Status(http.StatusOK).
		Body().Contains(req.Name)
}
