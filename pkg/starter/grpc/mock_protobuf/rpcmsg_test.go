package mock_protobuf

import (
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/examples/helloworld/helloworld"
	"testing"
)

func TestRpcMsg(t *testing.T) {
	rpcMsg := new(RpcMsg)
	req := &helloworld.HelloRequest{Name: "unit_test"}

	ok := rpcMsg.Matches(nil)
	assert.Equal(t, false, ok)

	rpcMsg.Message = req
	assert.Equal(t, "is name:\"unit_test\" ", rpcMsg.String())

	ok = rpcMsg.Matches(req)
	assert.Equal(t, true, ok)

}
