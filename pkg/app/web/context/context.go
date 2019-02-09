package context

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

// ExtendedContext extended context
type ExtendedContext interface {
	//RequestBody(data interface{}) error
	//RequestForm(data interface{}) error
	//RequestParams(request interface{}) error
	ResponseString(s string)
	ResponseBody(message string, data interface{})
	ResponseError(message string, code int)
}

type Context interface {
	context.Context
	ExtendedContext
}

type Handler func(Context)

// NewHandler will convert iris handler to our handler of func(*Context),
// in order to be compatible with the HTTP API.
func NewHandler(h iris.Handler) Handler {
	return func(ctx Context) {
		h(ctx.(context.Context))
	}
}