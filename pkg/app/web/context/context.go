package context

import (
	"net/http"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

// ExtendedContext extended context
type ExtendedContext interface {
	ResponseString(s string)
	ResponseBody(message string, data interface{})
	ResponseError(message string, code int)
	WrapHandler(h http.Handler)
	SetAnnotations(ann interface{})
	SetURLParam(name, value string)
	Annotations() interface{}
	InitResponses()
	SetResponse(idx int, response interface{})
	GetResponses() (responses []interface{})
	GetResponse(idx int) (response interface{})
//StaticResource(system http.FileSystem)
}

// Context is the interface of web app context
type Context interface {
	context.Context
	ExtendedContext
}

// Handler is the handler func type (for Middleware)
type Handler func(Context)

// NewHandler will convert iris handler to our handler of func(*Context),
// in order to be compatible with the HTTP API.
func NewHandler(h iris.Handler) Handler {
	return func(ctx Context) {
		h(ctx.(context.Context))
	}
}
