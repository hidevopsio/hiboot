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

// Create your own custom Context, put any fields you wanna need.
type MyContext struct {
	// Optional Part 1: embed (optional but required if you don't want to override all context's methods)
	context.Context // it's the context/context.go#context struct but you don't need to know it.
}

var _ context.Context = &MyContext{} // optionally: validate on compile-time if MyContext implements context.Context.

// The only one important if you will override the Context
// with an embedded context.Context inside it.
// Required in order to run the handlers via this "*MyContext".
func (ctx *MyContext) Do(handlers context.Handlers) {
	context.Do(ctx, handlers)
}

// The second one important if you will override the Context
// with an embedded context.Context inside it.
// Required in order to run the chain of handlers via this "*MyContext".
func (ctx *MyContext) Next() {
	context.Next(ctx)
}

// Override any context's method you want...
// [...]

func (ctx *MyContext) HTML(htmlContents string) (int, error) {
	ctx.Application().Logger().Infof("Executing .HTML function from MyContext")

	ctx.ContentType("text/html")
	return ctx.WriteString(htmlContents)
}


func (ctx *MyContext) Foo() {
	log.Print("foo")
}

func TestCustomContext(t *testing.T) {
	app := iris.New()
	// app.Logger().SetLevel("debug")

	// The only one Required:
	// here is how you define how your own context will
	// be created and acquired from the iris' generic context pool.
	app.ContextPool.Attach(func() context.Context {
		return &MyContext{
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

// should always print "($PATH) Handler is executing from 'MyContext'"
func recordWhichContextJustForProofOfConcept(ctx context.Context) {
	ctx.Application().Logger().Infof("(%s) Handler is executing from: '%s'", ctx.Path(), reflect.TypeOf(ctx).Elem().Name())
	ctx.Next()
}

