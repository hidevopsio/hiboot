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
	"net"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"github.com/hidevopsio/iris"
	"github.com/hidevopsio/iris/core/router"
)

const (
	pathSep = "/"

	beforeMethod       = "Before"
	afterMethod        = "After"
	applicationContext = "app.applicationContext"
)

type webApp struct {
	*iris.Application
}

func newWebApplication() *webApp {
	return &webApp{
		Application: iris.New(),
	}
}

// Application is the struct of web Application
type application struct {
	app.BaseApplication
	webApp     *webApp
	jwtEnabled bool
	//anonControllers []interface{}
	//jwtControllers  []interface{}
	controllers []interface{}
	dispatcher  *Dispatcher
	//controllerMap map[string][]interface{}
	startUpTime time.Time
	// fallback serves any request the iris router has no route for; nil disables
	// the seam. Set via SetFallback and installed as a pre-routing WrapRouter
	// wrapper in buildHandler.
	fallback http.Handler
	// buildOnce memoizes buildHandler so one application can power both Run()
	// (serve on a listener) and Handler() (hand the http.Handler to a caller-
	// owned listener) without building the webApp twice. build() mutates global
	// component/profile state and is not safe to run twice on one instance.
	buildOnce    sync.Once
	builtHandler http.Handler
	builtErr     error
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
func (a *application) SetProperty(name string, value ...interface{}) app.Application {
	a.BaseApplication.SetProperty(name, value...)
	return a
}

// SetFallback installs an http.Handler that serves any request the iris router
// has no matching route for. Stored here and wired as a pre-routing WrapRouter
// wrapper in buildHandler so unmatched requests (legacy verbs, hijack/streaming)
// are delegated before iris touches the connection.
func (a *application) SetFallback(handler http.Handler) app.Application {
	a.fallback = handler
	return a
}

// Initialize init application
func (a *application) Initialize() error {
	return a.BaseApplication.Initialize()
}

// buildHandler is the shared build core behind both Run and Handler: it builds
// the HiBoot application and the underlying webApp (iris) instance, then returns
// it as an http.Handler. A missing-controllers condition is tolerated so a web
// app whose routes come purely from auto-configuration still builds; any other
// build error is returned. Keeping Run (serve) and Handler (return) on this one
// path guarantees the two modes build the webApp identically.
func (a *application) buildHandler() (http.Handler, error) {
	a.buildOnce.Do(func() {
		a.builtHandler, a.builtErr = a.doBuildHandler()
	})
	return a.builtHandler, a.builtErr
}

// doBuildHandler performs the one-time build. Guarded by buildHandler's
// buildOnce so build() (which mutates global component/profile state) runs at
// most once per application instance.
func (a *application) doBuildHandler() (http.Handler, error) {
	err := a.build()
	if err != nil && !errors.Is(err, ErrControllersNotFound) {
		return nil, err
	}
	a.webApp.Configure(iris.WithConfiguration(defaultConfiguration()))
	// Install the fallback seam before Build: iris applies router wrappers at
	// BuildRouter time. For each request, if iris owns no route for the
	// method+path, delegate to the fallback handler before iris touches the
	// connection — preserving hijack/streaming for legacy (unmigrated) verbs.
	if a.fallback != nil {
		fallback := a.fallback
		webApp := a.webApp
		webApp.WrapRouter(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			ctx := webApp.ContextPool.Acquire(w, r)
			exists := webApp.RouteExists(ctx, r.Method, r.URL.Path)
			// ReleaseLight returns the context to the pool without running
			// EndRequest, which would otherwise fire an error status code and
			// flush it to w for the (status 0) probe context.
			webApp.ContextPool.ReleaseLight(ctx)
			if exists {
				next(w, r)
				return
			}
			fallback.ServeHTTP(w, r)
		})
	}
	if err = a.webApp.Build(); err != nil {
		return nil, err
	}
	return a.webApp, nil
}

// Run web application
func (a *application) Run() {
	serverPort := ":8080"
	// build HiBoot Application and its webApp handler via the shared core
	handler, err := a.buildHandler()
	conf := a.SystemConfig()
	if conf == nil || err != nil {
		return
	}

	// Unix-socket transport: serve the handler directly on a caller-owned
	// socket path instead of binding a TCP port. Do not register on
	// http.DefaultServeMux — the socket listener serves the handler directly.
	if conf.Server.UnixSocket != "" {
		socketPath := conf.Server.UnixSocket
		// unlink a stale socket left by a previous run before listening
		if _, statErr := os.Stat(socketPath); statErr == nil {
			_ = os.Remove(socketPath)
		}
		listener, lerr := net.Listen("unix", socketPath)
		if lerr != nil {
			log.Errorf("Failed to listen on unix socket %v: %v", socketPath, lerr)
			return
		}
		// The daemon may run as root while the CLI runs as the invoking user;
		// loosen the socket perms so non-root client processes can connect.
		_ = os.Chmod(socketPath, 0o660)
		log.Infof("Hiboot started on unix socket %v", socketPath)
		timeDiff := time.Since(a.startUpTime)
		log.Infof("Started %v in %f seconds", conf.App.Name, timeDiff.Seconds())
		srv := &http.Server{Handler: handler}
		if conf.Server.TlsCert != "" && conf.Server.TlsKey != "" {
			log.Infof("Serving Hiboot web application with TLS on unix socket")
			err = srv.ServeTLS(listener, conf.Server.TlsCert, conf.Server.TlsKey)
		} else {
			log.Infof("Serving Hiboot web application on unix socket")
			err = srv.Serve(listener)
		}
		log.Debug(err)
		return
	}

	if conf.Server.Port != "" {
		serverPort = fmt.Sprintf(":%v", conf.Server.Port)
	}
	log.Infof("Hiboot started on port(s) http://localhost%v", serverPort)
	timeDiff := time.Since(a.startUpTime)
	log.Infof("Started %v in %f seconds", conf.App.Name, timeDiff.Seconds())

	// handler to Serve HTTP
	http.Handle("/", handler)

	// serve web app with server port, default port number is 8080
	if conf.Server.TlsCert != "" && conf.Server.TlsKey != "" {
		log.Infof("Serving Hiboot web application with TLS")
		err = http.ListenAndServeTLS(serverPort, conf.Server.TlsCert, conf.Server.TlsKey, nil)
	} else {
		log.Infof("Serving Hiboot web application")
		err = http.ListenAndServe(serverPort, nil)
		log.Debug(err)
	}
}

// Handler builds the web application and returns its http.Handler without
// taking over transport (unlike Run, which calls http.ListenAndServe). The
// caller owns the net.Listener, so the returned handler can be served on a
// Unix-socket listener or mounted on a foreign http.ServeMux.
func (a *application) Handler() (http.Handler, error) {
	return a.buildHandler()
}

// NewHandler creates a web application from the given controllers and returns
// its built http.Handler, for mounting on a caller-owned listener or mux. The
// web profile is included by default so the web auto-configuration (router,
// dispatcher) is not filtered out — mirroring SetProperty(app.ProfilesInclude,
// web.Profile) in a standard web.NewApplication(...).Run() bootstrap.
func NewHandler(controllers ...interface{}) (http.Handler, error) {
	a := NewApplication(controllers...)
	a.SetProperty(app.ProfilesInclude, Profile)
	wa, ok := a.(*application)
	if !ok {
		return nil, ErrInvalidController
	}
	return wa.Handler()
}

func unique(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// Init init web application
func (a *application) build() (err error) {

	a.Build()

	// set custom properties
	a.PrintStartupMessages()

	systemConfig := a.SystemConfig()
	// should do deduplication
	systemConfig.App.Profiles.Include = append(systemConfig.App.Profiles.Include, app.Profiles...)
	systemConfig.App.Profiles.Include = unique(systemConfig.App.Profiles.Include)
	if systemConfig != nil {
		log.Infof("Starting Hiboot web application %v version %v on localhost with PID %v", systemConfig.App.Name, systemConfig.App.Version, os.Getpid())
		log.Infof("Working directory: %v", a.WorkDir)
		log.Infof("The following profiles are active: %v, %v", systemConfig.App.Profiles.Active, systemConfig.App.Profiles.Include)
	}
	log.Infof("Initializing Hiboot Application")
	f := a.ConfigurableFactory()
	f.AppendComponent(app.ApplicationContextName, a)
	f.SetInstance(app.ApplicationContextName, a)

	// fill controllers into component container
	for _, ctrl := range a.controllers {
		f.AppendComponent(ctrl)
	}

	// build auto configurations
	err = a.BuildConfigurations()
	if err != nil {
		return
	}

	// create dispatcher
	a.dispatcher = a.GetInstance(Dispatcher{}).(*Dispatcher)

	// first register anon controllers
	err = a.RegisterController(at.RestController{})
	if err == nil {
		// call AfterInitialization with factory interface
		a.AfterInitialization()
	}
	return
}

// RegisterController register controller, e.g. at.RestController, jwt.Controller, or other customized controller
func (a *application) RegisterController(controller interface{}) error {
	middleware := a.ConfigurableFactory().GetInstances(at.Middleware{})
	log.Debug(middleware)
	// get from controller map
	// parse controller type
	controllers := a.ConfigurableFactory().GetInstances(controller)
	if controllers != nil {
		return a.dispatcher.register(controllers, middleware)
	}
	return ErrControllersNotFound
}

// Use apply middleware
func (a *application) Use(handlers ...context.Handler) {
	// pass user's instances
	for _, hdl := range handlers {
		a.webApp.Use(Handler(hdl))
	}
}

func (a *application) WrapRouter(handler router.WrapperFunc) {
	a.webApp.WrapRouter(handler)
}

func (a *application) initialize(controllers ...interface{}) (err error) {
	io.EnsureWorkDir(3, "config/application.yml")

	// new iris app
	a.webApp = newWebApplication()
	app.Register(a.webApp)

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

// SetAddCommandLineProperties set add command line properties to be enabled or disabled
func (a *application) SetAddCommandLineProperties(enabled bool) app.Application {
	a.BaseApplication.SetAddCommandLineProperties(enabled)
	return a
}

// RestController register rest controller to controllers container
// Deprecated: please use app.Register() instead
var RestController = app.Register

// NewApplication create new web application instance and init it
func NewApplication(controllers ...interface{}) app.Application {
	log.SetLevel("error") // set debug level to error first
	a := new(application)
	app.Register(a)
	a.startUpTime = time.Now()
	_ = a.initialize(controllers...)
	return a
}
