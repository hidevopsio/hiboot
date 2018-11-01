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
	"fmt"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/model"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"github.com/hidevopsio/hiboot/pkg/utils/replacer"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
	"net/http"
	"reflect"
	"strings"
)

type request struct {
	typeName string
	name     string
	kind     reflect.Kind
	genKind  reflect.Kind // e.g. convert int16 to int
	typ      reflect.Type
	iTyp     reflect.Type
	val      reflect.Value
	iVal     reflect.Value
	pathIdx  int
}

type response struct {
	typeName string
	name     string
	kind     reflect.Kind
	typ      reflect.Type
}

type handler struct {
	controller      interface{}
	method          reflect.Method
	inputs          []reflect.Value
	numIn           int
	numOut          int
	pathParams      []string
	requests        []request
	responses       []response
	lenOfPathParams int
	hasCtxField     bool
}

func clean(in string) (out string) {
	out = strings.Replace(in, "//", "/", -1)
	if strings.Contains(out, "//") {
		out = clean(out)
	}
	lenOfOut := len(out) - 1
	if lenOfOut > 1 && out[lenOfOut:] == "/" {
		//log.Debug(out[:lenOfOut])
		out = out[:lenOfOut]
	}
	return
}

func (h *handler) parse(method reflect.Method, object interface{}, path string) {
	//log.Debug("NumIn: ", method.Type.NumIn())
	h.controller = object
	h.method = method
	h.numIn = method.Type.NumIn()
	h.numOut = method.Type.NumOut()
	h.inputs = make([]reflect.Value, h.numIn)
	h.inputs[0] = reflect.ValueOf(object)

	//log.Debugf("method: %v", method.Name)

	// TODO: should parse all of below request and response during router register to improve performance
	path = clean(path)
	//log.Debugf("path: %v", path)
	pps := strings.SplitN(path, "/", -1)
	//log.Debug(pps)
	pp := replacer.ParseVariables(path, compiledRegExp)
	h.pathParams = make([]string, len(pp))
	for i, pathParam := range pp {
		//log.Debugf("pathParm: %v", pathParam[1])
		h.pathParams[i] = pathParam[1]
	}

	typeNames := []string{
		model.RequestTypeForm,
		model.RequestTypeParams,
		model.RequestTypeBody,
		model.Context,
	}

	h.requests = make([]request, h.numIn)
	objVal := reflect.ValueOf(object)
	idv := reflect.Indirect(objVal)
	objTyp := idv.Type()
	h.requests[0].typeName = objTyp.Name()
	h.requests[0].typ = objTyp
	h.requests[0].val = objVal
	h.hasCtxField = reflector.HasEmbeddedFieldType(object, Controller{})

	lenOfPathParams := len(h.pathParams)
	for i := 1; i < h.numIn; i++ {
		typ := method.Type.In(i)
		iTyp := reflector.IndirectType(typ)
		h.requests[i].typ = typ
		h.requests[i].iTyp = iTyp
		if typ.Kind() == reflect.Slice {
			h.requests[i].kind = reflect.Slice
		} else {
			h.requests[i].kind = iTyp.Kind()
		}
		h.requests[i].val = reflect.New(typ)
		h.requests[i].iVal = reflect.New(iTyp)
		h.requests[i].genKind = reflector.GetKindByValue(h.requests[i].iVal) // TODO:

		pi := i - 1
		if pi < lenOfPathParams {
			h.requests[i].name = pp[pi][1]
			for idx, pv := range pps {
				if pv == pp[pi][0] {
					h.requests[i].pathIdx = idx
					break
				}
			}
		}
		h.requests[i].typeName = iTyp.Name()
		if iTyp.Kind() == reflect.Struct {
			for _, tn := range typeNames {
				if field, ok := iTyp.FieldByName(tn); ok && field.Anonymous {
					h.requests[i].typeName = tn
					break
				}
			}
		}
	}
	h.lenOfPathParams = lenOfPathParams

	h.responses = make([]response, h.numOut)
	for i := 0; i < h.numOut; i++ {
		typ := method.Type.Out(i)

		h.responses[i].typ = typ
		h.responses[i].kind = typ.Kind()
		h.responses[i].typeName = typ.Name()
		//log.Debug(h.responses[i])
	}
	if path != "" {
		log.Infof("Mapped \"%v\" onto %v.%v()", path, idv.Type(), method.Name)
	}
}

func (h *handler) responseData(ctx *Context, numOut int, results []reflect.Value) {
	if numOut == 0 {
		return
	}

	result := results[0]
	if !result.CanInterface() {
		ctx.ResponseError("response is invalid", http.StatusInternalServerError)
		return
	}

	respVal := result.Interface()
	if respVal == nil {
		log.Warn("response is nil")
		return
	}

	switch h.responses[0].typeName {
	case "string":
		ctx.ResponseString(result.Interface().(string))
	case "error":
		respErr := result.Interface().(error)
		ctx.ResponseError(respErr.Error(), http.StatusInternalServerError)
	case "Response":
		response := respVal.(model.Response)
		if numOut >= 2 {
			var respErr error
			errVal := results[1]
			if errVal.IsNil() {
				respErr = nil
			} else if errVal.Type().Name() == "error" {
				respErr = results[1].Interface().(error)
			}

			if respErr == nil {
				response.SetCode(http.StatusOK)
				response.SetMessage(ctx.translate("success"))
			} else {
				if response.GetCode() == 0 {
					response.SetCode(http.StatusInternalServerError)
				}
				// TODO: output error message directly? how about i18n
				response.SetMessage(ctx.translate(respErr.Error()))

				// TODO: configurable status code in application.yml
				ctx.StatusCode(response.GetCode())
			}
		}
		ctx.JSON(response)
	default:
		ctx.ResponseError("response type is not implemented!", http.StatusInternalServerError)
	}
}

func (h *handler) call(ctx *Context) {

	var request interface{}
	var reqErr error
	var path string
	var pvs []string

	if h.lenOfPathParams != 0 {
		path = ctx.Path()
		//log.Debugf("path: %v", path)
		pvs = strings.SplitN(path, "/", -1)
	}

	for i := 1; i < h.numIn; i++ {
		req := h.requests[i]
		request = req.iVal.Interface()
		if req.kind == reflect.Struct {
			switch req.typeName {
			case model.RequestTypeForm:
				reqErr = ctx.RequestForm(request)
			case model.RequestTypeParams:
				reqErr = ctx.RequestParams(request)
			case model.RequestTypeBody:
				reqErr = ctx.RequestBody(request)
			case model.Context:
				request = ctx
			}

			if reqErr != nil {
				e := reqErr.Error()
				log.Error(e)
				return
			}

			h.inputs[i] = reflect.ValueOf(request)
			break
		} else if h.lenOfPathParams != 0 {
			strVal := pvs[req.pathIdx]
			val := str.Convert(strVal, req.kind)
			h.inputs[i] = reflect.ValueOf(val)
		} else {
			msg := fmt.Sprintf("input type: %v is not supported!", req.typ)
			ctx.ResponseError(msg, http.StatusInternalServerError)
			return
		}
	}

	//var respErr error
	var results []reflect.Value
	if reqErr == nil {
		if h.hasCtxField {
			reflector.SetFieldValue(h.controller, "Ctx", ctx)
		}
		// call controller method
		results = h.method.Func.Call(h.inputs)
		h.responseData(ctx, h.numOut, results)
	}
}
