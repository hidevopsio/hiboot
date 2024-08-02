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
	"path/filepath"
	"sync"

	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/factory/depends"
	"github.com/hidevopsio/hiboot/pkg/inject"
	"github.com/hidevopsio/hiboot/pkg/inject/annotation"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/system"
	"github.com/hidevopsio/hiboot/pkg/system/types"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
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
	at.Qualifier `value:"github.com/hidevopsio/hiboot/pkg/factory.instantiateFactory"`

	instanceContainer       factory.InstanceContainer
	scopedInstanceContainer factory.InstanceContainer
	components              []*factory.MetaData
	resolved                []*factory.MetaData
	defaultProperties       cmap.ConcurrentMap
	categorized             map[string][]*factory.MetaData
	inject                  inject.Inject
	builder                 system.Builder
	mutex                   sync.Mutex
}

// NewInstantiateFactory the constructor of instantiateFactory
func NewInstantiateFactory(instanceMap cmap.ConcurrentMap, components []*factory.MetaData, defaultProperties cmap.ConcurrentMap) factory.InstantiateFactory {
	if defaultProperties == nil {
		defaultProperties = cmap.New()
	}

	f := &instantiateFactory{
		instanceContainer: newInstanceContainer(instanceMap),
		components:        components,
		defaultProperties: defaultProperties,
		categorized:       make(map[string][]*factory.MetaData),
	}
	f.inject = inject.NewInject(f)

	// create new builder
	workDir := io.GetWorkDir()

	sa := new(system.App)
	ss := new(system.Server)
	sl := new(system.Logging)
	syscfg := system.NewConfiguration()

	customProps := defaultProperties.Items()
	f.builder = system.NewPropertyBuilder(
		filepath.Join(workDir, config),
		customProps,
	)

	f.Append(syscfg, sa, ss, sl, f, f.builder)

	initScopedFactory(f)

	return f
}

// Initialized check if factory is initialized
func (f *instantiateFactory) Initialized() bool {
	return f.instanceContainer != nil
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

// SetDefaultProperty set default property
func (f *instantiateFactory) SetDefaultProperty(name string, value interface{}) factory.InstantiateFactory {
	f.builder.SetDefaultProperty(name, value)
	return f
}

// Append append to component and instanceContainer container
func (f *instantiateFactory) Append(i ...interface{}) {
	for _, inst := range i {
		f.AppendComponent(inst)
		_ = f.SetInstance(inst)
	}
}

// AppendComponent append component
func (f *instantiateFactory) AppendComponent(c ...interface{}) {
	metaData := factory.NewMetaData(c...)
	f.components = append(f.components, metaData)
}

// injectDependency inject dependency
func (f *instantiateFactory) injectDependency(instanceContainer factory.InstanceContainer, item *factory.MetaData) (err error) {
	var name string
	var inst interface{}
	switch item.Kind {
	case types.Func:
		inst, err = f.inject.IntoFunc(instanceContainer, item.MetaObject)
		name = item.Name
		// TODO: should report error when err is not nil
		if err == nil {
			log.Debugf("inject into func: %v %v", item.ShortName, item.Type)
		}
	case types.Method:
		inst, err = f.inject.IntoMethod(instanceContainer, item.ObjectOwner, item.MetaObject)
		name = item.Name
		if err != nil {
			return
		}
		log.Debugf("inject into method: %v %v", item.Name, item.Type)
	default:
		name, inst = item.Name, item.MetaObject
	}
	if inst != nil {
		// inject into object
		err = f.inject.IntoObject(instanceContainer, inst)
		// TODO: remove duplicated code
		qf := annotation.GetAnnotation(inst, at.Qualifier{})
		if qf != nil {
			name = qf.Field.StructField.Tag.Get("value")
			//log.Debugf("name: %v, Qualifier: %v", item.Name, name)
		}

		if name != "" {
			// save object
			item.Instance = inst
			// set item
			err = f.SetInstance(instanceContainer, name, item)
		}
	}
	return
}

// InjectDependency inject dependency
func (f *instantiateFactory) InjectDependency(instanceContainer factory.InstanceContainer, object interface{}) (err error) {
	return f.injectDependency(instanceContainer, factory.CastMetaData(object))
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
		// log.Debugf("build component: %v %v", idx, item.Type)
		if item.Scope != "" {
			//log.Debugf("at.Scope: %v", item.MetaObject)
			err = f.SetInstance(item)
		} else {
			// inject dependencies into function
			// components, controllers
			// TODO: should save the upstream dependencies that contains item.Scope annotation for runtime injection
			err = f.injectDependency(f.instanceContainer, item)
		}
	}
	if err == nil {
		log.Debugf("Injected dependencies")
	}
	return
}

// SetInstance save instanceContainer
func (f *instantiateFactory) SetInstance(params ...interface{}) (err error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	var instanceContainer factory.InstanceContainer
	switch params[0].(type) {
	case factory.InstanceContainer:
		instanceContainer = params[0].(factory.InstanceContainer)
		params = params[1:]
	default:
		instanceContainer = f.instanceContainer
		if len(params) > 1 && params[0] == nil {
			params = params[1:]
		}
	}

	name, inst := factory.ParseParams(params...)

	if inst == nil {
		return ErrNotInitialized
	}

	metaData := factory.CastMetaData(inst)
	if metaData == nil {
		metaData = factory.NewMetaData(inst)
	}

	if metaData != nil {
		if metaData.Scope != "" && f.scopedInstanceContainer != nil {
			_ = f.scopedInstanceContainer.Set(name, inst)
		} else {
			err = instanceContainer.Set(name, inst)
			// categorize instances
			obj := metaData.MetaObject
			if metaData.Instance != nil {
				obj = metaData.Instance
			}

			annotations := annotation.GetAnnotations(obj)
			if annotations != nil {
				for _, item := range annotations.Items {
					typeName := reflector.GetLowerCamelFullNameByType(item.Field.StructField.Type)
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

// GetInstance get instanceContainer by name
func (f *instantiateFactory) GetInstance(params ...interface{}) (retVal interface{}) {
	switch params[0].(type) {
	case factory.InstanceContainer:
		inst := params[0].(factory.InstanceContainer)
		params = params[1:]
		retVal = inst.Get(params...)
	default:
		if len(params) > 1 && params[0] == nil {
			params = params[1:]
		}

	}
	// if it does not found from instanceContainer, try to find it from f.instanceContainer
	if retVal == nil {
		retVal = f.instanceContainer.Get(params...)
	}
	return
}

// GetInstances get instanceContainer by name
func (f *instantiateFactory) GetInstances(params ...interface{}) (retVal []*factory.MetaData) {
	if f.Initialized() {
		name, _ := factory.ParseParams(params...)
		retVal = f.categorized[name]
	}
	return
}

// Items return instanceContainer map
func (f *instantiateFactory) Items() map[string]interface{} {
	return f.instanceContainer.Items()
}

// DefaultProperties return default properties
func (f *instantiateFactory) DefaultProperties() map[string]interface{} {
	dp := f.defaultProperties.Items()
	return dp
}

// InjectIntoObject inject into object
func (f *instantiateFactory) InjectIntoObject(instanceContainer factory.InstanceContainer, object interface{}) error {
	return f.inject.IntoObject(instanceContainer, object)
}

// InjectDefaultValue inject default value
func (f *instantiateFactory) InjectDefaultValue(object interface{}) error {
	return f.inject.DefaultValue(object)
}

// InjectIntoFunc inject into func
func (f *instantiateFactory) InjectIntoFunc(instanceContainer factory.InstanceContainer, object interface{}) (retVal interface{}, err error) {
	return f.inject.IntoFunc(instanceContainer, object)
}

// InjectIntoMethod inject into method
func (f *instantiateFactory) InjectIntoMethod(instanceContainer factory.InstanceContainer, owner, object interface{}) (retVal interface{}, err error) {
	return f.inject.IntoMethod(instanceContainer, owner, object)
}

func (f *instantiateFactory) Replace(source string) (retVal interface{}) {
	retVal = f.builder.Replace(source)
	return
}

// InjectScopedDependencies inject context aware objects
func (f *instantiateFactory) InjectScopedDependencies(instanceContainer factory.InstanceContainer, dps []*factory.MetaData) (err error) {
	for _, d := range dps {
		if len(d.DepMetaData) > 0 {
			err = f.InjectScopedDependencies(instanceContainer, d.DepMetaData)
			if err != nil {
				return
			}
		}
		if d.Scope != "" {
			// making sure that the scoped instanceContainer does not exist before the dependency injection
			if instanceContainer.Get(d.Name) == nil || d.Scope == factory.ScopePrototype {
				newItem := factory.CloneMetaData(d)
				err = f.InjectDependency(instanceContainer, newItem)
				if err != nil {
					return
				}
			}
		}
	}
	return
}

// InjectScopedObjects inject context aware objects
func (f *instantiateFactory) InjectScopedObjects(ctx context.Context, dps []*factory.MetaData, ic factory.InstanceContainer) (instanceContainer factory.InstanceContainer, err error) {
	log.Debugf(">>> InjectScopedObjects(%x) ...", &ctx)

	// create new runtime instanceContainer
	instanceContainer = ic
	if instanceContainer == nil {
		instanceContainer = newInstanceContainer(nil)
	}

	// update context
	if ctx != nil {
		err = instanceContainer.Set(reflector.GetLowerCamelFullName(new(context.Context)), ctx)
		if err != nil {
			log.Error(err)
			return
		}
	}

	err = f.InjectScopedDependencies(instanceContainer, dps)
	if err != nil {
		log.Error(err)
	}

	return
}
