package web

import (
	"github.com/hidevopsio/hiboot/pkg/log"
	"testing"
)

var (
	jwtMw *JwtMiddleware
	ctx FakeContext
)

type FakeContext struct {
}

func (c *FakeContext) Next()  {
	log.Debug("FakeContext.Next()")
}

func (c *FakeContext) StopExecution()  {
	log.Debug("FakeContext.Next()")
}

func init() {
	log.SetLevel(log.DebugLevel)

	jwtMw = new(JwtMiddleware)
}

func TestCheckJWT(t *testing.T) {
}
