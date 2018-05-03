package web

import (
	"github.com/kataras/iris/context"
	"github.com/kataras/iris"
	"reflect"
	"testing"
	"github.com/kataras/iris/httptest"
	"net/http"
	"github.com/hidevopsio/hiboot/pkg/log"
)

// should always print "($PATH) Handler is executing from 'Context'"
func recordWhichContextJustForProofOfConcept(ctx context.Context) {
	ctx.Application().Logger().Infof("(%s) Handler is executing from: '%s'", ctx.Path(), reflect.TypeOf(ctx).Elem().Name())
	ctx.Next()
}


func TestCustomContext(t *testing.T) {
	log.Debug("TestCustomContext")
	app := iris.New()
	// app.Logger().SetLevel("debug")

	// The only one Required:
	// here is how you define how your own context will
	// be created and acquired from the iris' generic context pool.
	app.ContextPool.Attach(func() context.Context {
		return &Context{
			// Optional Part 3:
			Context: context.NewContext(app),
		}
	})

	// register your route, as you normally do
	app.Handle("GET", "/health", recordWhichContextJustForProofOfConcept, func(ctx context.Context) {
		// use the context's overridden HTML method.
		health := Health{
			Status: "UP",
		}

		ctx.JSON(health)
	})

	e := httptest.New(t, app)

	e.Request("GET", "/health").Expect().Status(http.StatusOK).Body()

}


