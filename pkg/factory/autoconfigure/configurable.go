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

// Package autoconfigure implement ConfigurableFactory
package autoconfigure

import (
	"errors"
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/factory/instantiate"
	"github.com/hidevopsio/hiboot/pkg/inject"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/system"
	"github.com/hidevopsio/hiboot/pkg/system/types"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"github.com/hidevopsio/hiboot/pkg/utils/gotest"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"github.com/hidevopsio/hiboot/pkg/utils/replacer"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

const (
	// System configuration name
	System      = "system"
	application = "application"
	config      = "config"
	yaml        = "yaml"

	// EnvAppProfilesActive is the environment variable name APP_PROFILES_ACTIVE
	EnvAppProfilesActive = "APP_PROFILES_ACTIVE"

	// PostfixConfiguration is the Configuration postfix
	PostfixConfiguration = "Configuration"
)

var (
	// ErrInvalidMethod method is invalid
	ErrInvalidMethod = errors.New("[factory] method is invalid")

	// ErrFactoryCannotBeNil means that the InstantiateFactory can not be nil
	ErrFactoryCannotBeNil = errors.New("[factory] InstantiateFactory can not be nil")

	// ErrFactoryIsNotInitialized means that the InstantiateFactory is not initialized
	ErrFactoryIsNotInitialized = errors.New("[factory] InstantiateFactory is not initialized")

	// ErrInvalidObjectType means that the Configuration type is invalid, it should embeds app.PreConfiguration
	ErrInvalidObjectType = errors.New("[factory] invalid Configuration type, one of app.Configuration, app.PreConfiguration, or app.PostConfiguration need to be embedded")

	// ErrConfigurationNameIsTaken means that the configuration name is already taken
	ErrConfigurationNameIsTaken = errors.New("[factory] configuration name is already taken")

	// ErrComponentNameIsTaken means that the component name is already taken
	ErrComponentNameIsTaken = errors.New("[factory] component name is already taken")
)

type ConfigurableFactory struct {
	*instantiate.InstantiateFactory
	configurations cmap.ConcurrentMap
	systemConfig   *system.Configuration
	builder        *system.Builder

	preConfigureContainer  []*factory.MetaData
	configureContainer     []*factory.MetaData
	postConfigureContainer []*factory.MetaData
}

// Initialize initialize ConfigurableFactory
func (f *ConfigurableFactory) Initialize(configurations cmap.ConcurrentMap) (err error) {
	if f.InstantiateFactory == nil {
		return ErrFactoryCannotBeNil
	}
	if !f.Initialized() {
		return ErrFactoryIsNotInitialized
	}
	f.configurations = configurations
	f.SetInstance("configurations", configurations)

	return
}

// SystemConfiguration getter
func (f *ConfigurableFactory) SystemConfiguration() *system.Configuration {
	return f.systemConfig
}

// Configuration getter
func (f *ConfigurableFactory) Configuration(name string) interface{} {
	cfg, ok := f.configurations.Get(name)
	if ok {
		return cfg
	}
	return nil
}

// BuildSystemConfig build system configuration
func (f *ConfigurableFactory) BuildSystemConfig() (systemConfig *system.Configuration, err error) {
	workDir := io.GetWorkDir()
	systemConfig = new(system.Configuration)
	profile := os.Getenv(EnvAppProfilesActive)
	f.builder = &system.Builder{
		Path:       filepath.Join(workDir, config),
		Name:       application,
		FileType:   yaml,
		Profile:    profile,
		ConfigType: systemConfig,
	}

	f.SetInstance("systemConfiguration", systemConfig)
	inject.DefaultValue(systemConfig)
	_, err = f.builder.Build()
	if err == nil {
		// TODO: should separate instance to system and app
		inject.IntoObject(systemConfig)
		replacer.Replace(systemConfig, systemConfig)

		f.configurations.Set(System, systemConfig)

		f.systemConfig = systemConfig
	}
	return
}

// Build build all auto configurations
func (f *ConfigurableFactory) Build(configs []*factory.MetaData) {
	// categorize configurations first, then inject object if necessary
	for _, item := range configs {
		ifcField := reflector.GetEmbeddedInterfaceField(item.Object)

		if ifcField.Anonymous {
			switch ifcField.Name {
			case "Configuration":
				f.configureContainer = append(f.configureContainer, item)
			case "PreConfiguration":
				f.preConfigureContainer = append(f.preConfigureContainer, item)
			case "PostConfiguration":
				f.postConfigureContainer = append(f.postConfigureContainer, item)
			default:
				continue
			}
		} else {
			err := ErrInvalidObjectType
			log.Errorf("item: %v err: %v", item, err)
			continue
		}
	}

	f.build(f.preConfigureContainer)
	f.build(f.configureContainer)
	f.build(f.postConfigureContainer)
}

// Instantiate run instantiation by method
func (f *ConfigurableFactory) Instantiate(configuration interface{}) (err error) {
	cv := reflect.ValueOf(configuration)

	// inject configuration before instantiation
	configType := cv.Type()
	//log.Debug("type: ", configType)
	//name := configType.Elem().Name()
	//log.Debug("fieldName: ", name)

	// call Init
	numOfMethod := cv.NumMethod()
	//log.Debug("methods: ", numOfMethod)
	for mi := 0; mi < numOfMethod; mi++ {
		method := configType.Method(mi)
		// append inst to f.components
		f.AppendComponent(configuration, method)
	}
	return
}

// appProfilesActive getter
func (f *ConfigurableFactory) appProfilesActive() string {
	if f.systemConfig == nil {
		return os.Getenv(EnvAppProfilesActive)
	}
	return f.systemConfig.App.Profiles.Active
}

func (f *ConfigurableFactory) parseName(item *factory.MetaData) string {
	name := strings.Replace(item.TypeName, PostfixConfiguration, "", -1)
	name = str.ToLowerCamel(name)

	if name == "" || name == strings.ToLower(PostfixConfiguration) {
		name = item.PkgName
	}
	return name
}

// build
func (f *ConfigurableFactory) build(cfgContainer []*factory.MetaData) {

	isTestRunning := gotest.IsRunning()
	for _, item := range cfgContainer {
		name := f.parseName(item)
		config := item.Object

		// TODO: should check if profiles is enabled str.InSlice(name, sysconf.App.Profiles.Include)
		if !isTestRunning && f.systemConfig != nil && !str.InSlice(name, f.systemConfig.App.Profiles.Include) {
			continue
		}
		log.Infof("Auto configure %v starter", name)

		// inject into func
		if item.Kind == types.Func {
			config, _ = inject.IntoFunc(config)
		}

		// inject properties
		f.builder.ConfigType = config

		// inject default value
		inject.DefaultValue(config)

		// build properties, inject settings
		cf, _ := f.builder.Build(name, f.appProfilesActive())
		// No properties needs to build, use default config
		if cf == nil {
			confTyp := reflect.TypeOf(config)
			if confTyp != nil && confTyp.Kind() == reflect.Ptr {
				cf = config
			} else {
				log.Errorf("Unsupported type: %v", confTyp)
				continue
			}
		}

		// replace references and environment variables
		if f.systemConfig != nil {
			replacer.Replace(cf, f.systemConfig)
		}
		inject.IntoObject(cf)
		replacer.Replace(cf, cf)

		// instantiation
		f.Instantiate(cf)
		// save configuration
		if _, ok := f.configurations.Get(name); ok {
			log.Fatalf("[factory] configuration name %v is already taken", name)
		}
		f.configurations.Set(name, cf)

	}
	//}
}
