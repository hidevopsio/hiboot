package web

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/kataras/iris/context"
	"github.com/hidevopsio/hiboot/pkg/log"
)



type FooController struct{
	Controller
}
type BarController struct{
	Controller
}

func (c *FooController) PostSayHello(ctx context.Context)  {
	log.Print("SayHello")
}

func (c *BarController) GetSayHello(ctx context.Context)  {
	log.Print("SayHello")
}


type Controllers struct{
	Foo *FooController `controller:"foo",auth:"anon"`
	Bar *BarController `controller:"bar"`
}

func TestNewApplication(t *testing.T)  {
	_, err := NewApplication(&Controllers{})
	assert.Equal(t, nil, err)
}
