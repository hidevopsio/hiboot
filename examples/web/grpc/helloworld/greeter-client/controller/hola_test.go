package controller

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/hidevopsio/hiboot/examples/web/grpc/helloworld/mock"
	"github.com/hidevopsio/hiboot/examples/web/grpc/helloworld/protobuf"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/starter/grpc/mockgrpc"
	"github.com/hidevopsio/hiboot/pkg/starter/logging"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestHolaClient(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockHolaClient := mock.NewMockHolaServiceClient(ctrl)

	req := &protobuf.HolaRequest{Name: "Steve"}

	mockHolaClient.EXPECT().SayHola(
		gomock.Any(),
		&mockgrpc.RPCMsg{Message: req},
	).Return(&protobuf.HolaReply{Message: "Hola " + req.Name}, nil)

	t.Run("test hello grpc service", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := mockHolaClient.SayHola(ctx, req)
		assert.Equal(t, nil, err)
		assert.Equal(t, "Hola Steve", r.Message)
	})

}

func TestHolaController(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockHolaClient := mock.NewMockHolaServiceClient(ctrl)

	req := &protobuf.HolaRequest{Name: "Steve"}

	mockHolaClient.EXPECT().SayHola(
		gomock.Any(),
		&mockgrpc.RPCMsg{Message: req},
	).Return(&protobuf.HolaReply{Message: "Hola " + req.Name}, nil)

	t.Run("test hello controller", func(t *testing.T) {
		holaCtrl := newHolaController(mockHolaClient)
		testApp := web.NewTestApp(t, holaCtrl).SetProperty(logging.Level, logging.LevelDebug).Run(t)
		testApp.Get("/hola/{name}").
			WithPath("name", req.Name).
			Expect().Status(http.StatusOK).
			Body().Contains(req.Name)
	})

}
