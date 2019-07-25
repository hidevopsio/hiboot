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
	"hidevops.io/hiboot/pkg/factory"
	"hidevops.io/hiboot/pkg/utils/reflector"
	"hidevops.io/hiboot/pkg/utils/str"
	"net/http"
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

const Any = "ANY"

type Dispatcher struct {
	webApp *webApp
	// inject context aware dependencies
	configurableFactory factory.ConfigurableFactory

	//contextAwareInstances []interface{}
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

//TODO: scan apis and params to generate swagger api automatically by include swagger starter
func (d *Dispatcher) register(controllers []*factory.MetaData) (err error) {
	for _, metaData := range controllers {
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
		contextMapping, ok := reflector.FindEmbeddedFieldTag(controller, "ContextPath", "value")

		// parse method
		fieldNames := camelcase.Split(fieldName)
		controllerName := ""
		if len(fieldNames) >= 2 {
			controllerName = strings.Replace(fieldName, fieldNames[len(fieldNames)-1], "", 1)
			controllerName = str.LowerFirst(controllerName)
		}
		//log.Debug("controllerName: ", controllerName)
		// use controller's prefix as context mapping
		if contextMapping == "" {
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
			contextMapping = pathSep + cn
		}

		numOfMethod := field.NumMethod()
		//log.Debug("methods: ", numOfMethod)

		beforeMethod, ok := fieldType.MethodByName(beforeMethod)
		var party iris.Party
		if ok {
			//log.Debug("contextPath: ", contextMapping)
			//log.Debug("beforeMethod.Name: ", beforeMethod.Name)
			hdl := newHandler(d.configurableFactory)
			hdl.parse(beforeMethod, controller, "")
			party = d.webApp.Party(contextMapping, Handler(func(c context.Context) {
				hdl.call(c)
			}))
		} else {
			party = d.webApp.Party(contextMapping)
		}

		afterMethod, ok := fieldType.MethodByName(afterMethod)
		if ok {
			hdl := newHandler(d.configurableFactory)
			hdl.parse(afterMethod, controller, "")
			party.Done(Handler(func(c context.Context) {
				hdl.call(c)
			}))
		}

		for mi := 0; mi < numOfMethod; mi++ {
			method := fieldType.Method(mi)
			methodName := method.Name
			//log.Debug("method: ", methodName)

			ctxMap := camelcase.Split(methodName)
			httpMethod := strings.ToUpper(ctxMap[0])

			// apiContextMapping should add arguments
			//log.Debug("contextMapping: ", apiContextMapping)
			// check if it's valid http method
			hasAnyMethod := httpMethod == Any
			hasGenericMethod := str.InSlice(httpMethod, httpMethods)
			foundMethod := hasAnyMethod || hasGenericMethod
			if foundMethod {
				var apiContextMapping string
				if len(ctxMap) > 2 && ctxMap[1] == "By" {
					for _, pathParam := range ctxMap[2:] {
						lpp := strings.ToLower(pathParam)
						apiContextMapping = apiContextMapping + pathSep + lpp + pathSep + "{" + lpp + "}"
					}
				} else {
					apiContextMapping = strings.Replace(methodName, ctxMap[0], "", 1)
					apiContextMapping = pathSep + str.LowerFirst(apiContextMapping)
				}

				// parse all necessary requests and responses
				// create new method parser here
				hdl := newHandler(d.configurableFactory)
				hdl.parse(method, controller, contextMapping+apiContextMapping)
				methodHandler := Handler(func(c context.Context) {
					hdl.call(c)
					c.Next()
				})

				if hasAnyMethod {
					party.Any(apiContextMapping, methodHandler)
				} else if hasGenericMethod {
					route := party.Handle(httpMethod, apiContextMapping, methodHandler)
					route.MainHandlerName = fmt.Sprintf("%s/%s.%s", pkgPath, fieldName, methodName)
				}
			}
		}
	}
	return nil
}
