package mock_protobuf

import (
	"fmt"
	"github.com/golang/protobuf/proto"
)

// RpcMsg implements the gomock.Matcher interface
type RpcMsg struct {
	Message proto.Message
}

// Matches return matches message
func (r *RpcMsg) Matches(msg interface{}) bool {
	m, ok := msg.(proto.Message)
	if !ok {
		return false
	}
	return proto.Equal(m, r.Message)
}

// String return message in string
func (r *RpcMsg) String() string {
	return fmt.Sprintf("is %s", r.Message)
}
