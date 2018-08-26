package autoconfigure

import (
	"os"
	"path/filepath"
	"reflect"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"github.com/hidevopsio/hiboot/pkg/utils/replacer"
	"github.com/hidevopsio/hiboot/pkg/system"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
	"github.com/hidevopsio/hiboot/pkg/utils/gotest"
	"github.com/hidevopsio/hiboot/pkg/inject"
	"github.com/hidevopsio/hiboot/pkg/factory/inst"
)

const (
	System            = "system"
	application       = "application"
	config            = "config"
	yaml              = "yaml"
	appProfilesActive = "APP_PROFILES_ACTIVE"
)

type ConfigurableFactory struct {
	*inst.InstanceFactory
	configurations cmap.ConcurrentMap
	systemConfig   *system.Configuration
	builder        *system.Builder
}

func (f *ConfigurableFactory) Initialize(configurations cmap.ConcurrentMap)  {
	if f.InstanceFactory == nil {
		log.Fatal("[factory] InstanceFactory can not be nil")
	}
	if !f.Initialized() {
		log.Fatal("[factory] instances map can not be nil")
	}
	f.configurations = configurations
	f.SetInstance("configurations", configurations)
}

func (f *ConfigurableFactory) SystemConfig() *system.Configuration {
	return f.systemConfig
}

func (f *ConfigurableFactory) BuildSystemConfig(configType interface{}) (err error) {
	workDir := io.GetWorkDir()

	profile := os.Getenv(appProfilesActive)
	f.builder = &system.Builder{
		Path:       filepath.Join(workDir, config),
		Name:       application,
		FileType:   yaml,
		Profile:    profile,
		ConfigType: configType,
	}

	systemConfig, err := f.builder.Build()

	if err == nil {
		f.systemConfig = systemConfig.(*system.Configuration)
	} else {
		f.systemConfig = new(system.Configuration)
	}
	// TODO: should separate instance to system and app
	f.SetInstance("systemConfiguration", f.systemConfig)
	inject.IntoObject(f.systemConfig)
	replacer.Replace(f.systemConfig, f.systemConfig)

	f.configurations.Set(System, f.systemConfig)

	return err
}

func (f *ConfigurableFactory) Build(configs ...cmap.ConcurrentMap) {
	for _, configMap := range configs {
		f.build(configMap)
	}
}

func (f *ConfigurableFactory) Instantiate(configuration interface{}) {
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
				f.SetInstance(instanceName, instance)
			}
		}
	}
}

func (f *ConfigurableFactory) build(cfgContainer cmap.ConcurrentMap)  {
	isTestRunning := gotest.IsRunning()
	for item := range cfgContainer.IterBuffered() {
		name, configType := item.Key, item.Val
		// TODO: should check if profiles is enabled str.InSlice(name, sysconf.App.Profiles.Include)
		if !isTestRunning && f.systemConfig != nil && !str.InSlice(name, f.systemConfig.App.Profiles.Include) {
			continue
		}
		log.Infof("auto configure: %v", name)

		// inject properties
		f.builder.ConfigType = configType
		f.builder.Profile = name
		cf, err := f.builder.Build()

		// TODO: check if cf.DependsOn

		if cf == nil {
			log.Warnf("failed to build %v configuration with error %v", name, err)
		} else {
			// replace references and environment variables
			if f.systemConfig != nil {
				replacer.Replace(cf, f.systemConfig)
			}
			inject.IntoObject(cf)
			replacer.Replace(cf, cf)

			// instantiation
			if err == nil {
				// create instances
				f.Instantiate(cf)
				// save configuration
				if _, ok := f.configurations.Get(name); ok {
					log.Fatalf("[factory] configuration name % is already taken", name)
				}
				f.configurations.Set(name, cf)
			}
		}
	}
}


