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
	"strings"
	"regexp"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/middleware/i18n"
	"github.com/kataras/iris/middleware/logger"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/system"
	"github.com/hidevopsio/hiboot/pkg/starter"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
	"github.com/hidevopsio/hiboot/pkg/starter/grpc"
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
type Application interface {
	Init(controllers ...interface{}) error
	Config() *starter.SystemConfiguration
	Run()
}

// Application is the struct of web Application
// TODO: application should be singleton and private
type application struct {
	app             *iris.Application
	config          *starter.SystemConfiguration
	jwtEnabled      bool
	workDir         string
	httpMethods     []string
	factory         starter.Factory
	anonControllers []interface{}
	jwtControllers  []interface{}
	dispatcher      dispatcher
}

var (
	// controllers global controllers container
	webControllers []interface{}
	compiledRegExp = regexp.MustCompile(`\{(.*?)\}`)
)

// Config returns application config
func (wa *application) Config() *starter.SystemConfiguration {
	return wa.config
}

// Run run web application
func (wa *application) Run() {
	serverPort := ":8080"
	if wa.config != nil && wa.config.Server.Port != 0 {
		serverPort = fmt.Sprintf(":%v", wa.config.Server.Port)
	}
	// TODO: WithCharset should be configurable
	wa.app.Run(iris.Addr(fmt.Sprintf(serverPort)), iris.WithConfiguration(DefaultConfiguration()))
}

// Add add controller to controllers container
func Add(controllers ...interface{}) {
	webControllers = append(webControllers, controllers...)
}

func (wa *application) add(controllers ...interface{}) {
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
func (wa *application) Init(controllers ...interface{}) error {

	wa.workDir = io.GetWorkDir()

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

	wa.factory = starter.GetFactory()
	wa.factory.Build()

	config := wa.factory.Configuration(starter.System)
	if config != nil {
		//return errors.New("system configuration not found")
		wa.config = config.(*starter.SystemConfiguration)
		// ensure web is included
		if !str.InSlice("web", wa.config.App.Profiles.Include) {
			wa.config.App.Profiles.Include = append(wa.config.App.Profiles.Include, "web")
		}
		log.SetLevel(wa.config.Logging.Level)
	} else {
		wa.config = new(starter.SystemConfiguration)
		wa.config.App.Project = "hidevopsio"
		wa.config.App.Name = "hiboot"
		wa.config.App.Profiles.Include = append(wa.config.App.Profiles.Include, "web")
		log.Warnf("no config files in %v, e.g. application.yml", filepath.Join(wa.workDir, "config"))
	}

	if len(controllers) == 0 {
		wa.add(webControllers...)
	} else {
		wa.add(controllers...)
	}

	numJwtCtrl := len(wa.jwtControllers)
	// Init JWT
	err := InitJwt(wa.workDir)
	if err != nil {
		wa.jwtEnabled = false
		if numJwtCtrl != 0 {
			log.Warn(err.Error())
		}
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
		log.Debug(err)
	}

	// inject grpc services
	grpc.InjectIntoObject()

	if len(wa.anonControllers) == 0 &&
		len(wa.jwtControllers) == 0 {
		return nil
	}

	// first register anon controllers
	err = wa.dispatcher.register(wa.app, wa.anonControllers)
	if err != nil {
		return err
	}

	// then use jwt
	wa.app.Use(jwtHandler.Serve)

	// finally register jwt controllers
	err = wa.dispatcher.register(wa.app, wa.jwtControllers)
	if err != nil {
		return err
	}

	return nil
}

func (wa *application) initLocale() error {
	// TODO: localePath should be configurable in application.yml
	// locale:
	//   en-US: ./config/i18n/en-US.ini
	//   cn-ZH: ./config/i18n/cn-ZH.ini
	// TODO: or
	// locale:
	//   path: ./config/i18n/
	localePath := "config/i18n/"
	if io.IsPathNotExist(localePath) {
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
		lng = io.BaseDir(lng)
		lng = io.Basename(lng)

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

// NewApplication create new web application instance and init it
func NewApplication(controllers ...interface{}) *application {
	wa := new(application)
	err := wa.Init(controllers...)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	return wa
}
