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
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"github.com/kataras/iris/context"
	"reflect"
)

type Application interface {
	Init() error
	SetProperty(name string, value interface{}) Application
	Run()
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
	configContainer     [][]interface{}
	instanceContainer   cmap.ConcurrentMap

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
	instanceContainer = cmap.New()
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

// AutoConfiguration register auto configuration struct
func AutoConfiguration(params ...interface{}) (err error) {
	cfg := make([]interface{}, 2)

	hasTwoParams := len(params) == 2 && reflect.TypeOf(params[0]).Kind() == reflect.String
	if hasTwoParams {
		cfg[0] = params[0]
		cfg[1] = params[1]
	} else {
		cfg[0] = params[0]
	}
	if cfg[0] != nil {
		kind := reflect.TypeOf(cfg[0]).Kind()
		if kind == reflect.Func || kind == reflect.Ptr{
			configContainer = append(configContainer, cfg)
			return nil
		}
	}

	return InvalidObjectTypeError
}

// Component register a struct instance, so that it will be injectable.
// starter should register component type
func Component(params ...interface{}) error {
	if len(params) == 0 || params[0] == nil {
		return InvalidObjectTypeError
	}
	var name string
	var inst interface{}
	if len(params) == 2 && reflect.TypeOf(params[0]).Kind() == reflect.String {
		name = params[0].(string)
		inst = params[1]
	} else {
		inst = params[0]
		ifcField := reflector.GetEmbeddedInterfaceField(inst)
		if ifcField.Anonymous {
			name = ifcField.Name
		}
	}
	if _, ok := instanceContainer.Get(name); ok {
		return ComponentNameIsTakenError
	}

	err := validateObjectType(inst)
	if err == nil {
		instanceContainer.Set(name, inst)
	}
	return err
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
func (a *BaseApplication) Init() error {
	a.WorkDir = io.GetWorkDir()

	a.propertyMap = cmap.New()

	a.configurations = cmap.New()
	a.instances = instanceContainer

	instanceFactory := new(instantiate.InstantiateFactory)
	instanceFactory.Initialize(a.instances)
	a.instances.Set("instanceFactory", instanceFactory)

	configurableFactory := new(autoconfigure.ConfigurableFactory)
	configurableFactory.InstantiateFactory = instanceFactory
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
	a.configurableFactory.Build(configContainer)
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
