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
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/inject/annotation"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/model"
	"github.com/hidevopsio/hiboot/pkg/utils/mapstruct"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"github.com/hidevopsio/hiboot/pkg/utils/replacer"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
)

const (
	success       = "success"
	failed        = "failed"
	finalResponse = "finalResponse"
)

var (
	ErrCanNotInterface = errors.New("response can not interface")
)

type request struct {
	typeName     string
	name         string
	fullName     string
	kind         reflect.Kind
	genKind      reflect.Kind // e.g. convert int16 to int
	typ          reflect.Type
	iTyp         reflect.Type
	obj          interface{}
	iVal         reflect.Value
	pathIdx      int
	callback     func(ctx context.Context, data interface{}) error
	isAnnotation bool
}

type response struct {
	typeName               string
	name                   string
	kind                   reflect.Kind
	typ                    reflect.Type
	isResponseBody         bool
	implementsResponse     bool
	implementsResponseInfo bool
	value                  interface{}
}

type handler struct {
	object          interface{}
	method          *reflect.Method
	path            string
	objVal          reflect.Value
	numIn           int
	numOut          int
	pathVariable    []string
	requests        []request
	responses       []response
	lenOfPathParams int
	factory         factory.ConfigurableFactory
	runtimeInstance factory.Instance
	contextName     string
	dependencies    []*factory.MetaData

	injectableObject *injectableObject
	restMethod       *injectableMethod
	annotations      interface{}
}

type requestSet struct {
	name     string
	callback func(ctx context.Context, data interface{}) error
}

var requestSets []requestSet

func newRequestTypeName(in interface{}) string {
	return reflector.GetName(in)
}

func init() {
	requestSets = []requestSet{
		{newRequestTypeName(new(at.RequestForm)), RequestForm},
		{newRequestTypeName(new(at.RequestParams)), RequestParams},
		{newRequestTypeName(new(at.RequestBody)), RequestBody},
	}
}

func newHandler(factory factory.ConfigurableFactory, injectableObject *injectableObject, restMethod *injectableMethod, atType interface{}) *handler {
	hdl := &handler{
		contextName:      reflector.GetLowerCamelFullName(new(context.Context)),
		factory:          factory,
		injectableObject: injectableObject,
		restMethod:       restMethod,
	}

	hdl.parseMethod(injectableObject, restMethod, atType)
	return hdl
}

func clean(in string) (out string) {
	out = strings.Replace(in, "//", "/", -1)
	if strings.Contains(out, "//") {
		out = clean(out)
	}
	lenOfOut := len(out) - 1
	if lenOfOut > 1 && out[lenOfOut:] == "/" {
		out = out[:lenOfOut]
	}
	return
}

func (h *handler) parseMethod(injectableObject *injectableObject, injectableMethod *injectableMethod, atType interface{}) {
	//log.Debug("NumIn: ", method.Type.NumIn())
	method := injectableMethod.method
	path := ""
	if injectableMethod.requestMapping != nil && injectableMethod.requestMapping.Value != "" {
		path = injectableObject.pathPrefix + injectableMethod.requestMapping.Value
	}
	h.object = injectableObject.object
	h.method = method
	h.numIn = method.Type.NumIn()
	h.numOut = method.Type.NumOut()
	//h.inputs = make([]reflect.Value, h.numIn)
	objVal := reflect.ValueOf(h.object)
	h.objVal = objVal

	// TODO: should parse all of below request and response during router register to improve performance
	path = clean(path)
	h.path = path
	//log.Debugf("path: %v", path)
	pps := strings.SplitN(path, "/", -1)
	//log.Debug(pps)
	pp := replacer.ParseVariables(path, compiledRegExp)
	h.pathVariable = make([]string, len(pp))
	for i, pathParam := range pp {
		//log.Debugf("pathParm: %v", pathParam[1])
		h.pathVariable[i] = pathParam[1]
	}

	h.requests = make([]request, h.numIn)

	idv := reflect.Indirect(objVal)
	objTyp := idv.Type()
	h.requests[0].typeName = objTyp.Name()
	h.requests[0].typ = objTyp
	//h.requests[0].val = objVal

	lenOfPathParams := len(h.pathVariable)
	pathIdx := 0
	// parse request
	for i := 1; i < h.numIn; i++ {
		typ := method.Type.In(i)
		iTyp := reflector.IndirectType(typ)

		// parse embedded annotation at.ContextAware
		// append at.ContextAware dependencies
		dp := h.factory.GetInstance(iTyp, factory.MetaData{})
		if dp != nil {
			cdp := dp.(*factory.MetaData)
			if cdp.ContextAware {
				h.dependencies = append(h.dependencies, dp.(*factory.MetaData))
			}
		}

		// basic info
		h.requests[i].typ = typ
		h.requests[i].iTyp = iTyp
		if typ.Kind() == reflect.Slice {
			h.requests[i].kind = reflect.Slice
		} else {
			h.requests[i].kind = iTyp.Kind()
		}
		//h.requests[i].val = reflect.New(typ)
		h.requests[i].iVal = reflect.New(iTyp)
		h.requests[i].obj = h.requests[i].iVal.Interface()
		h.requests[i].typeName = iTyp.Name()
		h.requests[i].genKind = reflector.GetKindByValue(h.requests[i].iVal) // TODO:
		h.requests[i].fullName = reflector.GetLowerCamelFullNameByType(iTyp)

		// handle annotations
		request := h.requests[i].iVal.Interface()
		// TODO: use annotation.Contains(request, at.Annotation{}) instead, need to test more cases
		// check if it's annotation at.RequestMapping
		if ann := annotation.GetAnnotation(request, atType); ann != nil {
			// TODO: should confirm if value is usable
			h.requests[i].iVal = injectableMethod.annotations.Items[0].Parent.Value //ann.Parent.Value.Elem()
			h.requests[i].isAnnotation = true

			// operation
			_ = annotation.InjectAll(request)
			h.annotations = request
			continue
		}

		// parse path variable
		if pathIdx < lenOfPathParams {
			for idx, pv := range pps {
				if pv == pp[pathIdx][0] {
					h.requests[i].name = pp[pathIdx][1]
					h.requests[i].pathIdx = idx
					pathIdx = pathIdx + 1
					break
				}
			}
		}

		// request params, body, form
		if iTyp.Kind() == reflect.Struct {
			if typ.Kind() == reflect.Slice {
				//h.requests[i].obj = h.requests[i].iVal.Interface()
				h.requests[i].iVal = reflect.MakeSlice(typ, 0, 0)
				h.requests[i].typeName = "RequestBody"
				h.requests[i].callback = RequestBody
			} else {
				for _, tn := range requestSets {
					if field, ok := iTyp.FieldByName(tn.name); ok && field.Anonymous {
						h.requests[i].typeName = tn.name
						h.requests[i].callback = tn.callback
						break
					}
				}
			}
		}

	}
	h.lenOfPathParams = lenOfPathParams

	// parse response
	h.responses = make([]response, h.numOut)
	for i := 0; i < h.numOut; i++ {
		typ := method.Type.Out(i)

		h.responses[i].typ = typ
		h.responses[i].kind = typ.Kind()
		h.responses[i].typeName = typ.Name()
		//log.Debug(h.responses[i])
		if annotation.Contains(typ, at.ResponseBody{}) {
			h.responses[i].isResponseBody = true
		}
		modelType := reflect.TypeOf(new(model.Response)).Elem()
		h.responses[i].implementsResponse = typ.Implements(modelType)
		modelType = reflect.TypeOf(new(model.ResponseInfo)).Elem()
		h.responses[i].implementsResponseInfo = typ.Implements(modelType)
	}

	// check if configured annotation for starters
	//for _, subscriber := range h.factory.GetInstances() {
	//
	//}

	// finally, print mapped method
	if path != "" {
		log.Infof("Mapped %v \"%v\" onto %v.%v()", injectableMethod.requestMapping.Method, path, idv.Type(), method.Name)
	}
}

func (h *handler) responseData(ctx context.Context, numOut int, results []reflect.Value) (err error) {
	var res model.ResponseInfo

	for idx, val := range results {
		var objVal interface{}
		if val.CanInterface() {
			objVal = val.Interface()
		} else {
			err = ErrCanNotInterface
			return
		}

		switch objVal.(type) {
		case model.ResponseInfo:
			res = objVal.(model.ResponseInfo)
		case model.Response:
			res = objVal.(model.Response)
		case error:
			respErr := objVal.(error)
			if res != nil {
				if respErr == nil {
					if res.GetCode() == 0 {
						res.SetCode(http.StatusOK)
					}
					if res.GetMessage() == "" {
						res.SetMessage(ctx.Translate(success))
					}
				} else {
					h.setErrorResponseCode(ctx, res)
					// TODO: output error message directly? how about i18n
					res.SetMessage(ctx.Translate(respErr.Error()))
				}
			}
		}
		ctx.SetResponse(idx, objVal)
	}
	return
}

func (h *handler) finalizeResponse(ctx context.Context) {

	idx, size := ctx.HandlerIndex(-1), len(ctx.Handlers())-1

	if idx == size {

		// TODO: find out the root cause of strings type being called twice
		done := ctx.Values().GetString("done") == "true"
		if !done {
			//log.Debugf("%v / %v", idx, size)
			ctx.Values().Set("done", "true")
			res := ctx.GetResponse(0)
			if res != nil {
				switch res.(type) {
				case string:
					ctx.ResponseString(res.(string))
				case error:
					respErr := res.(error)
					ctx.ResponseError(respErr.Error(), http.StatusInternalServerError)
				case model.Response:
					r := res.(model.Response)
					if r.GetCode() == 0 {
						r.SetCode(http.StatusOK)
					}
					if r.GetMessage() == "" {
						r.SetMessage(ctx.Translate(success))
					}
					ctx.StatusCode(r.GetCode())
					_, _ = ctx.JSON(r)
				case model.ResponseInfo:
					r := res.(model.ResponseInfo)
					if r.GetCode() == 0 {
						r.SetCode(http.StatusOK)
					}
					if r.GetMessage() == "" {
						r.SetMessage(ctx.Translate(success))
					}
					ctx.StatusCode(r.GetCode())
					_, _ = ctx.JSON(r)
				case map[string]interface{}:
					_, _ = ctx.JSON(res)
				default:
					_, _ = ctx.JSON(res)
				}
			} else {
				e := ctx.GetResponse(1)
				switch e.(type) {
				case error:
					err := e.(error)
					ctx.ResponseError(err.Error(), http.StatusInternalServerError)
				}
			}
		}
	}
	return
}

func (h *handler) responseWithError(ctx context.Context, numOut int, results []reflect.Value) {
	if numOut >= 2 {
		responseObj := results[0].Interface()
		res := responseObj.(model.ResponseInfo)
		var respErr error
		//errVal := results[1]
		errObj := results[1].Interface()

		switch errObj.(type) {
		case error:
			respErr = errObj.(error)
		}

		if respErr == nil {
			if res.GetCode() == 0 {
				res.SetCode(http.StatusOK)
			}
			if res.GetMessage() == "" {
				res.SetMessage(ctx.Translate(success))
			}
		} else {
			h.setErrorResponseCode(ctx, res)
			// TODO: output error message directly? how about i18n
			res.SetMessage(ctx.Translate(respErr.Error()))

			// TODO: configurable status code in application.yml
		}
	}
}


func (h *handler) setErrorResponseCode(ctx context.Context, response model.ResponseInfo) {
	if response.GetCode() == 0 {
		prevStatusCode := ctx.GetStatusCode()
		if prevStatusCode == http.StatusOK || prevStatusCode == 0 {
			response.SetCode(http.StatusInternalServerError)
		} else {
			response.SetCode(prevStatusCode)
		}
	}
}

func (h *handler) call(ctx context.Context) {

	var input interface{}
	var reqErr error
	var path string
	var pvs []string
	var runtimeInstance factory.Instance

	// init responses
	log.Debugf("HTTP Handler: %v: %v %v%v", ctx.HandlerIndex(-1), ctx.Method(), ctx.Host(), ctx.Path())
	idx := ctx.HandlerIndex(-1)
	 if idx <= 1 {
		ctx.InitResponses()
	}

	if h.lenOfPathParams != 0 {
		path = ctx.Path()
		pvs = strings.SplitN(path, "/", -1)
	}
	var results []reflect.Value
	if len(h.dependencies) > 0 {
		runtimeInstance, reqErr = h.factory.InjectContextAwareObjects(ctx, h.dependencies)
		if reqErr != nil && h.numOut > 0 {
			results = make([]reflect.Value, h.numOut)
			if h.numOut > 1 {
				results[0] = reflect.New(h.responses[0].typ.Elem())
			}
			results[h.numOut - 1] = reflect.ValueOf(reqErr)
		}
	}
	if reqErr == nil {
		inputs := make([]reflect.Value, h.numIn)
		if h.numIn != 0 {
			inputs[0] = h.objVal
		}

		lenOfPathParams := h.lenOfPathParams
		for i := 1; i < h.numIn; i++ {
			req := h.requests[i]
			input = reflect.New(req.iTyp).Interface()

			// inject params
			//log.Debugf("%v, %v", i, h.requests[i].iVal.Type())
			if req.kind == reflect.Slice {
				var res reflect.Value
				if lenOfPathParams != 0 {
					res = h.decodePathVariable(&lenOfPathParams, pvs, req, req.kind)
				} else {
					res, reqErr = h.decodeSlice(ctx, h.requests[i].iTyp, h.requests[i].iVal)
				}
				inputs[i] = res
			} else if req.callback != nil {
				_ = h.factory.InjectDefaultValue(input) // support default value injection for request body/params/form
				reqErr = req.callback(ctx, input)
				inputs[i] = reflect.ValueOf(input)
			} else if req.kind == reflect.Interface && model.Context == req.typeName {
				input = ctx
				inputs[i] = reflect.ValueOf(input)
			} else if lenOfPathParams != 0 && !req.isAnnotation {
				// allow inject other dependencies after number of lenOfPathParams
				res := h.decodePathVariable(&lenOfPathParams, pvs, req, req.kind)

				inputs[i] = res
			} else {
				// inject instances
				var inst interface{}
				if runtimeInstance != nil {
					inst = runtimeInstance.Get(req.fullName)
				}
				if inst == nil {
					inst = h.factory.GetInstance(req.fullName) // TODO: primitive types does not need to get instance for the sake of performance
				}
				if inst != nil {
					inputs[i] = reflect.ValueOf(inst)
				} else {
					inputs[i] = h.requests[i].iVal
				}
			}
		}
		if reqErr == nil {
			// call controller method
			results = h.method.Func.Call(inputs)
		}
	}

	if h.numOut > 0 && len(results) > 0 {
		_ = h.responseData(ctx, h.numOut, results)
	}

	// finalize response
	h.finalizeResponse(ctx)
}

func (h *handler) decodePathVariable(lenOfPathParams *int, pvs []string, req request, kind reflect.Kind) reflect.Value {
	*lenOfPathParams = *lenOfPathParams - 1
	strVal := pvs[req.pathIdx]
	val := str.Convert(strVal, kind)
	res := reflect.ValueOf(val)
	return res
}

func (h *handler) decodeSlice(ctx context.Context, iTyp reflect.Type, input reflect.Value) (retVal reflect.Value, err error) {
	var m []interface{}
	err = ctx.ReadJSON(&m)
	for _, v := range m {
		item := reflect.New(iTyp).Interface()
		// TODO: Known issue - time.Time is not decoded
		err = mapstruct.Decode(item, v, mapstruct.WithSquash, mapstruct.WithWeaklyTypedInput)
		if err == nil {
			input = reflect.Append(input, reflect.ValueOf(item))
		}
	}
	retVal = input
	return
}
