// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package web

import (
	ctx "github.com/kataras/iris/context"
	"github.com/kataras/iris/middleware/i18n"
	"hidevops.io/hiboot/pkg/app/web/context"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/model"
	"hidevops.io/hiboot/pkg/utils/mapstruct"
	"hidevops.io/hiboot/pkg/utils/validator"
	"net/http"
)

// Context Create your own custom Context, put any fields you wanna need.
type Context struct {
	// Optional Part 1: embed (optional but required if you don't want to override all context's methods)
	// it's the context/context.go#context struct but you don't need to know it.
	ctx.Context
}

var _ ctx.Context = &Context{} // optionally: validate on compile-time if Context implements context.Context.

// NewContext constructor of context.Context
func NewContext(app ctx.Application) context.Context {
	return &Context{
		Context: ctx.NewContext(app),
	}
}

// Do The only one important if you will override the Context
// with an embedded context.Context inside it.
// Required in order to run the handlers via this "*Context".
func (c *Context) Do(handlers ctx.Handlers) {
	log.Debug("Context.Do()")
	ctx.Do(c, handlers)
}

// Next The second one important if you will override the Context
// with an embedded context.Context inside it.
// Required in order to run the chain of handlers via this "*Context".
func (c *Context) Next() {
	ctx.Next(c)
}

// HTML Override any context's method you want...
// [...]
func (c *Context) HTML(htmlContents string) (int, error) {
	c.Application().Logger().Infof("Executing .HTML function from Context")

	c.ContentType("text/html")
	return c.WriteString(htmlContents)
}

//// RequestBody get RequestBody
//func (c *Context) RequestBody(data interface{}) error {
//	return RequestBody(c, data)
//}
//
//// RequestForm get RequestFrom
//func (c *Context) RequestForm(data interface{}) error {
//	return RequestForm(c, data)
//}
//
//// RequestParams get RequestParams
//func (c *Context) RequestParams(data interface{}) error {
//	return RequestParams(c, data)
//}

// handle i18n
func (c *Context) translate(message string) string {

	message = i18n.Translate(c, message)

	return message
}

// Translate override base context method Translate to return format if i18n is not enabled
func (c *Context) Translate(format string, args ...interface{}) string {

	msg := c.Context.Translate(format, args...)

	if msg == "" {
		msg = format
	}

	return msg
}

// ResponseString set response
func (c *Context) ResponseString(data string) {
	c.WriteString(c.translate(data))
}

// ResponseBody set response
func (c *Context) ResponseBody(message string, data interface{}) {

	// TODO: check if data is a string, should we translate it?
	response := new(model.BaseResponse)
	response.SetCode(c.GetStatusCode())
	response.SetMessage(c.translate(message))
	response.SetData(data)

	c.JSON(response)
}

// ResponseError response with error
func (c *Context) ResponseError(message string, code int) {

	response := new(model.BaseResponse)
	response.SetCode(code)
	response.SetMessage(c.translate(message))

	c.StatusCode(code)
	c.JSON(response)
}

// RequestEx get RequestBody
func requestEx(c context.Context, data interface{}, cb func() error) error {
	if cb != nil {
		err := cb()
		if err != nil {
			c.ResponseError(err.Error(), http.StatusInternalServerError)
			return err
		}

		err = validator.Validate.Struct(data)
		if err != nil {
			c.ResponseError(err.Error(), http.StatusBadRequest)
			return err
		}
	}
	return nil
}

// RequestBody get RequestBody
func RequestBody(c context.Context, data interface{}) error {

	return requestEx(c, data, func() error {
		return c.ReadJSON(data)
	})
}

// RequestForm get RequestFrom
func RequestForm(c context.Context, data interface{}) error {

	return requestEx(c, data, func() error {
		return c.ReadForm(data)
	})
}

// RequestParams get RequestParams
func RequestParams(c context.Context, data interface{}) error {

	return requestEx(c, data, func() error {

		values := c.URLParams()
		if len(values) != 0 {
			return mapstruct.Decode(data, values)
		}
		return nil
	})
}
