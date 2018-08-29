package app

import (
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"reflect"
	"github.com/hidevopsio/hiboot/pkg/inject"
	"github.com/hidevopsio/hiboot/pkg/system"
	"errors"
	"github.com/hidevopsio/hiboot/pkg/factory/autoconfigure"
	"github.com/hidevopsio/hiboot/pkg/factory/instance"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/kataras/iris/context"
)

type Application interface {
	Init(args ...interface{}) error
	Run()
	RegisterController(controller interface{}) error
	Use(handlers ...context.Handler)
}

type Configuration interface {}
type PreConfiguration interface {}
type PostConfiguration interface {}

type BaseApplication struct {
	WorkDir             string
	configurations      cmap.ConcurrentMap
	instances           cmap.ConcurrentMap
	potatoes		 	cmap.ConcurrentMap
	configurableFactory *autoconfigure.ConfigurableFactory
	systemConfig 		*system.Configuration
	postProcessor  		postProcessor
}

var (
	preConfigContainer       cmap.ConcurrentMap
	configContainer          cmap.ConcurrentMap
	postConfigContainer      cmap.ConcurrentMap
	instanceContainer		 cmap.ConcurrentMap

	InvalidObjectTypeError        = errors.New("[app] invalid object type")
	ConfigurationNameIsTakenError = errors.New("[app] configuration name is already taken")
	ComponentNameIsTakenError     = errors.New("[app] component name is already taken")

	logo = `

`
)

func init() {
	preConfigContainer = cmap.New()
	configContainer = cmap.New()
	postConfigContainer = cmap.New()
	instanceContainer = cmap.New()
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

func validateObjectType(inst interface{}) error  {
	val := reflect.ValueOf(inst)
	//log.Println(val.Kind())
	//log.Println(reflect.Indirect(val).Kind())
	if val.Kind() == reflect.Ptr && reflect.Indirect(val).Kind() == reflect.Struct {
		return nil
	}
	return InvalidObjectTypeError
}

// AutoConfiguration
func AutoConfiguration(params ...interface{}) (err error) {
	if len(params) == 0 || params[0] == nil {
		err = InvalidObjectTypeError
		log.Error(err)
		return
	}
	name, inst := parseInstance("Configuration", params...)
	if name == "" || name == "configuration" {
		name = reflector.ParseObjectPkgName(params[0])
	}

	ifcField := reflector.GetEmbeddedInterfaceField(inst)
	var c cmap.ConcurrentMap
	if ifcField.Anonymous {
		switch ifcField.Name {
		case "Configuration":
			c = configContainer
		case "PreConfiguration":
			c = preConfigContainer
		case "PostConfiguration":
			c = postConfigContainer
		default:
			err = InvalidObjectTypeError
			return
		}
	} else {
		err = InvalidObjectTypeError
		log.Error(err)
		return
	}

	if _, ok := c.Get(name); ok {
		err = ConfigurationNameIsTakenError
		log.Error(err)
		return
	}

	err = validateObjectType(inst)
	if err == nil {
		c.Set(name, inst)
	} else {
		log.Error(err)
	}

	return err
}

// Component
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
func (a *BaseApplication) Init(args ...interface{}) error  {
	a.WorkDir = io.GetWorkDir()

	a.configurations = cmap.New()
	a.instances = instanceContainer

	instanceFactory := new(instance.InstanceFactory)
	instanceFactory.Initialize(a.instances)
	a.instances.Set("instanceFactory", instanceFactory)

	configurableFactory := new(autoconfigure.ConfigurableFactory)
	configurableFactory.InstanceFactory = instanceFactory
	a.instances.Set("configurableFactory", configurableFactory)

	inject.SetInstance(a.instances)

	configurableFactory.Initialize(a.configurations)

	a.systemConfig = new(system.Configuration)
	configurableFactory.BuildSystemConfig(a.systemConfig)

	a.configurableFactory = configurableFactory

	return nil
}

// Config returns application config
func (a *BaseApplication) SystemConfig() *system.Configuration {
	return a.systemConfig
}

func (a *BaseApplication) BuildConfigurations()  {
	a.configurableFactory.Build(preConfigContainer, configContainer, postConfigContainer)
}

func (a *BaseApplication) ConfigurableFactory() *autoconfigure.ConfigurableFactory  {
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
