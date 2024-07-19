package controller

import (
	"github.com/golang/mock/gomock"
	"github.com/hidevopsio/hiboot/examples/web/grpc/helloworld/mock"
	"github.com/hidevopsio/hiboot/examples/web/grpc/helloworld/protobuf"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/starter/grpc/mockgrpc"
	"github.com/hidevopsio/hiboot/pkg/starter/logging"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"net/http"
	"testing"
	"time"
)

func TestHelloClient(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockHelloClient := mock.NewMockHelloServiceClient(ctrl)

	req := &protobuf.HelloRequest{Name: "Steve"}
	mockHelloClient.EXPECT().SayHello(
		gomock.Any(),
		&mockgrpc.RPCMsg{Message: req},
	).Return(&protobuf.HelloReply{Message: "Hello " + req.Name}, nil)

	t.Run("test hello grpc service", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := mockHelloClient.SayHello(ctx, req)
		assert.Equal(t, nil, err)
		assert.Equal(t, "Hello Steve", r.Message)
	})
}

func TestHelloController(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockHelloClient := mock.NewMockHelloServiceClient(ctrl)

	req := &protobuf.HelloRequest{Name: "Steve"}
	mockHelloClient.EXPECT().SayHello(
		gomock.Any(),
		&mockgrpc.RPCMsg{Message: req},
	).Return(&protobuf.HelloReply{Message: "Hello " + req.Name}, nil)

	t.Run("test hello controller", func(t *testing.T) {
		helloCtrl := newHelloController(mockHelloClient)
		testApp := web.NewTestApp(t, helloCtrl).SetProperty(logging.Level, logging.LevelDebug).Run(t)
		testApp.Get("/hello/{name}").
			WithPath("name", req.Name).
			Expect().Status(http.StatusOK).
			Body().Contains(req.Name)
	})
}
