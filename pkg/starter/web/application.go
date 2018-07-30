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
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/fatih/camelcase"
	"github.com/hidevopsio/hiboot/pkg/inject"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/system"
	"github.com/hidevopsio/hiboot/pkg/utils"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/core/router"
	"github.com/kataras/iris/middleware/i18n"
	"github.com/kataras/iris/middleware/logger"
	"github.com/hidevopsio/hiboot/pkg/starter"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"github.com/hidevopsio/hiboot/pkg/model"
)

const (
	pathSep = "/"

	AuthType        = "AuthType"
	AuthTypeDefault = ""
	AuthTypeAnon    = "anon"
	AuthTypeJwt     = "jwt"

	initMethodName = "Init"

	BeforeMethod = "Before"
	AfterMethod  = "After"
)

// ApplicationInterface is the interface of web application
type ApplicationInterface interface {
	Init()
	Config() *starter.SystemConfiguration
	Run()
}

// Application is the struct of web Application
// TODO: application should be singleton and private
type Application struct {
	app               *iris.Application
	config            *starter.SystemConfiguration
	jwtEnabled        bool
	workDir           string
	httpMethods       []string
	autoConfiguration starter.AutoConfiguration
	anonControllers   []interface{}
	jwtControllers    []interface{}
}

// Health is the health check struct
type Health struct {
	Status string `json:"status"`
}

var (
	// controllers global controllers container
	webControllers []interface{}
)

// Config returns application config
func (wa *Application) Config() *starter.SystemConfiguration {
	return wa.config
}

// Run run web application
func (wa *Application) Run() {
	serverPort := ":8080"
	if wa.config != nil && wa.config.Server.Port != 0 {
		serverPort = fmt.Sprintf(":%v", wa.config.Server.Port)
	}
	// TODO: WithCharset should be configurable
	wa.app.Run(iris.Addr(fmt.Sprintf(serverPort)), iris.WithCharset("UTF-8"), iris.WithoutVersionChecker)
}

func healthHandler(app *iris.Application) *router.Route {
	return app.Get("/health", func(ctx context.Context) {
		health := Health{
			Status: "UP",
		}
		ctx.JSON(health)
	})
}

func (wa *Application) handle(method reflect.Method, object interface{}, ctx *Context) {
	//log.Debug("NumIn: ", method.Type.NumIn())
	numIn := method.Type.NumIn()
	numOut := method.Type.NumOut()
	inputs := make([]reflect.Value, numIn)
	inputs[0] = reflect.ValueOf(object)
	var reqErr error
	if numIn >= 2 {
		// switch input type
		// find if the input type contains specific identifier RequestBody,RequestParams,RequestForm
		// then create new request instance
		// after that, call ctx.RequestBody() parse request
		requestType := reflector.IndirectType(method.Type.In(1))
		reqVal := reflect.New(requestType)
		var request interface{}
		request = reqVal.Interface()

		// TODO: evaluate the performance with below two solution
		inputs[1] = reqVal
		if field, ok := requestType.FieldByName(model.RequestTypeForm); ok && field.Anonymous {
			reqErr = ctx.RequestForm(request)
		} else if field, ok := requestType.FieldByName(model.RequestTypeParams); ok && field.Anonymous {
			reqErr = ctx.RequestParams(request)
		} else if field, ok := requestType.FieldByName(model.RequestTypeBody); ok && field.Anonymous {
			reqErr = ctx.RequestBody(request)
		} else {
			// assume that ctx is presented if it does not find above requests
			inputs[1] = reflect.ValueOf(ctx)
		}

		//result, err := reflector.CallMethodByName(request, "RequestType")
		//if err == nil {
		//	rt := fmt.Sprintf("%v", result)
		//	switch rt {
		//	case model.RequestTypeForm:
		//		reqErr = ctx.RequestForm(request)
		//	case model.RequestTypeParams:
		//		reqErr = ctx.RequestParams(request)
		//	default:
		//		reqErr = ctx.RequestBody(request)
		//	}
		//	inputs[1] = reqVal
		//} else {
		//	inputs[1] = reflect.ValueOf(ctx)
		//}

		//fmt.Printf("\nMethod: %v\nKind: %v\nName: %v\n-----------", method.Name, requestType.Kind(), requestType.Name())
	}

	var respErr error
	var results []reflect.Value
	if reqErr == nil {
		reflector.SetFieldValue(object, "Ctx", ctx)
		results = method.Func.Call(inputs)
		if numOut == 1 {
			result := results[0]
			if result.Kind() == reflect.String && result.CanInterface() {
				ctx.ResponseString(result.Interface().(string))
			}
		} else if numOut >= 2 {
			if method.Type.Out(1).Name() == "error" {
				if results[1].Kind() == reflect.Struct {
					fmt.Println(results[1].Type(), " ", results[1].Type().Name())
				}
				if results[1].IsNil() {
					respErr = nil
				} else {
					respErr = results[1].Interface().(error)
				}
			}
			responseTypeName := method.Type.Out(0).Name()
			var response model.Response
			switch responseTypeName {
			case "Response":
				response = results[0].Interface().(model.Response)
				if respErr == nil {
					response.Code = http.StatusOK
					response.Message = "success"
				} else {
					if response.Code == 0 {
						response.Code = http.StatusInternalServerError
					}
					// TODO: output error message directly? how about i18n
					response.Message = respErr.Error()
				}
			}
			ctx.JSON(response)
		}
	}
}

// Add add controller to controllers container
func Add(controllers ...interface{}) {
	webControllers = append(webControllers, controllers...)
}

func (wa *Application) add(controllers ...interface{}) {
	for _, controller := range controllers {
		authType := AuthTypeAnon
		result, err := reflector.CallMethodByName(controller, "AuthType")
		if err == nil {
			authType = fmt.Sprintf("%v", result)
		}
		// separate jwt controllers from anon
		if authType == AuthTypeJwt {
			wa.jwtControllers = append(wa.jwtControllers, controller)
		} else {
			wa.anonControllers = append(wa.anonControllers, controller)
		}
	}
}

// Init init web application
func (wa *Application) Init(controllers ...interface{}) error {

	wa.workDir = utils.GetWorkDir()

	wa.httpMethods = []string{
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

	wa.autoConfiguration = starter.GetAutoConfiguration()
	wa.autoConfiguration.Build()

	config := wa.autoConfiguration.Configuration(starter.System)
	if config != nil {
		//return errors.New("system configuration not found")
		wa.config = config.(*starter.SystemConfiguration)
	} else {
		log.Warnf("no application config files in %v", filepath.Join(wa.workDir, "config"))
	}

	// Init JWT
	err := InitJwt(wa.workDir)
	if err != nil {
		wa.jwtEnabled = false
		log.Warn(err.Error())
	} else {
		wa.jwtEnabled = true
	}

	wa.app = iris.New()

	customLogger := logger.New(logger.Config{
		// Status displays status code
		Status: true,
		// IP displays request's remote address
		IP: true,
		// Method displays the http method
		Method: true,
		// Path displays the request path
		Path: true,
		// Query appends the url query to the Path.
		//Query: true,

		//Columns: true,

		// if !empty then its contents derives from `ctx.Values().Get("logger_message")
		// will be added to the logs.
		MessageContextKeys: []string{"logger_message"},

		// if !empty then its contents derives from `ctx.GetHeader("User-Agent")
		MessageHeaderKeys: []string{"User-Agent"},
	})

	wa.app.Use(customLogger)

	// The only one Required:
	// here is how you define how your own context will
	// be created and acquired from the iris' generic context pool.
	wa.app.ContextPool.Attach(func() context.Context {
		return &Context{
			// Optional Part 3:
			Context: context.NewContext(wa.app),
		}
	})

	err = wa.initLocale()
	if err != nil {
		log.Warn(err)
	}

	healthHandler(wa.app)
	if len(controllers) == 0 {
		wa.add(webControllers...)
	} else {
		wa.add(controllers...)
	}

	if len(wa.anonControllers) == 0 &&
		len(wa.jwtControllers) == 0 {
		return &system.NotFoundError{Name: "controller"}
	}

	// first register anon controllers
	err = wa.register(wa.anonControllers)
	if err != nil {
		return err
	}

	// then use jwt
	wa.app.Use(jwtHandler.Serve)

	// finally register jwt controllers
	err = wa.register(wa.jwtControllers)
	if err != nil {
		return err
	}

	return nil
}

func (wa *Application) initLocale() error {
	// TODO: localePath should be configurable in application.yml
	// locale:
	//   en-US: ./config/i18n/en-US.ini
	//   cn-ZH: ./config/i18n/cn-ZH.ini
	// TODO: or
	// locale:
	//   path: ./config/i18n/
	localePath := "config/i18n/"
	if utils.IsPathNotExist(localePath) {
		return &system.NotFoundError{Name: localePath}
	}

	// parse language files
	languages := make(map[string]string)
	err := filepath.Walk(localePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		//*files = append(*files, path)
		lng := strings.Replace(path, localePath, "", 1)
		lng = utils.BaseDir(lng)
		lng = utils.Basename(lng)

		if lng != "" && path != localePath+lng {
			//languages[lng] = path
			if languages[lng] == "" {
				languages[lng] = path
			} else {
				languages[lng] = languages[lng] + ", " + path
			}
			//log.Debugf("%v, %v", lng, languages[lng])
		}
		return nil
	})
	if err != nil {
		return err
	}

	globalLocale := i18n.New(i18n.Config{
		Default:      "en-US",
		URLParameter: "lang",
		Languages:    languages,
	})

	wa.app.Use(globalLocale)

	return nil
}

func (wa *Application) register(controllers []interface{}) error {
	app := wa.app
	for _, c := range controllers {
		field := reflect.ValueOf(c)

		fieldType := field.Type()
		//log.Debug("fieldType: ", fieldType)
		fieldName := fieldType.Elem().Name()
		//log.Debug("fieldName: ", fieldName)

		controller := field.Interface()
		//log.Debug("controller: ", controller)

		// inject component
		err := inject.IntoObject(field)
		if err != nil {
			return err
		}

		//// call Init
		//initMethod, ok := fieldType.MethodByName(initMethodName)
		//if ok {
		//	inputs := make([]reflect.Value, initMethod.Type.NumIn())
		//	inputs[0] = reflect.ValueOf(controller)
		//	initMethod.Func.Call(inputs)
		//}

		fieldValue := field.Elem()

		// get context mapping
		cp := fieldValue.FieldByName("ContextMapping")
		if !cp.IsValid() {
			return &system.InvalidControllerError{Name: fieldName}
		}
		contextMapping := fmt.Sprintf("%v", cp.Interface())

		// parse method
		fieldNames := camelcase.Split(fieldName)
		controllerName := ""
		if len(fieldNames) >= 2 {
			controllerName = strings.Replace(fieldName, fieldNames[len(fieldNames)-1], "", 1)
			controllerName = utils.LowerFirst(controllerName)
		}
		//log.Debug("controllerName: ", controllerName)
		// use controller's prefix as context mapping
		if contextMapping == "" {
			contextMapping = pathSep + controllerName
		}

		numOfMethod := field.NumMethod()
		//log.Debug("methods: ", numOfMethod)

		beforeMethod, ok := fieldType.MethodByName(BeforeMethod)
		var party iris.Party
		if ok {
			//log.Debug("contextPath: ", contextMapping)
			//log.Debug("beforeMethod.Name: ", beforeMethod.Name)
			party = app.Party(contextMapping, func(ctx context.Context) {
				wa.handle(beforeMethod, controller, ctx.(*Context))
			})
		} else {
			party = app.Party(contextMapping)
		}

		afterMethod, ok := fieldType.MethodByName(AfterMethod)
		if ok {
			party.Done(func(ctx context.Context) {
				wa.handle(afterMethod, controller, ctx.(*Context))
			})
		}

		for mi := 0; mi < numOfMethod; mi++ {
			method := fieldType.Method(mi)
			methodName := method.Name
			//log.Debug("method: ", methodName)

			ctxMap := camelcase.Split(methodName)
			httpMethod := strings.ToUpper(ctxMap[0])
			apiContextMapping := strings.Replace(methodName, ctxMap[0], "", 1)
			apiContextMapping = pathSep + utils.LowerFirst(apiContextMapping)

			// check if it's valid http method
			if utils.StringInSlice(httpMethod, wa.httpMethods) {

				//log.Debug("contextMapping: ", apiContextMapping)
				party.Handle(httpMethod, apiContextMapping, func(ctx context.Context) {
					wa.handle(method, controller, ctx.(*Context))
					ctx.Next()
				})

			}
		}
	}
	return nil
}

// NewApplication create new web application instance and init it
func NewApplication(controllers ...interface{}) *Application {
	wa := new(Application)
	err := wa.Init(controllers...)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	return wa
}
