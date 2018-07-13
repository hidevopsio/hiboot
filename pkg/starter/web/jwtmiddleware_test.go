package web

import (
	"github.com/hidevopsio/hiboot/pkg/log"
)

var (
	jm *JwtMiddleware
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

	jm = new(JwtMiddleware)
}
