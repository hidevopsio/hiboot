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

// Package app provides abstract layer for cli/web application
package app

import (
	"errors"
	"fmt"
	"hidevops.io/hiboot/pkg/app/web/context"
	"hidevops.io/hiboot/pkg/factory"
	"hidevops.io/hiboot/pkg/factory/autoconfigure"
	"hidevops.io/hiboot/pkg/factory/instantiate"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/system"
	"hidevops.io/hiboot/pkg/utils/cmap"
	"hidevops.io/hiboot/pkg/utils/io"
	"os"
	"reflect"
	"strings"
	"sync"
)

const (
	// ApplicationContextName is the application context instance name
	ApplicationContextName = "app.applicationContext"
)

// Application is the base application interface
type Application interface {
	Initialize() error
	SetProperty(name string, value ...interface{}) Application
	GetProperty(name string) (value interface{}, ok bool)
	SetAddCommandLineProperties(enabled bool) Application
	Run()
}

// ApplicationContext is the alias interface of Application
type ApplicationContext interface {
	RegisterController(controller interface{}) error
	Use(handlers ...context.Handler)
	GetProperty(name string) (value interface{}, ok bool)
	GetInstance(params ...interface{}) (instance interface{})
}

// BaseApplication is the base application
type BaseApplication struct {
	WorkDir             string
	configurations      cmap.ConcurrentMap
	instances           cmap.ConcurrentMap
	potatoes            cmap.ConcurrentMap
	configurableFactory factory.ConfigurableFactory
	systemConfig        *system.Configuration
	postProcessor       *postProcessor
	properties          cmap.ConcurrentMap
	mu                  sync.Mutex
	// SetAddCommandLineProperties
	addCommandLineProperties bool
}

var (
	configContainer    []*factory.MetaData
	componentContainer []*factory.MetaData

	// ErrInvalidObjectType indicates that configuration type is invalid
	ErrInvalidObjectType = errors.New("[app] invalid Configuration type, one of app.Configuration need to be embedded")

	banner = `
______  ____________             _____
___  / / /__(_)__  /_______________  /_
__  /_/ /__  /__  __ \  __ \  __ \  __/   
_  __  / _  / _  /_/ / /_/ / /_/ / /_     Hiboot Application Framework
/_/ /_/  /_/  /_.___/\____/\____/\__/     https://hidevops.io/hiboot

`
)

// PrintStartupMessages prints startup messages
func (a *BaseApplication) PrintStartupMessages() {
	if !a.systemConfig.App.Banner.Disabled {
		fmt.Print(banner)
	}
}

// SetProperty set application property
// TODO: should set property from source by SetProperty or accept from program argument, e.g. myapp --app.profiles.active=dev
func (a *BaseApplication) SetProperty(name string, value ...interface{}) Application {
	var val interface{}
	if len(value) == 1 {
		val = value[0]
	} else {
		val = value
	}

	kind := reflect.TypeOf(val).Kind()
	if kind == reflect.String && strings.Contains(val.(string), ",") {
		val = strings.SplitN(val.(string), ",", -1)
	}
	a.properties.Set(name, val)

	return a
}

// GetProperty get application property
func (a *BaseApplication) GetProperty(name string) (value interface{}, ok bool) {
	value, ok = a.properties.Get(name)
	return
}

// Initialize init application
func (a *BaseApplication) Initialize() (err error) {
	log.SetLevel(log.InfoLevel)
	a.properties = cmap.New()
	a.configurations = cmap.New()
	a.instances = cmap.New()
	// set add command line properties to true as default
	a.SetAddCommandLineProperties(true)
	return nil
}

// Initialize init application
func (a *BaseApplication) Build() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.WorkDir = io.GetWorkDir()

	// set custom properties from args
	a.setCustomPropertiesFromArgs()

	instantiateFactory := instantiate.NewInstantiateFactory(a.instances, componentContainer, a.properties)
	// TODO: should set or get instance by passing object instantiateFactory
	instantiateFactory.SetInstance(factory.InstantiateFactoryName, instantiateFactory)
	instantiateFactory.AppendComponent(factory.InstantiateFactoryName, instantiateFactory)

	configurableFactory := autoconfigure.NewConfigurableFactory(instantiateFactory, a.configurations)
	instantiateFactory.SetInstance(factory.ConfigurableFactoryName, configurableFactory)
	instantiateFactory.AppendComponent(factory.ConfigurableFactoryName, configurableFactory)
	a.configurableFactory = configurableFactory

	a.postProcessor = newPostProcessor(instantiateFactory)

	a.systemConfig, _ = configurableFactory.BuildSystemConfig()

	// set logging level
	log.SetLevel(a.systemConfig.Logging.Level)
}

// SystemConfig returns application config
func (a *BaseApplication) setCustomPropertiesFromArgs() {
	//log.Println(os.Args)
	if a.addCommandLineProperties {
		for _, val := range os.Args {
			prefix := val[:2]
			if prefix == "--" {
				kv := val[2:]
				kvPair := strings.Split(kv, "=")
				// --property equal to --property=true
				if len(kvPair) == 1 {
					kvPair = append(kvPair, "true")
				}
				a.SetProperty(kvPair[0], kvPair[1])
			}
		}
	}
}

// SystemConfig returns application config
func (a *BaseApplication) SystemConfig() *system.Configuration {
	return a.systemConfig
}

// BuildConfigurations get BuildConfigurations
func (a *BaseApplication) BuildConfigurations() {
	// build configurations
	a.configurableFactory.Build(configContainer)
	// build components
	a.configurableFactory.BuildComponents()
}

// ConfigurableFactory get ConfigurableFactory
func (a *BaseApplication) ConfigurableFactory() factory.ConfigurableFactory {
	return a.configurableFactory
}

// AfterInitialization post initialization
func (a *BaseApplication) AfterInitialization(configs ...cmap.ConcurrentMap) {
	// pass user's instances
	a.postProcessor.Init()
	a.postProcessor.AfterInitialization()

	if a.addCommandLineProperties {
		log.Info("Add command line properties is enabled")
	} else {
		log.Info("Add command line properties is disabled")
	}
}

// RegisterController register controller by interface
func (a *BaseApplication) RegisterController(controller interface{}) error {
	return nil
}

// Use use middleware handlers
func (a *BaseApplication) Use(handlers ...context.Handler) {
}

// SetAddCommandLineProperties set add command line properties to be enabled or disabled
func (a *BaseApplication) SetAddCommandLineProperties(enabled bool) Application {
	a.addCommandLineProperties = enabled
	return a
}

// Run run the application
func (a *BaseApplication) Run() {
	log.Warn("application is not implemented!")
}

// GetInstance get application instance by name
func (a *BaseApplication) GetInstance(params ...interface{}) (instance interface{}) {
	if a.configurableFactory != nil {
		instance = a.configurableFactory.GetInstance(params...)
	}
	return
}
