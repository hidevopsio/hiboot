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

// Package instantiate implement InstantiateFactory
package instantiate

import (
	"errors"
	"hidevops.io/hiboot/pkg/app/web/context"
	"hidevops.io/hiboot/pkg/factory"
	"hidevops.io/hiboot/pkg/factory/depends"
	"hidevops.io/hiboot/pkg/inject"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/system"
	"hidevops.io/hiboot/pkg/system/types"
	"hidevops.io/hiboot/pkg/utils/cmap"
	"hidevops.io/hiboot/pkg/utils/io"
	"hidevops.io/hiboot/pkg/utils/reflector"
	"path/filepath"
)

var (
	// ErrNotInitialized InstantiateFactory is not initialized
	ErrNotInitialized = errors.New("[factory] InstantiateFactory is not initialized")

	// ErrInvalidObjectType invalid object type
	ErrInvalidObjectType = errors.New("[factory] invalid object type")
)

const (
	application = "application"
	config      = "config"
	yaml        = "yaml"
)

// InstantiateFactory is the factory that responsible for object instantiation
type instantiateFactory struct {
	instance             factory.Instance
	contextAwareInstance factory.Instance
	components           []*factory.MetaData
	resolved             []*factory.MetaData
	customProperties     cmap.ConcurrentMap
	categorized          map[string][]*factory.MetaData
	inject               inject.Inject
	builder              system.Builder
}

// NewInstantiateFactory the constructor of instantiateFactory
func NewInstantiateFactory(instanceMap cmap.ConcurrentMap, components []*factory.MetaData, customProperties cmap.ConcurrentMap) factory.InstantiateFactory {
	if customProperties == nil {
		customProperties = cmap.New()
	}
	f := &instantiateFactory{
		instance:         newInstance(instanceMap),
		components:       components,
		customProperties: customProperties,
		categorized:      make(map[string][]*factory.MetaData),
	}
	f.inject = inject.NewInject(f)

	// create new builder
	workDir := io.GetWorkDir()
	systemConfig := new(system.Configuration)
	f.SetInstance(systemConfig)
	customProps := customProperties.Items()
	f.builder = system.NewBuilder(systemConfig,
		filepath.Join(workDir, config),
		application,
		yaml,
		customProps,
	)
	return f
}

// Initialized check if factory is initialized
func (f *instantiateFactory) Initialized() bool {
	return f.instance != nil
}

// Builder get builder
func (f *instantiateFactory) Builder() (builder system.Builder) {
	return f.builder
}

// GetProperty get property
func (f *instantiateFactory) GetProperty(name string) (retVal interface{}) {
	retVal = f.builder.GetProperty(name)
	return
}

// SetProperty get property
func (f *instantiateFactory) SetProperty(name string, value interface{}) factory.InstantiateFactory {
	f.builder.SetProperty(name, value)
	return f
}

// AppendComponent append component
func (f *instantiateFactory) AppendComponent(c ...interface{}) {
	metaData := factory.NewMetaData(c...)
	f.components = append(f.components, metaData)
}

// injectDependency inject dependency
func (f *instantiateFactory) injectDependency(item *factory.MetaData) (err error) {
	var name string
	var inst interface{}
	switch item.Kind {
	case types.Func:
		inst, err = f.inject.IntoFunc(item.MetaObject)
		name = item.Name
		// TODO: should report error when err is not nil
		if err == nil {
			log.Debugf("inject into func: %v %v", item.ShortName, item.Type)
		}
	case types.Method:
		inst, err = f.inject.IntoMethod(item.ObjectOwner, item.MetaObject)
		name = item.Name
		if err == nil {
			log.Debugf("inject into method: %v %v", item.ShortName, item.Type)
		}
	default:
		name, inst = item.Name, item.MetaObject
	}
	if inst != nil {
		// inject into object
		err = f.inject.IntoObject(inst)
		tagName, ok := reflector.FindEmbeddedFieldTag(inst, "Qualifier", "name")
		if ok {
			name = tagName
			log.Debugf("name: %v, Qualifier: %v, ok: %v", item.Name, name, ok)
		}

		if name != "" {
			// save object
			item.Instance = inst
			// set item
			err = f.SetInstance(name, item)
		}
	}
	return
}

// InjectDependency inject dependency
func (f *instantiateFactory) InjectDependency(object interface{}) (err error) {
	return f.injectDependency(factory.CastMetaData(object))
}

// BuildComponents build all registered components
func (f *instantiateFactory) BuildComponents() (err error) {
	// first resolve the dependency graph
	var resolved []*factory.MetaData
	log.Debugf("Resolving dependencies")
	resolved, err = depends.Resolve(f.components)
	f.resolved = resolved
	log.Debugf("Injecting dependencies")
	// then build components
	for _, item := range resolved {
		log.Debugf("build component: %v", item.Type)
		if item.ContextAware {
			//log.Debugf("at.ContextAware: %v", item.MetaObject)
			f.SetInstance(item)
		} else {
			// inject dependencies into function
			// components, controllers
			f.injectDependency(item)
		}
	}
	if err == nil {
		log.Debugf("Injected dependencies")
	}
	return
}

// SetInstance save instance
func (f *instantiateFactory) SetInstance(params ...interface{}) (err error) {
	name, inst := factory.ParseParams(params...)

	if inst == nil {
		return ErrNotInitialized
	}

	metaData := factory.CastMetaData(inst)
	if metaData == nil {
		metaData = factory.NewMetaData(inst)
	}

	if metaData != nil {
		if metaData.ContextAware && f.contextAwareInstance != nil {
			f.contextAwareInstance.Set(name, inst)
		} else {
			err = f.instance.Set(name, inst)
			// categorize instances
			if metaData != nil {
				obj := metaData.MetaObject
				if metaData.Instance != nil {
					obj = metaData.Instance
				}
				fields := reflector.GetEmbeddedFields(obj)
				for _, field := range fields {
					typeName := reflector.GetLowerCamelFullNameByType(field.Type)
					categorised, ok := f.categorized[typeName]
					if !ok {
						categorised = make([]*factory.MetaData, 0)
					}
					f.categorized[typeName] = append(categorised, metaData)
				}
			}
		}
	}

	return
}

// GetInstance get instance by name
func (f *instantiateFactory) GetInstance(params ...interface{}) (retVal interface{}) {
	if f.contextAwareInstance != nil {
		retVal = f.contextAwareInstance.Get(params...)
	}

	if retVal == nil {
		retVal = f.instance.Get(params...)
	}
	return
}

// GetInstances get instance by name
func (f *instantiateFactory) GetInstances(params ...interface{}) (retVal []*factory.MetaData) {
	if f.Initialized() {
		name, _ := factory.ParseParams(params...)
		retVal = f.categorized[name]
	}
	return
}

// Items return instance map
func (f *instantiateFactory) Items() map[string]interface{} {
	return f.instance.Items()
}

// Items return instance map
func (f *instantiateFactory) CustomProperties() map[string]interface{} {
	return f.customProperties.Items()
}

// InjectIntoObject inject into object
func (f *instantiateFactory) InjectIntoObject(object interface{}) error {
	return f.inject.IntoObject(object)
}

// InjectDefaultValue inject default value
func (f *instantiateFactory) InjectDefaultValue(object interface{}) error {
	return f.inject.DefaultValue(object)
}

// InjectIntoFunc inject into func
func (f *instantiateFactory) InjectIntoFunc(object interface{}) (retVal interface{}, err error) {
	return f.inject.IntoFunc(object)
}

// InjectIntoMethod inject into method
func (f *instantiateFactory) InjectIntoMethod(owner, object interface{}) (retVal interface{}, err error) {
	return f.inject.IntoMethod(owner, object)
}

func (f *instantiateFactory) Replace(source string) (retVal interface{}) {
	retVal = f.builder.Replace(source)
	return
}

// InjectContextAwareObject inject context aware objects
func (f *instantiateFactory) injectContextAwareDependencies(dps []*factory.MetaData) (err error) {
	for _, d := range dps {
		if len(d.DepMetaData) > 0 {
			f.injectContextAwareDependencies(d.DepMetaData)
		}
		if d.ContextAware {
			// making sure that the context aware instance does not exist before the dependency injection
			if f.contextAwareInstance.Get(d.Name) == nil {
				newItem := factory.CloneMetaData(d)
				err = f.InjectDependency(newItem)
			}
		}
	}
	return
}

// InjectContextAwareObject inject context aware objects
func (f *instantiateFactory) InjectContextAwareObjects(ctx context.Context, dps []*factory.MetaData) (contextAwareInstance factory.Instance, err error) {
	log.Debugf(">>> InjectContextAwareObjects(%x) ...", &ctx)

	// create new runtime instance
	f.contextAwareInstance = newInstance(nil)

	// update context
	f.contextAwareInstance.Set(reflector.GetLowerCamelFullName(new(context.Context)), ctx)

	err = f.injectContextAwareDependencies(dps)

	contextAwareInstance = f.contextAwareInstance
	return
}
