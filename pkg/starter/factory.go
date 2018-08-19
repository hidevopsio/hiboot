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
	"github.com/hidevopsio/hiboot/pkg/utils/gotest"
)

type Factory interface {
	Build()
	Instantiate(configuration interface{})
	Configurations() map[string]interface{}
	Configuration(name string) interface{}
	Instances() map[string]interface{}
	Instance(name string) interface{}
}

type factory struct {
	configurations map[string]interface{}
	instances map[string]interface{}
}

const (
	System      = "system"
	application = "application"
	config      = "config"
	yaml        = "yaml"
	appProfilesActive = "APP_PROFILES_ACTIVE"
)

var (
	bootFactory *factory
	container   map[string]interface{}
	once        sync.Once
)

func init() {
	container = make(map[string]interface{})
}

func GetFactory() Factory {
	once.Do(func() {
		bootFactory = new(factory)
		bootFactory.configurations = make(map[string]interface{})
		bootFactory.instances = make(map[string]interface{})
	})
	return bootFactory
}

func parseInstance(eliminator string, params ...interface{}) (name string, inst interface{})  {

	if len(params) == 2 && reflect.TypeOf(params[0]).Kind() == reflect.String {
		name = params[0].(string)
		inst = params[1]
	} else {
		inst = params[0]
		name = reflector.ParseObjectName(inst, eliminator)
	}
	return
}

func AddConfig(params ...interface{})  {

	name, inst := parseInstance("Configuration", params...)

	if container[name] != nil {
		log.Fatalf("configuration name %v is already taken!", name)
	}
	container[name] = inst
}

func Add(params ...interface{})  {

	name, inst := parseInstance("Impl", params...)

	f := GetFactory()
	instances := f.Instances()

	if instances[name] != nil {
		log.Fatalf("instance name %v is already taken!", name)
	}
	instances[name] = inst
}

func (c *factory) Build()  {

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
		c.configurations[System] = defaultConfig
		replacer.Replace(defaultConfig, defaultConfig)
		sysconf = defaultConfig.(*SystemConfiguration)
		sysconf.App.Profiles.Active = profile
		log.Infof("profiles{active: %v, include: %v}", sysconf.App.Profiles.Active, sysconf.App.Profiles.Include)
	} else {
		log.Warnf("%v", err)
	}

	isTestRunning := gotest.IsRunning()
	for name, configType := range container {
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
			replacer.Replace(cf, defaultConfig)
			replacer.Replace(cf, cf)

			// instantiation
			if err == nil {
				// create instances
				c.Instantiate(cf)
				// save configuration
				c.configurations[name] = cf
			}
		}
	}
}

func (c *factory) Instantiate(configuration interface{})  {
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
			if retVal[0].CanInterface() {
				instance := retVal[0].Interface()
				//log.Debugf("instantiated: %v", instance)
				instanceName := str.LowerFirst(methodName)
				if c.instances[instanceName] != nil {
					log.Fatalf("method name %v is already taken!", instanceName)
				}
				c.instances[instanceName] = instance
			}
		}
	}
}

func (c *factory) Configurations() map[string]interface{} {
	return c.configurations
}

func (c *factory) Configuration(name string) interface{} {
	return c.configurations[name]
}

func (c *factory) Instances() map[string]interface{} {
	return c.instances
}

func (c *factory) Instance(name string) interface{} {
	return c.instances[name]
}