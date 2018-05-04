package web

import (
	"github.com/kataras/iris/context"
	"github.com/hidevopsio/hiboot/pkg/model"
	"github.com/hidevopsio/hiboot/pkg/utils"
	"net/http"
)

type ContextInterface interface {
	RequestBody(data interface{}) error
	Response(message string, data interface{})
	ResponseError(message string, code int)
}

// Create your own custom Context, put any fields you wanna need.
type Context struct {
	// Optional Part 1: embed (optional but required if you don't want to override all context's methods)
	context.Context // it's the context/context.go#context struct but you don't need to know it.
	ContextInterface

}

var _ context.Context = &Context{} // optionally: validate on compile-time if Context implements context.Context.


// The only one important if you will override the Context
// with an embedded context.Context inside it.
// Required in order to run the handlers via this "*Context".
func (ctx *Context) Do(handlers context.Handlers) {
	context.Do(ctx, handlers)
}

// The second one important if you will override the Context
// with an embedded context.Context inside it.
// Required in order to run the chain of handlers via this "*Context".
func (ctx *Context) Next() {
	context.Next(ctx)
}

// Override any context's method you want...
// [...]

func (ctx *Context) HTML(htmlContents string) (int, error) {
	ctx.Application().Logger().Infof("Executing .HTML function from Context")

	ctx.ContentType("text/html")
	return ctx.WriteString(htmlContents)
}


// get response
func (ctx *Context) RequestBody(data interface{}) error {
	err := ctx.ReadJSON(&data)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusInternalServerError)
		return err
	}

	err = utils.Validate.Struct(data)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusBadRequest)
		return err
	}
	return nil
}

// set response
func (ctx *Context) Response(message string, data interface{}) {
	response := &model.Response{
		Code:    ctx.GetStatusCode(),
		Message: message,
		Data:    data,
	}

	// just for debug now
	ctx.JSON(response)
}

// set response
func (ctx *Context) ResponseError(message string, code int) {
	response := &model.Response{
		Code:    code,
		Message: message,
	}

	// just for debug now
	ctx.StatusCode(code)
	ctx.JSON(response)
}



