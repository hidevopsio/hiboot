// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/hidevopsio/hiboot/examples/grpc/helloworld/protobuf (interfaces: HelloServiceClient)

// Package mock_protobuf is a generated GoMock package.
package mock

import (
	"context"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc"
	protobuf2 "github.com/hidevopsio/hiboot/examples/web/grpc/helloworld/protobuf"
	"reflect"
)

// MockHelloServiceClient is a mock of HelloServiceClient interface
type MockHelloServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockHelloServiceClientMockRecorder
}

// MockHelloServiceClientMockRecorder is the mock recorder for MockHelloServiceClient
type MockHelloServiceClientMockRecorder struct {
	mock *MockHelloServiceClient
}

// NewMockHelloServiceClient creates a new mock instance
func NewMockHelloServiceClient(ctrl *gomock.Controller) *MockHelloServiceClient {
	mock := &MockHelloServiceClient{ctrl: ctrl}
	mock.recorder = &MockHelloServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockHelloServiceClient) EXPECT() *MockHelloServiceClientMockRecorder {
	return m.recorder
}

// SayHello mocks base method
func (m *MockHelloServiceClient) SayHello(arg0 context.Context, arg1 *protobuf2.HelloRequest, arg2 ...grpc.CallOption) (*protobuf2.HelloReply, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SayHello", varargs...)
	ret0, _ := ret[0].(*protobuf2.HelloReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SayHello indicates an expected call of SayHello
func (mr *MockHelloServiceClientMockRecorder) SayHello(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SayHello", reflect.TypeOf((*MockHelloServiceClient)(nil).SayHello), varargs...)
}
