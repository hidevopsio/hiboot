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
	"fmt"
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/factory/depends"
	"github.com/hidevopsio/hiboot/pkg/inject"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/system/types"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"reflect"
)

var (
	// ErrNotInitialized InstantiateFactory is not initialized
	ErrNotInitialized = errors.New("[factory] InstantiateFactory is not initialized")

	// ErrInvalidObjectType invalid object type
	ErrInvalidObjectType = errors.New("[factory] invalid object type")
)

// InstantiateFactory is the factory that responsible for object instantiation
type InstantiateFactory struct {
	instanceMap cmap.ConcurrentMap
	components  []*factory.MetaData
	categorized map[string][]interface{}
}

// Initialize init the factory
func (f *InstantiateFactory) Initialize(instanceMap cmap.ConcurrentMap, components []*factory.MetaData) {
	f.instanceMap = instanceMap
	f.components = components
	f.categorized = make(map[string][]interface{})
}

// Initialized check if factory is initialized
func (f *InstantiateFactory) Initialized() bool {
	return f.instanceMap != nil
}

// IsValidObjectType check if is valid object type
func (f *InstantiateFactory) IsValidObjectType(inst interface{}) bool {
	val := reflect.ValueOf(inst)
	//log.Println(val.Kind())
	//log.Println(reflect.Indirect(val).Kind())
	if val.Kind() == reflect.Ptr && reflect.Indirect(val).Kind() == reflect.Struct {
		return true
	}
	return false
}

// AppendComponent append component
func (f *InstantiateFactory) AppendComponent(c ...interface{}) {
	metaData := factory.NewMetaData(c...)
	f.components = append(f.components, metaData)
}

// BuildComponents build all registered components
func (f *InstantiateFactory) BuildComponents() (err error) {
	// first resolve the dependency graph
	var resolved []*factory.MetaData
	resolved, err = depends.Resolve(f.components)
	// then build components
	var obj interface{}
	var name string
	for i, item := range resolved {
		// inject dependencies into function
		// components, controllers
		switch item.Kind {
		case types.Func:
			obj, err = inject.IntoFunc(item.Object)
			name = item.Name
			log.Debugf("%d: inject into func: %v - %v", i, item.Name, item.Type)
		case types.Method:
			obj, err = inject.IntoMethod(item.Context, item.Object)
			name = item.Name
			log.Debugf("%d: inject into method: %v - %v", i, item.Name, item.Type)
		default:
			name, obj = item.Name, item.Object
		}
		if obj != nil {
			// inject into object
			err = inject.IntoObject(obj)

			//field := reflector.GetEmbeddedField(obj)
			//if field.Anonymous {
			//	// use interface name if it's available as use does not specify its name
			//	name = io.DirName(field.PkgPath) + "." + field.Name
			//	log.Debugf("component %v has embedded field: %v", obj, name)
			//}
			if name != "" {
				err = f.SetInstance(name, obj)
			}
		}
	}
	return
}

// SetInstance save instance
func (f *InstantiateFactory) SetInstance(params ...interface{}) (err error) {
	if !f.Initialized() {
		return ErrNotInitialized
	}

	name, instance := factory.ParseParams(params...)

	if _, ok := f.instanceMap.Get(name); ok {
		return fmt.Errorf("instance name %v is already taken", name)
	}

	f.instanceMap.Set(name, instance)

	ifcField := reflector.GetEmbeddedField(instance)
	if ifcField.Anonymous {
		typeName := ifcField.Name
		categorised, ok := f.categorized[typeName]
		if !ok {
			categorised = make([]interface{}, 0)
		}
		f.categorized[typeName] = append(categorised, instance)
	}
	return
}

// GetInstance get instance by name
func (f *InstantiateFactory) GetInstance(params ...interface{}) (retVal interface{}) {
	if !f.Initialized() {
		return nil
	}

	name, _ := factory.ParseParams(params...)

	//items := f.Items()
	//log.Debug(items)

	var ok bool
	if retVal, ok = f.instanceMap.Get(name); !ok {
		return nil
	}
	return
}

// GetInstances get instance by name
func (f *InstantiateFactory) GetInstances(name string) (retVal []interface{}) {
	//items := f.Items()
	//log.Debug(items)
	if f.Initialized() {
		retVal = f.categorized[name]
	}
	return
}

// Items return instance map
func (f *InstantiateFactory) Items() map[string]interface{} {
	if !f.Initialized() {
		return nil
	}
	return f.instanceMap.Items()
}
