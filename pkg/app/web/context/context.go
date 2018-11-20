package context

import "github.com/kataras/iris/context"

// ExtendedContext extended context
type ExtendedContext interface {
	RequestBody(data interface{}) error
	RequestForm(data interface{}) error
	ResponseBody(message string, data interface{})
	ResponseError(message string, code int)
	RequestParams(request interface{}) error
	ResponseString(s string)
}

type Context interface {
	context.Context
	ExtendedContext
}
