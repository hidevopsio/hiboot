package mockgrpc

import (
	"fmt"
	"github.com/golang/protobuf/proto"
)

// RPCMsg implements the gomock.Matcher interface
type RPCMsg struct {
	Message proto.Message
}

// Matches return matches message
func (r *RPCMsg) Matches(msg interface{}) bool {
	m, ok := msg.(proto.Message)
	if !ok {
		return false
	}
	return proto.Equal(m, r.Message)
}

// String return message in string
func (r *RPCMsg) String() string {
	return fmt.Sprintf("is %s", r.Message)
}
