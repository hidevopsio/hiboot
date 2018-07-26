package starter

import (
	"github.com/hidevopsio/hiboot/pkg/system"
	"github.com/hidevopsio/hiboot/pkg/utils"
	"path/filepath"
	"os"
	"reflect"
	"github.com/hidevopsio/hiboot/pkg/log"
	"sync"
)

const (
	System      = "system"
	application = "application"
	config      = "config"
	yaml        = "yaml"
	appProfilesActive = "APP_PROFILES_ACTIVE"
)

type Profiles struct {
	Include []string `json:"include"`
	Active  string   `json:"active"`
}

type App struct {
	Project        string   `json:"project"`
	Name           string   `json:"name"`
	Profiles       Profiles `json:"profiles"`
	DataSourceType string   `json:"data_source_type"`
}

type Server struct {
	Port int32 `json:"port"`
}

type Logging struct {
	Level string `json:"level"`
}

type Env struct {
	Name  string
	Value string
}

type DataSource map[string]interface{}

type SystemConfiguration struct {
	App         App          `mapstructure:"app"`
	Server      Server       `mapstructure:"server"`
	Logging     Logging      `mapstructure:"logging"`
	DataSources []DataSource `mapstructure:"dataSources"`
}

type AutoConfiguration interface {
	Build()
	Instantiate(configuration interface{})
	Configurations() map[string]interface{}
	Configuration(name string) interface{}
	Instances() map[string]interface{}
	Instance(name string) interface{}
}

type autoConfiguration struct {
	configurations map[string]interface{}
	instances map[string]interface{}
}

var (
	configuration  *autoConfiguration
	configurations map[string]interface{}
	once sync.Once
)

func init() {
	configurations = make(map[string]interface{})
}

func GetAutoConfiguration() AutoConfiguration {
	once.Do(func() {
		configuration = new(autoConfiguration)
		configuration.configurations = make(map[string]interface{})
		configuration.instances = make(map[string]interface{})
	})
	return configuration
}

func Add(name string, conf interface{})  {
	configurations[name] = conf
}

func (c *autoConfiguration) Build()  {

	workDir := utils.GetWorkDir()

	builder := &system.Builder{
		Path:       filepath.Join(workDir, config),
		Name:       application,
		FileType:   yaml,
		Profile:    os.Getenv(appProfilesActive),
		ConfigType: SystemConfiguration{},
	}
	defaultConfig, err := builder.Build()
	if err == nil {
		c.configurations[System] = defaultConfig
	}
	utils.Replace(defaultConfig, defaultConfig)

	for name, configType := range configurations {
		// inject properties
		builder.ConfigType = configType
		builder.Profile = name
		cf, err := builder.BuildWithProfile()

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

func (c *autoConfiguration) Instantiate(configuration interface{})  {
	cv := reflect.ValueOf(configuration)

	configType := cv.Type()
	log.Debug("type: ", configType)
	name := configType.Elem().Name()
	log.Debug("fieldName: ", name)

	// call Init
	numOfMethod := cv.NumMethod()
	log.Debug("methods: ", numOfMethod)

	for mi := 0; mi < numOfMethod; mi++ {
		method := configType.Method(mi)
		methodName := method.Name
		log.Debugf("method: %v", methodName)
		numIn := method.Type.NumIn()
		// only 1 arg is supported so far
		if numIn == 1 {
			argv := make([]reflect.Value, numIn)
			argv[0] = reflect.ValueOf(configuration)
			retVal := method.Func.Call(argv)
			// save instance
			if retVal[0].CanInterface() {
				instance := retVal[0].Interface()
				log.Debugf("instantiated: %v", instance)
				c.instances[utils.LowerFirst(methodName)] = instance
			}
		}
	}
}

func (c *autoConfiguration) Configurations() map[string]interface{} {
	return c.configurations
}

func (c *autoConfiguration) Configuration(name string) interface{} {
	return c.configurations[name]
}

func (c *autoConfiguration) Instances() map[string]interface{} {
	return c.instances
}

func (c *autoConfiguration) Instance(name string) interface{} {
	return c.instances[name]
}