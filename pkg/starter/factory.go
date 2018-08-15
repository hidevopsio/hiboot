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
	"github.com/hidevopsio/hiboot/pkg/utils"
	"github.com/hidevopsio/hiboot/pkg/system"
	"path/filepath"
	"os"
	"reflect"
	"github.com/hidevopsio/hiboot/pkg/log"
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

func Add(name string, conf interface{})  {
	if container[name] != nil {
		log.Fatalf("configuration name %v is already taken!", name)
	}
	container[name] = conf
}

func (c *factory) Build()  {

	workDir := utils.GetWorkDir()
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
		utils.Replace(defaultConfig, defaultConfig)
		sysconf = defaultConfig.(*SystemConfiguration)
		sysconf.App.Profiles.Active = profile
		log.Infof("profiles{active: %v, include: %v}", sysconf.App.Profiles.Active, sysconf.App.Profiles.Include)
	} else {
		log.Warn(err)
	}

	for name, configType := range container {
		// TODO: should check if profiles is enabled utils.StringInSlice(name, sysconf.App.Profiles.Include)
		//if sysconf != nil && !utils.StringInSlice(name, sysconf.App.Profiles.Include) {
		//	continue
		//}

		// inject properties
		builder.ConfigType = configType
		builder.Profile = name
		cf, err := builder.BuildWithProfile()

		// TODO: check if cf.DependsOn

		if cf == nil {
			log.Warnf("failed to build %v configuration with error %v", name, err)
		} else {
			// replace references and environment variables
			utils.Replace(cf, defaultConfig)
			utils.Replace(cf, cf)

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
				if c.instances[methodName] != nil {
					log.Fatalf("method name %v is already taken!", methodName)
				}
				c.instances[methodName] = instance
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