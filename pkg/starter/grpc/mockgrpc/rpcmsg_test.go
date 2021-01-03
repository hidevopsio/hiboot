package mockgrpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/examples/helloworld/helloworld"
)

func TestRpcMsg(t *testing.T) {
	rpcMsg := new(RPCMsg)
	req := &helloworld.HelloRequest{Name: "unit_test"}

	ok := rpcMsg.Matches(nil)
	assert.Equal(t, false, ok)

	rpcMsg.Message = req
	assert.Contains(t, rpcMsg.Message.String(), "unit_test", )

	ok = rpcMsg.Matches(req)
	assert.Equal(t, true, ok)

}
