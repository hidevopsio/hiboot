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
	"github.com/fatih/camelcase"
	"github.com/kataras/iris"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web/context"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/factory"
	"hidevops.io/hiboot/pkg/inject/annotation"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/utils/copier"
	"hidevops.io/hiboot/pkg/utils/str"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"
)

var httpMethods = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

const (
	Any             = "ANY"
	RequestMapping  = "RequestMapping"
	ContextPathRoot = "/"
	UrlSep          = "/"
)

type Dispatcher struct {
	webApp *webApp
	// inject context aware dependencies
	configurableFactory factory.ConfigurableFactory

	ContextPath       string `value:"${server.context_path:/}"`
	ContextPathFormat string `value:"${server.context_path_format}" `
}

type requestMapping struct {
	Method string
	Value  string
}

type injectableMethod struct {
	method              *reflect.Method
	annotation          *methodAnnotation
	handler             iris.Handler
	hasMethodAnnotation bool
	requestMapping       *requestMapping
}

type middlewareHandler struct {
	annotation *methodAnnotation
}

type injectableObject struct {
	object     interface{}
	name       string
	pkgPath    string
	pathPrefix string
	before     *injectableMethod
	after      *injectableMethod
	methods    []*injectableMethod
}

type methodAnnotation struct {
	index  int
	fields []*annotation.Field
	object interface{}
	value  reflect.Value
}

func newDispatcher(webApp *webApp, configurableFactory factory.ConfigurableFactory) *Dispatcher {
	d := &Dispatcher{
		webApp:              webApp,
		configurableFactory: configurableFactory,
	}
	return d
}

func init() {
	app.Register(newDispatcher)
}

func (d *Dispatcher) parseAnnotation(object interface{}, method *reflect.Method) (ma *methodAnnotation) {
	ma = new(methodAnnotation)
	numIn := method.Type.NumIn()
	inputs := make([]reflect.Value, numIn)
	inputs[0] = reflect.ValueOf(object)
	for n := 1; n < numIn; n++ {
		typ := method.Type.In(n)
		if typ.Name() == "" && typ.Kind() == reflect.Struct {
			ma.value = reflect.New(typ)
			ma.object = ma.value.Interface()
			ma.fields = annotation.GetFields(ma.object)
			ma.index = n
			if len(ma.fields) != 0 {
				_ = d.configurableFactory.InjectIntoObject(ma.object)
				log.Debug(ma.object)
				break
			}
		}
	}
	return
}

// TODO: controller method handler, middleware handler, and event handler, they all use the same way to achieve handler dispatch.
func (d *Dispatcher) parseMiddleware(m *factory.MetaData) (middleware *injectableObject) {
	middleware = new(injectableObject)

	mwi := reflect.ValueOf(m.Instance)
	middleware.object = mwi.Interface()
	mwType := mwi.Type()
	numOfMethod := mwi.NumMethod()
	for mi := 0; mi < numOfMethod; mi++ {
		methodHandler := new(injectableMethod)
		method := mwType.Method(mi)
		methodHandler.method = &method
		ma := d.parseAnnotation(mi, &method)
		methodHandler.annotation = ma
		mwh := annotation.Filter(ma.fields, at.MiddlewareHandler{})
		if len(mwh) > 0 {
			methodHandler.hasMethodAnnotation = true
			hdl := newHandler(d.configurableFactory, middleware, methodHandler, at.MiddlewareHandler{})
			methodHandler.handler = Handler(func(c context.Context) {
				hdl.call(c)
			})
			middleware.methods = append(middleware.methods, methodHandler)
		}
	}
	return
}

func (d *Dispatcher) parseRestController(ctl *factory.MetaData) (restController *injectableObject) {

	restController = new(injectableObject)

	c := ctl.Instance
	field := reflect.ValueOf(c)

	fieldType := field.Type()
	//log.Debug("fieldType: ", fieldType)
	ift := fieldType.Elem()
	fieldName := ift.Name()
	restController.pkgPath = ift.PkgPath()
	//log.Debug("fieldName: ", fieldName)

	controller := field.Interface()
	restController.object = controller
	//log.Debug("controller: ", controller)

	// get context mapping
	var customizedControllerPath bool
	pathPrefix := d.ContextPath
	af, ok := annotation.GetField(controller, at.RequestMapping{})
	if ok {
		customizedControllerPath = true
		pathPrefix = filepath.Join(pathPrefix, af.StructField.Tag.Get("value"))
	}

	// parse method
	fieldNames := camelcase.Split(fieldName)
	controllerName := ""
	if len(fieldNames) >= 2 {
		controllerName = strings.Replace(fieldName, fieldNames[len(fieldNames)-1], "", 1)
		controllerName = str.LowerFirst(controllerName)
	}
	//log.Debug("controllerName: ", controllerName)
	// use controller's prefix as context mapping
	if !customizedControllerPath {
		cn := controllerName
		switch d.ContextPathFormat {
		case app.ContextPathFormatKebab:
			cn = str.ToKebab(controllerName)
		case app.ContextPathFormatSnake:
			cn = str.ToSnake(controllerName)
		case app.ContextPathFormatCamel:
			cn = str.ToCamel(controllerName)
		case app.ContextPathFormatLowerCamel:
			cn = str.ToLowerCamel(controllerName)
		}
		contextPath := d.ContextPath
		if contextPath == ContextPathRoot {
			contextPath = ""
		}
		pathPrefix = fmt.Sprintf("%v/%v", contextPath, cn)
	}
	restController.pathPrefix = pathPrefix
	restController.name = fieldName

	numOfMethod := field.NumMethod()
	//log.Debug("methods: ", numOfMethod)

	// find before, after method
	before, ok := fieldType.MethodByName(beforeMethod)
	if ok {
		restMethod := new(injectableMethod)
		restMethod.method = &before
		restController.before = restMethod
	}

	after, ok := fieldType.MethodByName(afterMethod)
	if ok {
		restMethod := new(injectableMethod)
		restMethod.method = &after
		restController.after = restMethod
	}

	var methods []*injectableMethod
	for mi := 0; mi < numOfMethod; mi++ {
		restMethod := new(injectableMethod)

		method := fieldType.Method(mi)
		methodName := method.Name
		restMethod.annotation = d.parseAnnotation(controller, &method)
		restMethod.method = &method
		httpMethodAnnotation := annotation.Filter(restMethod.annotation.fields, at.HttpMethod{})
		restMethod.hasMethodAnnotation = len(httpMethodAnnotation) > 0

		reqMap := new(requestMapping)
		if restMethod.hasMethodAnnotation {
			// only one HttpMethod should be annotated
			_ = copier.Copy(reqMap, httpMethodAnnotation[0].Value.Interface())
		}

		beforeMethod := annotation.Filter(restMethod.annotation.fields, at.BeforeMethod{})
		if len(beforeMethod) > 0 {
			restMethod.method = &method
			restController.before = restMethod
			continue
		}
		afterMethod := annotation.Filter(restMethod.annotation.fields, at.AfterMethod{})
		if len(afterMethod) > 0 {
			restMethod.method = &method
			restController.after = restMethod
			continue
		}

		if !restMethod.hasMethodAnnotation {
			ctxMap := camelcase.Split(methodName)
			reqMap.Method = strings.ToUpper(ctxMap[0])
			var apiPath string
			if len(ctxMap) > 2 && ctxMap[1] == "By" {
				for _, pathParam := range ctxMap[2:] {
					lpp := strings.ToLower(pathParam)
					apiPath = apiPath + pathSep + lpp + pathSep + "{" + lpp + "}"
				}
			} else {
				apiPath = strings.Replace(methodName, ctxMap[0], "", 1)
				apiPath = pathSep + str.LowerFirst(apiPath)
			}
			reqMap.Value = apiPath
		}
		restMethod.requestMapping = reqMap

		hasAnyMethod := reqMap.Method == Any
		hasRegularMethod := str.InSlice(reqMap.Method, httpMethods)
		foundMethod := hasAnyMethod || hasRegularMethod
		if foundMethod {
			methods = append(methods, restMethod)
		}

	}
	restController.methods = methods
	return
}

//TODO: scan apis and params to generate swagger api automatically by include swagger starter
func (d *Dispatcher) register(controllers []*factory.MetaData, middleware []*factory.MetaData) (err error) {
	var mws []*injectableObject
	for _, m := range middleware {
		mw := d.parseMiddleware(m)
		if mw != nil {
			mws = append(mws, mw)
		}
	}

	log.Debug("register rest controller")
	for _, ctl := range controllers {
		// get and parse all controller methods
		restController := d.parseRestController(ctl)

		var party iris.Party
		if restController.before != nil {
			hdl := newHandler(d.configurableFactory, restController, restController.before, at.BeforeMethod{})
			party = d.webApp.Party(restController.pathPrefix, Handler(func(c context.Context) {
				hdl.call(c)
			}))
		} else {
			party = d.webApp.Party(restController.pathPrefix)
		}

		if restController.after != nil {
			hdl := newHandler(d.configurableFactory, restController, restController.after, at.AfterMethod{})
			party.Done(Handler(func(c context.Context) {
				hdl.call(c)
			}))
		}

		// bind regular http method handlers with router
		for _, m := range restController.methods {
			var handlers []iris.Handler

			// 1. pass all annotations to registered starter for further implementations, e.g. swagger

			// 2. get all handlers from middleware, then use middleware
			// first check if controller annotated to use at.Middleware
			// second check if method annotated to use at.MiddlewareHandler
			// handlers = append(handlers, middleware...)
			if len(mws) > 0 {
				for _, mw := range mws {
					for _, mth := range mw.methods {
						handlers = append(handlers, mth.handler)
					}
				}
			}

			// 3. create new handler for rest controller method
			hdl := newHandler(d.configurableFactory, restController, m, at.HttpMethod{})
			httpHandler := Handler(func(c context.Context) {
				hdl.call(c)
				c.Next()
			})
			handlers = append(handlers, httpHandler)

			// 4. finally, handle all method handlers
			d.handleControllerMethod(restController, m, party, handlers)
		}
	}
	return
}

func (d *Dispatcher) handleControllerMethod(restController *injectableObject, m *injectableMethod, party iris.Party, httpHandler []iris.Handler) {
	if m.requestMapping.Method == Any {
		party.Any(m.requestMapping.Value, httpHandler...)
	} else {
		route := party.Handle(m.requestMapping.Method, m.requestMapping.Value, httpHandler...)
		route.MainHandlerName = fmt.Sprintf("%s/%s.%s", restController.pkgPath, restController.name, m.method.Name)
	}
}
