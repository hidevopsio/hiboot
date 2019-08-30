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
	UrlSep			= "/"

)
type Dispatcher struct {
	webApp *webApp
	// inject context aware dependencies
	configurableFactory factory.ConfigurableFactory

	ContextPath string `value:"${server.context_path:/}"`
	ContextPathFormat string `value:"${server.context_path_format}" `
}

type requestMapping struct {
	Method     string
	Value      string
}

type restMethod struct {
	method *reflect.Method
	annotations []*annotation.Field
	hasMethodAnnotation bool
	requestMapping *requestMapping
}

type restController struct {
	controller interface{}
	fieldName string
	pkgPath string
	pathPrefix string
	before *restMethod
	after *restMethod
	methods []*restMethod
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


func (d *Dispatcher) parseRequestMapping(object interface{}, method *reflect.Method) (reqMap *requestMapping) {
	//reflector.
	reqMap = new(requestMapping)
	numIn := method.Type.NumIn()
	inputs := make([]reflect.Value, numIn)
	inputs[0] = reflect.ValueOf(object)
	for n := 1; n < numIn; n++ {
		typ := method.Type.In(n)
		o := reflect.New(typ).Interface()
		annotations := annotation.Find(o, at.HttpMethod{})
		if len(annotations) != 0 {
			err := d.configurableFactory.InjectIntoObject(o)
			if err == nil {
				_ = copier.Copy(reqMap, annotations[0].Value.Interface())
				break
			}
		}
	}
	return
}


func (d *Dispatcher) getAnnotations(object interface{}, method *reflect.Method) (annotations []*annotation.Field) {
	numIn := method.Type.NumIn()
	inputs := make([]reflect.Value, numIn)
	inputs[0] = reflect.ValueOf(object)
	for n := 1; n < numIn; n++ {
		typ := method.Type.In(n)
		log.Debugf("type: %v", typ)
		log.Debugf("type pkgPath: %v", typ.PkgPath())
		log.Debugf("type kind: %v", typ.Kind())
		if typ.Name() == "" && typ.Kind() == reflect.Struct {
			o := reflect.New(typ).Interface()
			annotations = annotation.GetFields(o)
			if len(annotations) != 0 {
				_ = d.configurableFactory.InjectIntoObject(o)
				break
			}
		}
	}
	return
}

func (d *Dispatcher) getRestMethods(metaData *factory.MetaData) (restCtl *restController) {

	restCtl = new(restController)

	c := metaData.Instance
	field := reflect.ValueOf(c)

	fieldType := field.Type()
	//log.Debug("fieldType: ", fieldType)
	ift := fieldType.Elem()
	fieldName := ift.Name()
	restCtl.pkgPath = ift.PkgPath()
	//log.Debug("fieldName: ", fieldName)

	controller := field.Interface()
	restCtl.controller = controller
	//log.Debug("controller: ", controller)

	// get context mapping
	var customizedControllerPath bool
	controllerPath := d.ContextPath
	af, ok := annotation.GetField(controller, at.RequestMapping{})
	if ok {
		customizedControllerPath = true
		controllerPath = filepath.Join(controllerPath, af.StructField.Tag.Get("value"))
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
		controllerPath = fmt.Sprintf("%v/%v", contextPath, cn)
	}
	restCtl.pathPrefix = controllerPath
	restCtl.fieldName = fieldName

	numOfMethod := field.NumMethod()
	//log.Debug("methods: ", numOfMethod)

	// find before, after method
	before, ok := fieldType.MethodByName(beforeMethod)
	if ok {
		restMethod := new(restMethod)
		restMethod.method = &before
		restCtl.before = restMethod
	}

	after, ok := fieldType.MethodByName(afterMethod)
	if ok {
		restMethod := new(restMethod)
		restMethod.method = &after
		restCtl.after = restMethod
	}

	var methods []*restMethod
	for mi := 0; mi < numOfMethod; mi++ {
		restMethod := new(restMethod)

		method := fieldType.Method(mi)
		methodName := method.Name
		annotations := d.getAnnotations(controller, &method)
		restMethod.annotations = annotations
		restMethod.method = &method
		httpMethodAnnotation := annotation.Filter(annotations, at.HttpMethod{})
		restMethod.hasMethodAnnotation = len(httpMethodAnnotation) > 0

		reqMap := new(requestMapping)
		if restMethod.hasMethodAnnotation {
			// only one HttpMethod should be annotated
			_ = copier.Copy(reqMap, httpMethodAnnotation[0].Value.Interface())
		}

		beforeMethod := annotation.Filter(annotations, at.BeforeMethod{})
		if len(beforeMethod) > 0 {
			restMethod.method = &method
			restCtl.before = restMethod
			continue
		}
		afterMethod := annotation.Filter(annotations, at.AfterMethod{})
		if len(afterMethod) > 0 {
			restMethod.method = &method
			restCtl.after = restMethod
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
		methods = append(methods, restMethod)
	}
	restCtl.methods = methods
	return
}

//TODO: scan apis and params to generate swagger api automatically by include swagger starter
func (d *Dispatcher) register(controllers []*factory.MetaData) (err error) {
	log.Debug("register rest controller")
	for _, metaData := range controllers {
		// get and parse all controller methods
		restController := d.getRestMethods(metaData)

		var party iris.Party
		if restController.before != nil {
			hdl := newHandler(d.configurableFactory)
			hdl.parse("", restController.before.method, restController.controller, "")
			party = d.webApp.Party(restController.pathPrefix, Handler(func(c context.Context) {
				hdl.call(c)
			}))
		} else {
			party = d.webApp.Party(restController.pathPrefix)
		}

		if restController.after != nil {
			hdl := newHandler(d.configurableFactory)
			hdl.parse("", restController.after.method, restController.controller, "")
			party.Done(Handler(func(c context.Context) {
				hdl.call(c)
			}))
		}

		// bind method handlers with router
		for _, m := range restController.methods{
			method := m.method
			methodName := method.Name
			//log.Debug("method: ", methodName)

			// check if it's valid http method
			hasAnyMethod := m.requestMapping.Method == Any
			hasRegularMethod := str.InSlice(m.requestMapping.Method, httpMethods)
			foundMethod := hasAnyMethod || hasRegularMethod
			if foundMethod {
				// parse all necessary requests and responses
				// create new method parser here
				hdl := newHandler(d.configurableFactory)
				hdl.parse(m.requestMapping.Method, m.method, restController.controller, restController.pathPrefix + m.requestMapping.Value)
				hdlHttpMethod := Handler(func(c context.Context) {
					hdl.call(c)
					c.Next()
				})

				if hasAnyMethod {
					party.Any(m.requestMapping.Value, hdlHttpMethod)
				} else if hasRegularMethod {
					route := party.Handle(m.requestMapping.Method, m.requestMapping.Value, hdlHttpMethod)
					route.MainHandlerName = fmt.Sprintf("%s/%s.%s", restController.pkgPath, restController.fieldName, methodName)
				}
			}
		}
	}
	return nil
}
