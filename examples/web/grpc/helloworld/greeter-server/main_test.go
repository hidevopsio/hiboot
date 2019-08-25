package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	protobuf2 "hidevops.io/hiboot/examples/web/grpc/helloworld/protobuf"
	"testing"
	"time"
)

func TestRunMain(t *testing.T) {
	go main()

	time.Sleep(time.Second)
}

func TestHelloServer(t *testing.T) {

	serviceServer := newHelloServiceServer()

	name := "Steve"
	expected := "Hello " + name
	req := &protobuf2.HelloRequest{Name: name}
	res, err := serviceServer.SayHello(context.Background(), req)

	assert.Equal(t, nil, err)
	assert.Equal(t, expected, res.Message)
}

func TestHolaServer(t *testing.T) {

	serviceServer := newHolaServiceServer()

	name := "Steve"
	expected := "Hola " + name
	req := &protobuf2.HolaRequest{Name: name}
	res, err := serviceServer.SayHola(context.Background(), req)

	assert.Equal(t, nil, err)
	assert.Equal(t, expected, res.Message)
}
