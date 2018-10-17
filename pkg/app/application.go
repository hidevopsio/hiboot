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
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/factory/autoconfigure"
	"github.com/hidevopsio/hiboot/pkg/factory/instantiate"
	"github.com/hidevopsio/hiboot/pkg/inject"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/system"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"github.com/kataras/iris/context"
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
	SetProperty(name string, value interface{}) Application
	GetProperty(name string) (value interface{}, ok bool)
	Run() error
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
	postProcessor       postProcessor
	propertyMap         cmap.ConcurrentMap
	mu                  sync.Mutex
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
/_/ /_/  /_/  /_.___/\____/\____/\__/     https://github.com/hidevopsio/hiboot

`
)

// PrintStartupMessages prints startup messages
func (a *BaseApplication) PrintStartupMessages() {
	prop, ok := a.GetProperty(PropertyBannerDisabled)
	if !(ok && prop.(bool)) {
		fmt.Print(banner)
	}
}

// SetProperty set application property
func (a *BaseApplication) SetProperty(name string, value interface{}) Application {
	a.propertyMap.Set(name, value)
	return a
}

// GetProperty get application property
func (a *BaseApplication) GetProperty(name string) (value interface{}, ok bool) {
	value, ok = a.propertyMap.Get(name)
	return
}

// Initialize init application
func (a *BaseApplication) Initialize() (err error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.WorkDir = io.GetWorkDir()

	a.propertyMap = cmap.New()

	a.configurations = cmap.New()
	a.instances = cmap.New()

	instantiateFactory := instantiate.NewInstantiateFactory(a.instances, componentContainer)
	// TODO: should set or get instance by passing object instantiateFactory
	a.instances.Set(factory.InstantiateFactoryName, instantiateFactory)

	configurableFactory := autoconfigure.NewConfigurableFactory(instantiateFactory, a.configurations)
	a.instances.Set(factory.ConfigurableFactoryName, configurableFactory)
	inject.SetFactory(configurableFactory)
	a.configurableFactory = configurableFactory

	a.systemConfig, _ = configurableFactory.BuildSystemConfig()
	return nil
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
	a.postProcessor.AfterInitialization(a.configurableFactory)
}

// RegisterController register controller by interface
func (a *BaseApplication) RegisterController(controller interface{}) error {
	return nil
}

// Use use middleware handlers
func (a *BaseApplication) Use(handlers ...context.Handler) {
}

// Run run the application
func (a *BaseApplication) Run() error {
	log.Warn("application is not implemented!")
	return nil
}

// GetInstance get application instance by name
func (a *BaseApplication) GetInstance(params ...interface{}) (instance interface{}) {
	if a.configurableFactory != nil {
		instance = a.configurableFactory.GetInstance(params...)
	}
	return
}

// AppendProfiles Run run the application
func (a *BaseApplication) AppendProfiles(app Application) error {
	profiles, ok := app.GetProperty(PropertyAppProfilesInclude)
	if ok {
		appProfilesInclude := strings.SplitN(profiles.(string), ",", -1)
		if a.systemConfig != nil {
			a.systemConfig.App.Profiles.Include =
				append(a.systemConfig.App.Profiles.Include, appProfilesInclude...)
		}
	}
	return nil
}
