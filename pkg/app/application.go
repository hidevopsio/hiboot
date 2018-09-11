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

package app

import (
	"errors"
	"fmt"
	"github.com/hidevopsio/hiboot/pkg/factory/autoconfigure"
	"github.com/hidevopsio/hiboot/pkg/factory/instantiate"
	"github.com/hidevopsio/hiboot/pkg/inject"
	"github.com/hidevopsio/hiboot/pkg/system"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"github.com/kataras/iris/context"
	"reflect"
)

type Application interface {
	Initialize() error
	SetProperty(name string, value interface{}) Application
	Run() error
}

type ApplicationContext interface {
	RegisterController(controller interface{}) error
	Use(handlers ...context.Handler)
	GetProperty(name string) (value interface{}, ok bool)
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
}

var (
	configContainer    [][]interface{}
	componentContainer [][]interface{}

	InvalidObjectTypeError        = errors.New("[app] invalid Configuration type, one of app.Configuration, app.PreConfiguration, or app.PostConfiguration need to be embedded")
	ConfigurationNameIsTakenError = errors.New("[app] configuration name is already taken")
	ComponentNameIsTakenError     = errors.New("[app] component name is already taken")

	banner = `
______  ____________             _____
___  / / /__(_)__  /_______________  /_
__  /_/ /__  /__  __ \  __ \  __ \  __/   
_  __  / _  / _  /_/ / /_/ / /_/ / /_     Hiboot Application Framework
/_/ /_/  /_/  /_.___/\____/\____/\__/     https://github.com/hidevopsio/hiboot

`
)

func init() {
	//instanceContainer = cmap.New()
}

//
//func parseObjectName(eliminator string, inst interface{}) string  {
//	name := reflector.ParseObjectName(inst, eliminator)
//	if reflect.TypeOf(inst).Kind() == reflect.Func {
//
//	}
//
//	if name == "" || name == eliminator {
//		name = reflector.ParseObjectPkgName(inst)
//	}
//	return name
//}
//
//func parseInstance(eliminator string, params ...interface{}) (name string, inst interface{}) {
//
//	hasTwoParams := len(params) == 2 && reflect.TypeOf(params[0]).Kind() == reflect.String
//
//	if hasTwoParams {
//		inst = params[1]
//		name = params[0].(string)
//	} else {
//		inst = params[0]
//	}
//
//	if !hasTwoParams {
//		name = parseObjectName(eliminator, inst)
//	}
//
//	return
//}

func validateObjectType(inst interface{}) error {
	val := reflect.ValueOf(inst)
	//log.Println(val.Kind())
	//log.Println(reflect.Indirect(val).Kind())
	if val.Kind() == reflect.Ptr && reflect.Indirect(val).Kind() == reflect.Struct {
		return nil
	}
	return InvalidObjectTypeError
}

func hasTwoParams(params ...interface{}) bool {
	return len(params) == 2 && reflect.TypeOf(params[0]).Kind() == reflect.String
}

func appendParams(container [][]interface{}, params ...interface{}) (retVal [][]interface{}, err error) {
	retVal = container
	if len(params) == 0 || params[0] == nil {
		err = InvalidObjectTypeError
		return
	}

	item := make([]interface{}, 2)
	inst := params[0]
	if hasTwoParams(params...) {
		item[0] = params[0]
		item[1] = params[1]
		inst = params[1]
	} else {
		item[0] = params[0]
	}
	if inst != nil {
		kind := reflect.TypeOf(inst).Kind()
		if kind == reflect.Func || kind == reflect.Ptr {
			retVal = append(container, item)
			return
		}
	}
	err = InvalidObjectTypeError
	return
}

// AutoConfiguration register auto configuration struct
func AutoConfiguration(params ...interface{}) (err error) {
	configContainer, err = appendParams(configContainer, params...)
	return
}

// Component register a struct instance, so that it will be injectable.
// starter should register component type
func Component(params ...interface{}) (err error) {
	componentContainer, err = appendParams(componentContainer, params...)
	return
}

// BeforeInitialization ?
func (a *BaseApplication) PrintStartupMessages() {
	prop, ok := a.GetProperty(PropertyBannerDisabled)
	if !(ok && prop.(bool)) {
		fmt.Print(banner)
	}
}

// SetProperty
func (a *BaseApplication) SetProperty(name string, value interface{}) {
	a.propertyMap.Set(name, value)
}

// GetProperty
func (a *BaseApplication) GetProperty(name string) (value interface{}, ok bool) {
	value, ok = a.propertyMap.Get(name)
	return
}

// Init
func (a *BaseApplication) Initialize() error {
	a.WorkDir = io.GetWorkDir()

	a.propertyMap = cmap.New()

	a.configurations = cmap.New()
	a.instances = cmap.New()

	instantiateFactory := new(instantiate.InstantiateFactory)
	instantiateFactory.Initialize(a.instances)
	a.instances.Set("instantiateFactory", instantiateFactory)

	configurableFactory := new(autoconfigure.ConfigurableFactory)
	configurableFactory.InstantiateFactory = instantiateFactory
	a.instances.Set("configurableFactory", configurableFactory)
	inject.SetFactory(configurableFactory)
	a.configurableFactory = configurableFactory

	a.BeforeInitialization()

	err := configurableFactory.Initialize(a.configurations)
	if err != nil {
		return err
	}

	a.systemConfig = new(system.Configuration)
	configurableFactory.BuildSystemConfig(a.systemConfig)

	return nil
}

// Config returns application config
func (a *BaseApplication) SystemConfig() *system.Configuration {
	return a.systemConfig
}

func (a *BaseApplication) BuildConfigurations() {
	// build configurations
	a.configurableFactory.Build(configContainer)
	// build components
	a.configurableFactory.BuildComponents(componentContainer)
}

func (a *BaseApplication) ConfigurableFactory() *autoconfigure.ConfigurableFactory {
	return a.configurableFactory
}

func (a *BaseApplication) BeforeInitialization() {
	// pass user's instances
	a.postProcessor.BeforeInitialization(a.configurableFactory)
}

func (a *BaseApplication) AfterInitialization(configs ...cmap.ConcurrentMap) {
	// pass user's instances
	a.postProcessor.AfterInitialization(a.configurableFactory)
}

func (a *BaseApplication) RegisterController(controller interface{}) error {
	return nil
}

func (a *BaseApplication) Use(handlers ...context.Handler) {
}
