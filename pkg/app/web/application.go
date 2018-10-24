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
	"fmt"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"os"
	"regexp"
)

const (
	pathSep = "/"

	beforeMethod       = "Before"
	afterMethod        = "After"
	applicationContext = "app.applicationContext"
)

// Application is the struct of web Application
type application struct {
	app.BaseApplication
	webApp     *iris.Application
	jwtEnabled bool
	//anonControllers []interface{}
	//jwtControllers  []interface{}
	controllers []interface{}
	dispatcher  dispatcher
	//controllerMap map[string][]interface{}
}

var (
	// controllers global controllers container
	registeredControllers []interface{}
	compiledRegExp        = regexp.MustCompile(`\{(.*?)\}`)

	// ErrControllersNotFound controller not found
	ErrControllersNotFound = errors.New("[app] controllers not found")

	// ErrInvalidController invalid controller
	ErrInvalidController = errors.New("[app] invalid controller")
)

// SetProperty set application property
func (a *application) SetProperty(name string, value interface{}) app.Application {
	a.BaseApplication.SetProperty(name, value)
	return a
}

// Initialize init application
func (a *application) Initialize() error {
	return a.BaseApplication.Initialize()
}

// Run run web application
func (a *application) Run() (err error) {
	serverPort := ":8080"
	conf := a.SystemConfig()
	if conf != nil && conf.Server.Port != "" {
		serverPort = fmt.Sprintf(":%v", conf.Server.Port)
	}

	err = a.build()
	if err == nil {
		err = a.webApp.Run(iris.Addr(fmt.Sprintf(serverPort)), iris.WithConfiguration(defaultConfiguration()))
	}

	return
}

// Init init web application
func (a *application) build() (err error) {
	a.PrintStartupMessages()

	a.AppendProfiles(a)
	systemConfig := a.SystemConfig()
	if systemConfig != nil {
		log.SetLevel(systemConfig.Logging.Level)
		log.Infof("Starting Hiboot web application %v on localhost with PID %v", systemConfig.App.Name, os.Getpid())
		log.Infof("Working directory: %v", a.WorkDir)
		log.Infof("The following profiles are active: %v, %v", systemConfig.App.Profiles.Active, systemConfig.App.Profiles.Include)
	}

	f := a.ConfigurableFactory()
	f.AppendComponent(app.ApplicationContextName, a)
	f.SetInstance(app.ApplicationContextName, a)

	// fill controllers into component container
	for _, ctrl := range a.controllers {
		f.AppendComponent(ctrl)
	}

	// build auto configurations
	a.BuildConfigurations()

	// The only one Required:
	// here is how you define how your own context will
	// be created and acquired from the iris' generic context pool.
	a.webApp.ContextPool.Attach(func() context.Context {
		return &Context{
			// Optional Part 3:
			Context: context.NewContext(a.webApp),
		}
	})

	// first register anon controllers
	a.RegisterController(new(at.RestController))

	// call AfterInitialization with factory interface
	a.AfterInitialization()
	return err
}

// RegisterController register controller, e.g. web.Controller, jwt.Controller, or other customized controller
func (a *application) RegisterController(controller interface{}) error {
	// get from controller map
	// parse controller type
	controllerInterfaceName := reflector.GetName(controller)
	controllers := a.ConfigurableFactory().GetInstances(controllerInterfaceName)
	if controllers != nil {
		return a.dispatcher.register(a.webApp, controllers)
	}
	return nil
}

// Use apply middleware
func (a *application) Use(handlers ...context.Handler) {
	// pass user's instances
	for _, hdl := range handlers {
		a.webApp.Use(hdl)
	}
}

func (a *application) initialize(controllers ...interface{}) (err error) {
	io.EnsureWorkDir(3, "config/application.yml")

	// new iris app
	a.webApp = iris.New()

	err = a.Initialize()

	if err == nil {
		if len(controllers) == 0 {
			a.controllers = registeredControllers
		} else {
			a.controllers = controllers
		}
	}
	return
}

// RestController register rest controller to controllers container
// Deprecated: please use app.Register() instead
var RestController = app.Register

// NewApplication create new web application instance and init it
func NewApplication(controllers ...interface{}) app.Application {
	log.SetLevel("error") // set debug level to error first
	a := new(application)
	err := a.initialize(controllers...)
	if err != nil {
		os.Exit(1)
	}
	return a
}
