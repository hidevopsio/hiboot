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
	"reflect"
	"github.com/hidevopsio/hiboot/pkg/model"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"net/http"
	"strings"
	"strconv"
	"path/filepath"
	"github.com/hidevopsio/hiboot/pkg/utils/replacer"
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
	path = filepath.Clean(path)
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

	lenOfPathParams := len(h.pathParams)
	for i := 1; i < h.numIn; i++ {
		typ := method.Type.In(i)
		iTyp := reflector.IndirectType(typ)
		h.requests[i].typ = typ
		h.requests[i].iTyp = iTyp
		h.requests[i].kind = iTyp.Kind()
		h.requests[i].val = reflect.New(typ)
		h.requests[i].iVal = reflect.New(iTyp)
		h.requests[i].genKind = reflector.GetKind(h.requests[i].iVal) // TODO:
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
}

func (h *handler) call(ctx *Context) {

	var request interface{}
	var reqErr error
	var err error
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
			strVal := pvs[req.pathIdx] //TODO: out of scope
			var val interface{}
			switch req.typeName {
			case "int", "int16", "int32", "int64":
				val, err = strconv.Atoi(strVal)
				if err != nil {
					log.Error(err)
					return
				}
			case "uint", "uint16", "uint32", "uint64":
				val, err = strconv.ParseUint(strVal, 10, 64)
				if err != nil {
					log.Error(err)
					return
				}
			default:
				val = strVal
			}

			h.inputs[i] = reflect.ValueOf(val)
		} else {
			log.Warn("Not implemented!")
			return
		}
	}

	var respErr error
	var results []reflect.Value
	if reqErr == nil {
		reflector.SetFieldValue(h.controller, "Ctx", ctx)
		results = h.method.Func.Call(h.inputs)
		if h.numOut == 1 {
			result := results[0]
			if result.CanInterface() {
				if h.responses[0].kind == reflect.String {
					ctx.ResponseString(result.Interface().(string))
				} else if h.responses[0].typeName == "error" && !result.IsNil() {
					respErr = result.Interface().(error)
					ctx.ResponseError(respErr.Error(), http.StatusInternalServerError)
				}
			}
		} else if h.numOut >= 2 {
			if h.method.Type.Out(1).Name() == "error" {
				if results[1].IsNil() {
					respErr = nil
				} else {
					respErr = results[1].Interface().(error)
				}
			}
			var response model.Response
			switch h.responses[0].typeName {
			case "Response":
				if results[0].Interface() == nil {
					// TODO: add unit test
					log.Warn("response is nil")
					return
				} else {
					response = results[0].Interface().(model.Response)
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
			}
			ctx.JSON(response)
		}
	}
}
