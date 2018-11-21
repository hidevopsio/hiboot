package context

import "github.com/kataras/iris/context"

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
