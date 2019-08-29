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

	//contextAwareInstances []interface{}
	ContextPath string `value:"${server.context_path:/}"`
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

type requestMapping struct {
	customized bool
	Prefix     string
	Method     string
	Value      string
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
		annotations := annotation.Find(o, at.RequestMapping{})
		if len(annotations) != 0 {
			err := d.configurableFactory.InjectIntoObject(o)
			if err == nil {
				_ = copier.Copy(reqMap, annotations[0].Value.Interface())
				reqMap.customized = true
				break
			}
		}
	}
	return
}

//TODO: scan apis and params to generate swagger api automatically by include swagger starter
func (d *Dispatcher) register(controllers []*factory.MetaData) (err error) {
	for _, metaData := range controllers {
		reqMap := new(requestMapping)

		c := metaData.Instance
		field := reflect.ValueOf(c)

		fieldType := field.Type()
		//log.Debug("fieldType: ", fieldType)
		ift := fieldType.Elem()
		fieldName := ift.Name()
		pkgPath := ift.PkgPath()
		//log.Debug("fieldName: ", fieldName)

		controller := field.Interface()
		//log.Debug("controller: ", controller)

		//fieldValue := field.Elem()

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
			cpf := d.configurableFactory.GetProperty(app.ContextPathFormat)
			if cpf != nil {
				contextPathFormat := cpf.(string)

				switch contextPathFormat {
				case app.ContextPathFormatKebab:
					cn = str.ToKebab(controllerName)
				case app.ContextPathFormatSnake:
					cn = str.ToSnake(controllerName)
				case app.ContextPathFormatCamel:
					cn = str.ToCamel(controllerName)
				case app.ContextPathFormatLowerCamel:
					cn = str.ToLowerCamel(controllerName)
				}
			}
			contextPath := d.ContextPath
			if contextPath == ContextPathRoot {
				contextPath = ""
			}
			controllerPath = fmt.Sprintf("%v/%v", contextPath, cn)
		}

		numOfMethod := field.NumMethod()
		//log.Debug("methods: ", numOfMethod)

		beforeMethod, ok := fieldType.MethodByName(beforeMethod)
		var party iris.Party
		if ok {
			//log.Debug("contextPath: ", requestMapping)
			//log.Debug("beforeMethod.Name: ", beforeMethod.Name)
			hdl := newHandler(d.configurableFactory)
			hdl.parse("", beforeMethod, controller, "")
			party = d.webApp.Party(controllerPath, Handler(func(c context.Context) {
				hdl.call(c)
			}))
		} else {
			party = d.webApp.Party(controllerPath)
		}

		afterMethod, ok := fieldType.MethodByName(afterMethod)
		if ok {
			hdl := newHandler(d.configurableFactory)
			hdl.parse("", afterMethod, controller, "")
			party.Done(Handler(func(c context.Context) {
				hdl.call(c)
			}))
		}

		for mi := 0; mi < numOfMethod; mi++ {
			method := fieldType.Method(mi)
			methodName := method.Name
			//log.Debug("method: ", methodName)
			if methodName == "Options" {
				log.Debug("===")
			}

			reqMap = d.parseRequestMapping(controller, &method)
			if !reqMap.customized {
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
			if reqMap.Prefix == "" {
				reqMap.Prefix = controllerPath
			}

			// apirequestMapping should add arguments
			//log.Debug("requestMapping: ", apirequestMapping)
			// check if it's valid http method
			hasAnyMethod := reqMap.Method == Any
			hasGenericMethod := str.InSlice(reqMap.Method, httpMethods)
			foundMethod := hasAnyMethod || hasGenericMethod
			if foundMethod {
				// parse all necessary requests and responses
				// create new method parser here
				hdl := newHandler(d.configurableFactory)
				hdl.parse(reqMap.Method, method, controller, reqMap.Prefix+reqMap.Value)
				methodHandler := Handler(func(c context.Context) {
					hdl.call(c)
					c.Next()
				})

				if hasAnyMethod {
					party.Any(reqMap.Value, methodHandler)
				} else if hasGenericMethod {
					route := party.Handle(reqMap.Method, reqMap.Value, methodHandler)
					route.MainHandlerName = fmt.Sprintf("%s/%s.%s", pkgPath, fieldName, methodName)
				}
			}
		}
	}
	return nil
}
