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

type Application interface {
	Initialize() error
	SetProperty(name string, value interface{}) Application
	GetProperty(name string) (value interface{}, ok bool)
	Run() error
}

type ApplicationContext interface {
	RegisterController(controller interface{}) error
	Use(handlers ...context.Handler)
	GetProperty(name string) (value interface{}, ok bool)
	GetInstance(name string) (instance interface{})
}

type Configuration interface{}
type PreConfiguration interface{}
type PostConfiguration interface{}

type BaseApplication struct {
	WorkDir             string
	configurations      cmap.ConcurrentMap
	instances           cmap.ConcurrentMap
	potatoes            cmap.ConcurrentMap
	configurableFactory *autoconfigure.ConfigurableFactory
	systemConfig        *system.Configuration
	postProcessor       postProcessor
	propertyMap         cmap.ConcurrentMap
	mu                  sync.Mutex
}

var (
	configContainer    []*factory.MetaData
	componentContainer []*factory.MetaData

	// ErrInvalidObjectType indicates that configuration type is invalid
	ErrInvalidObjectType = errors.New("[app] invalid Configuration type, one of app.Configuration, app.PreConfiguration, or app.PostConfiguration need to be embedded")

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
func (a *BaseApplication) Initialize() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.WorkDir = io.GetWorkDir()

	a.propertyMap = cmap.New()

	a.configurations = cmap.New()
	a.instances = cmap.New()

	instantiateFactory := new(instantiate.InstantiateFactory)
	instantiateFactory.Initialize(a.instances, componentContainer)
	a.instances.Set("instantiateFactory", instantiateFactory)

	configurableFactory := new(autoconfigure.ConfigurableFactory)
	configurableFactory.InstantiateFactory = instantiateFactory
	a.instances.Set("configurableFactory", configurableFactory)
	inject.SetFactory(configurableFactory)
	a.configurableFactory = configurableFactory

	a.BeforeInitialization()

	err := configurableFactory.Initialize(a.configurations)
	if err == nil {
		a.systemConfig, err = configurableFactory.BuildSystemConfig()
	}
	return nil
}

// Config returns application config
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
func (a *BaseApplication) ConfigurableFactory() *autoconfigure.ConfigurableFactory {
	return a.configurableFactory
}

// BeforeInitialization pre initialization
func (a *BaseApplication) BeforeInitialization() {
	// pass user's instances
	a.postProcessor.BeforeInitialization(a.configurableFactory)
}

// AfterInitialization post initialization
func (a *BaseApplication) AfterInitialization(configs ...cmap.ConcurrentMap) {
	// pass user's instances
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
func (a *BaseApplication) GetInstance(name string) (instance interface{}) {
	if a.configurableFactory != nil {
		instance = a.configurableFactory.GetInstance(name)
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
