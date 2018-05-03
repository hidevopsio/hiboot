package web

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/kataras/iris/context"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/kataras/iris"
	"github.com/kataras/iris/httptest"
	"net/http"
)



type RootController struct{
	Controller
}


type FooController struct{
	Controller
}

type BarController struct{
	Controller
}

func (c *RootController) GetHealth(ctx context.Context)  {
	health := Health{
		Status: "UP",
	}
	ctx.JSON(health)
}

func (c *FooController) PostSayHello(ctx context.Context)  {
	log.Print("SayHello")
}

func (c *BarController) GetSayHello(ctx context.Context)  {
	log.Print("SayHello")
}

type Controllers struct{

	Foo *FooController `controller:"foo" auth:"anon"`
	Bar *BarController `controller:"bar" auth:"anon"`
}

func TestApplicationHealth(t *testing.T) {
	app := iris.New()

	rc := &RootController{}

	urlPath := "/health"
	app.Handle(http.MethodGet, urlPath, rc.GetHealth)

	e := httptest.New(t, app)

	e.Request("GET", "/health").Expect().Status(http.StatusOK).Body().Contains("UP")
}


// HiCmdApplication
// HiWebApplication

func TestNewApplication(t *testing.T)  {
	controllers := &Controllers{}
	log.Println("PostSayHello: ", controllers.Foo.PostSayHello)
	log.Println("GetSayHello: ", controllers.Bar.GetSayHello)
	wa, err := NewWebApplication(controllers)
	assert.Equal(t, nil, err)

	e := httptest.New(t, wa.App())

	e.Request("POST", "/foo/sayHello").Expect().Status(http.StatusOK)
	e.Request("GET", "/bar/sayHello").Expect().Status(http.StatusOK)
}
