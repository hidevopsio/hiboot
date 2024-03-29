// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/hidevopsio/hiboot/examples/grpc/helloworld/protobuf (interfaces: HolaServiceClient)

// Package mock_protobuf is a generated GoMock package.
package mock

import (
	"context"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc"
	protobuf2 "github.com/hidevopsio/hiboot/examples/web/grpc/helloworld/protobuf"
	"reflect"
)

// MockHolaServiceClient is a mock of HolaServiceClient interface
type MockHolaServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockHolaServiceClientMockRecorder
}

// MockHolaServiceClientMockRecorder is the mock recorder for MockHolaServiceClient
type MockHolaServiceClientMockRecorder struct {
	mock *MockHolaServiceClient
}

// NewMockHolaServiceClient creates a new mock instance
func NewMockHolaServiceClient(ctrl *gomock.Controller) *MockHolaServiceClient {
	mock := &MockHolaServiceClient{ctrl: ctrl}
	mock.recorder = &MockHolaServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockHolaServiceClient) EXPECT() *MockHolaServiceClientMockRecorder {
	return m.recorder
}

// SayHola mocks base method
func (m *MockHolaServiceClient) SayHola(arg0 context.Context, arg1 *protobuf2.HolaRequest, arg2 ...grpc.CallOption) (*protobuf2.HolaReply, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SayHola", varargs...)
	ret0, _ := ret[0].(*protobuf2.HolaReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SayHola indicates an expected call of SayHola
func (mr *MockHolaServiceClientMockRecorder) SayHola(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SayHola", reflect.TypeOf((*MockHolaServiceClient)(nil).SayHola), varargs...)
}
