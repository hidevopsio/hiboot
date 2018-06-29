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
	"os"
	"fmt"
	"reflect"
	"strings"
	"net/http"
	"crypto/rsa"
	"path/filepath"
	"github.com/fatih/camelcase"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/router"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/middleware/i18n"
	"github.com/hidevopsio/hiboot/pkg/system"
	"github.com/hidevopsio/hiboot/pkg/utils"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/starter/db"
	"github.com/kataras/iris/middleware/logger"
)

const (
	pathSep = "/"

	AuthTypeDefault = ""
	AuthTypeAnon    = "anon"
	AuthTypeJwt     = "jwt"

	initMethodName  = "Init"
)

type ApplicationInterface interface {
	Init()
	Config() *system.Configuration
	GetSignKey() *rsa.PrivateKey
	Run()
}

type DataSources map[string]interface{}

type Application struct {
	app    *iris.Application
	config *system.Configuration
	jwtEnabled bool
	workDir string
	httpMethods []string
	dataSources DataSources
	inject Inject
}

type Health struct {
	Status string `json:"status"`
}

const (
	application = "application"
	config      = "config"
	yaml        = "yaml"
)

var (
	Controllers []interface{}
)

func (wa *Application) Config() *system.Configuration {
	return wa.config
}

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
	inputs := make([]reflect.Value, method.Type.NumIn())

	inputs[0] = reflect.ValueOf(object)
	inputs[1] = reflect.ValueOf(ctx)
	method.Func.Call(inputs)
}

func Add(controller interface{})  {
	Controllers = append(Controllers, controller)
}

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

	builder := &system.Builder{
		Path:       filepath.Join(wa.workDir, config),
		Name:       application,
		FileType:   yaml,
		Profile:    os.Getenv("APP_PROFILES_ACTIVE"),
		ConfigType: system.Configuration{},
	}
	cp, err := builder.Build()
	if err == nil {
		wa.config = cp.(*system.Configuration)
		log.SetLevel(wa.config.Logging.Level)
	} else {
		log.SetLevel(log.DebugLevel)
	}

	// Init DataSource
	if wa.config != nil && len(wa.config.DataSources) != 0 {
		factory := db.DataSourceFactory{}
		dataSources := wa.config.DataSources
		wa.dataSources = make(DataSources)
		for _, dataSourceConfig := range dataSources {
			dataSourceType := dataSourceConfig["type"].(string)
			//log.Debug(dataSourceType)
			dataSource, err := factory.New(dataSourceType)
			if err == nil {
				err = dataSource.Open(dataSourceConfig)
				if err != nil {
					return err
				}
			}
			wa.dataSources[dataSourceType] = dataSource
		}
	}

	// Init JWT
	err = InitJwt(wa.workDir)
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
		controllers = Controllers
		if len(controllers) == 0 {
			return &system.NotFoundError{Name: "controller"}
		}
	}

	if ! wa.jwtEnabled {
		err := wa.register(controllers, AuthTypeAnon, AuthTypeDefault, AuthTypeJwt)
		if err != nil {
			return err
		}
	} else {
		err := wa.register(controllers, AuthTypeAnon, AuthTypeDefault)
		if err != nil {
			return err
		}

		wa.app.Use(jwtHandler.Serve)

		err = wa.register(controllers, AuthTypeJwt)
		if err != nil {
			return err
		}
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
		lng := strings.Replace(path, localePath,"", 1)
		lng = utils.BaseDir(lng)
		lng = utils.Basename(lng)

		if lng != "" && path != localePath + lng {
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

func (wa *Application) register(controllers []interface{}, auths... string) error {
	app := wa.app
	for _, c := range controllers {
		field := reflect.ValueOf(c)

		fieldType := field.Type()
		//log.Debug("fieldType: ", fieldType)
		fieldName := fieldType.Elem().Name()
		//log.Debug("fieldName: ", fieldName)

		controller := field.Interface()
		//log.Debug("controller: ", controller)

		// call Init
		initMethod, ok := fieldType.MethodByName(initMethodName)
		if ok {
			inputs := make([]reflect.Value, initMethod.Type.NumIn())
			inputs[0] = reflect.ValueOf(controller)
			initMethod.Func.Call(inputs)
		}

		fieldValue := field.Elem()
		fieldAuth := fieldValue.FieldByName("AuthType")
		if ! fieldAuth.IsValid() {
			return &system.InvalidControllerError{Name: fieldName}
		}
		a := fmt.Sprintf("%v", fieldAuth.Interface())
		//log.Debug("authType: ", a)
		if ! utils.StringInSlice(a, auths) {
			continue
		}

		// inject component
		err := wa.inject.IntoObject(field, wa.dataSources)
		if err != nil {
			return err
		}

		// get context mapping
		cp := fieldValue.FieldByName("ContextMapping")
		if ! cp.IsValid() {
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

		beforeMethod, ok := fieldType.MethodByName("Before")
		var party iris.Party
		if ok {
			//log.Debug("contextPath: ", contextMapping)
			//log.Debug("beforeMethod.Name: ", beforeMethod.Name)
			party = app.Party(contextMapping, func(ctx context.Context) {
				wa.handle(beforeMethod, controller, ctx.(*Context))
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
				if party == nil {
					relativePath := filepath.Join(contextMapping, apiContextMapping)
					//log.Debug("relativePath: ", relativePath)
					app.Handle(httpMethod, relativePath, func(ctx context.Context) {
						wa.handle(method, controller, ctx.(*Context))
					})
				} else {
					//log.Debug("contextMapping: ", apiContextMapping)
					party.Handle(httpMethod, apiContextMapping, func(ctx context.Context) {
						wa.handle(method, controller, ctx.(*Context))
					})
				}
			}
		}

	}
	return nil
}


func (wa *Application) createApplication(controllers []interface{}) error {
	err := wa.Init(controllers)
	log.Debugf("workDir: %v", wa.workDir)
	return err
}

func NewApplication(controllers ...interface{}) *Application {
	wa := new(Application)
	err := wa.Init(controllers...)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	return wa
}

