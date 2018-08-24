package app

import (
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"reflect"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/inject"
	"github.com/hidevopsio/hiboot/pkg/system"
)

type Application interface {
	Init(args ...interface{}) error
	Run()
}

type BaseApplication struct {
	WorkDir             string
	configurations      cmap.ConcurrentMap
	instances           cmap.ConcurrentMap
	configurableFactory *factory.ConfigurableFactory
}

var (
	preConfigContainer       cmap.ConcurrentMap
	configContainer          cmap.ConcurrentMap
	postConfigContainer      cmap.ConcurrentMap
	instanceContainer		 cmap.ConcurrentMap
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

	if _, ok := instanceContainer.Get(name); ok {
		log.Fatalf("instance name %v is already taken!", name)
	}
	instanceContainer.Set(name, inst)
}

// BeforeInitialization ?
func (a *BaseApplication) Init(args ...interface{}) error  {
	a.WorkDir = io.GetWorkDir()

	a.configurations = cmap.New()
	a.instances = instanceContainer

	instanceFactory := new(factory.InstanceFactory)
	instanceFactory.Init(a.instances)
	a.instances.Set("instanceFactory", instanceFactory)

	configurableFactory := new(factory.ConfigurableFactory)
	configurableFactory.InstanceFactory = instanceFactory
	a.instances.Set("configurableFactory", configurableFactory)

	configurableFactory.Init(a.configurations)
	configurableFactory.BuildSystemConfig(system.Configuration{})

	a.configurableFactory = configurableFactory

	inject.SetInstance(a.instances)

	return nil
}

func (a *BaseApplication) BuildConfigurations()  {
	a.configurableFactory.Build(preConfigContainer, configContainer, postConfigContainer)
}

func (a *BaseApplication) ConfigurableFactory() *factory.ConfigurableFactory  {
	return a.configurableFactory
}