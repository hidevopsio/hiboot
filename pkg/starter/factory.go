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

package starter

import (
	"sync"
	"github.com/hidevopsio/hiboot/pkg/system"
	"path/filepath"
	"os"
	"reflect"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
	"github.com/hidevopsio/hiboot/pkg/utils/replacer"
	"errors"
	"github.com/hidevopsio/hiboot/pkg/utils/gotest"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
)

type Factory interface {
	Build(configs ...interface{})
	Instantiate(configuration interface{})
	Configurations() cmap.ConcurrentMap
	Configuration(name string) interface{}
	Instances() cmap.ConcurrentMap
	Instance(name string) interface{}
	AddInstance(name string, instance interface{}) error
}

type factory struct {
	configurations cmap.ConcurrentMap
	instances      cmap.ConcurrentMap
}

const (
	System            = "system"
	application       = "application"
	config            = "config"
	yaml              = "yaml"
	appProfilesActive = "APP_PROFILES_ACTIVE"
)

var (
	bootFactory              *factory
	preConfigContainer       cmap.ConcurrentMap
	configContainer          cmap.ConcurrentMap
	postConfigContainer      cmap.ConcurrentMap
	once                     sync.Once
	InstanceNameIsTakenError = errors.New("[factory] instance name is already taken")
)

func init() {
	preConfigContainer = cmap.New()
	configContainer = cmap.New()
	postConfigContainer = cmap.New()
}

func GetFactory() Factory {
	once.Do(func() {
		bootFactory = new(factory)
		bootFactory.configurations = cmap.New()
		bootFactory.instances = cmap.New()
	})
	return bootFactory
}

func parseInstance(eliminator string, params ...interface{}) (name string, inst interface{}) {

	if len(params) == 2 && reflect.TypeOf(params[0]).Kind() == reflect.String {
		name = params[0].(string)
		inst = params[1]
	} else {
		name = reflector.ParseObjectName(params[0], eliminator)
		inst = params[0]
	}
	return
}

func addConfig(c cmap.ConcurrentMap, params ...interface{}) {

	name, inst := parseInstance("Configuration", params...)
	if name == "" && params != nil {
		name = reflector.ParseObjectPkgName(params[0])
	}

	if _, ok := c.Get(name); ok {
		log.Fatalf("configuration name %v is already taken!", name)
	}
	c.Set(name, inst)
}

func AddConfig(params ...interface{}) {
	addConfig(configContainer, params...)
}

func AddPreConfig(params ...interface{}) {
	addConfig(preConfigContainer, params...)
}

func AddPostConfig(params ...interface{}) {
	addConfig(postConfigContainer, params...)
}

func Add(params ...interface{}) {
	name, inst := parseInstance("Impl", params...)

	f := GetFactory()
	instances := f.Instances()

	if _, ok := instances.Get(name); ok {
		log.Fatalf("instance name %v is already taken!", name)
	}
	instances.Set(name, inst)
}

func (f *factory) build(cfgContainer cmap.ConcurrentMap, builder *system.Builder, sysconf *SystemConfiguration)  {
	isTestRunning := gotest.IsRunning()
	for item := range cfgContainer.IterBuffered() {
		name, configType := item.Key, item.Val
		// TODO: should check if profiles is enabled str.InSlice(name, sysconf.App.Profiles.Include)
		if !isTestRunning && sysconf != nil && !str.InSlice(name, sysconf.App.Profiles.Include) {
			continue
		}
		log.Infof("auto configure: %v", name)

		// inject properties
		builder.ConfigType = configType
		builder.Profile = name
		cf, err := builder.BuildWithProfile()

		// TODO: check if cf.DependsOn

		if cf == nil {
			log.Warnf("failed to build %v configuration with error %v", name, err)
		} else {
			// replace references and environment variables
			if sysconf != nil {
				replacer.Replace(cf, sysconf)
			}
			replacer.Replace(cf, cf)

			// instantiation
			if err == nil {
				// create instances
				f.Instantiate(cf)
				// save configuration
				f.configurations.Set(name, cf)
			}
		}
	}
}

func (f *factory) Build(configs ...interface{}) {

	workDir := io.GetWorkDir()
	profile := os.Getenv(appProfilesActive)
	builder := &system.Builder{
		Path:       filepath.Join(workDir, config),
		Name:       application,
		FileType:   yaml,
		Profile:    profile,
		ConfigType: SystemConfiguration{},
	}
	var sysconf *SystemConfiguration
	defaultConfig, err := builder.Build()
	if err == nil {
		f.configurations.Set(System, defaultConfig)
		replacer.Replace(defaultConfig, defaultConfig)
		sysconf = defaultConfig.(*SystemConfiguration)
		sysconf.App.Profiles.Active = profile
		log.Infof("profiles{active: %v, include: %v}", sysconf.App.Profiles.Active, sysconf.App.Profiles.Include)
	} else {
		log.Warnf("%v", err)
	}

	if len(configs) != 0 {
		name, inst := parseInstance("Configuration", configs...)
		configContainer.Set(name, inst)
	}

	f.build(preConfigContainer, builder, sysconf)
	f.build(configContainer, builder, sysconf)
	f.build(postConfigContainer, builder, sysconf)
}

func (f *factory) Instantiate(configuration interface{}) {
	cv := reflect.ValueOf(configuration)

	configType := cv.Type()
	//log.Debug("type: ", configType)
	//name := configType.Elem().Name()
	//log.Debug("fieldName: ", name)

	// call Init
	numOfMethod := cv.NumMethod()
	//log.Debug("methods: ", numOfMethod)

	for mi := 0; mi < numOfMethod; mi++ {
		method := configType.Method(mi)
		methodName := method.Name
		//log.Debugf("method: %v", methodName)
		numIn := method.Type.NumIn()
		// only 1 arg is supported so far
		if numIn == 1 {
			argv := make([]reflect.Value, numIn)
			argv[0] = reflect.ValueOf(configuration)
			retVal := method.Func.Call(argv)
			// save instance
			if retVal != nil && retVal[0].CanInterface() {
				instance := retVal[0].Interface()
				//log.Debugf("instantiated: %v", instance)
				instanceName := str.LowerFirst(methodName)
				if _, ok := f.instances.Get(instanceName); ok {
					log.Fatalf("method name %v is already taken!", instanceName)
				}
				f.instances.Set(instanceName, instance)
			}
		}
	}
}

func (f *factory) Configurations() cmap.ConcurrentMap {
	return f.configurations
}

func (f *factory) Configuration(name string) interface{} {
	retVal, ok := f.configurations.Get(name)
	if ok {
		return retVal
	}
	return nil
}

func (f *factory) Instances() cmap.ConcurrentMap {
	return f.instances
}

func (f *factory) Instance(name string) interface{} {
	retVal, ok := f.instances.Get(name)
	if ok {
		return retVal
	}
	return nil
}

func (f *factory) AddInstance(name string, instance interface{}) error {
	if _, ok := f.instances.Get(name); ok {
		return InstanceNameIsTakenError
	}
	f.instances.Set(name, instance)

	return nil
}
