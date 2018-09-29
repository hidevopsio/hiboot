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
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
	"reflect"
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/factory/depends"
	"github.com/hidevopsio/hiboot/pkg/inject"
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
	components []*factory.MetaData
}

// Initialize init the factory
func (f *InstantiateFactory) Initialize(instanceMap cmap.ConcurrentMap, components []*factory.MetaData) {
	f.instanceMap = instanceMap
	f.components = components
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


// Initialized check if factory is initialized
func (f *InstantiateFactory) AppendComponent(c ...interface{}) {
	f.components = append(f.components, factory.ParseParams("", c...))
}

// BuildComponents build all registered components
func (f *InstantiateFactory) BuildComponents() (err error) {
	//TODO: should sort components according to dependency tree first
	f.components, err = depends.Resolve(f.components)
	// then build components
	var obj interface{}
	var name string
	for _, item := range f.components {
		// inject dependencies into function
		// components, controllers
		if item.Kind == reflect.Func {
			obj, err = inject.IntoFunc(item.Object)
			name = item.Name
		} else {
			name, obj = item.Name, item.Object
		}

		if obj == nil {
			return ErrInvalidObjectType
		}
		// use interface name if it's available as use does not specify its name
		field := reflector.GetEmbeddedInterfaceField(obj)
		if field.Anonymous {
			err = f.SetInstance(field.Name, obj)
			//log.Debugf("component %v has embedded field: %v", inst, name)
		}
		if name == "" {
			continue
		}
		if f.IsValidObjectType(obj) {
			err = f.SetInstance(name, obj)
			if err != nil {
				return
			}
		}
	}
	return
}

// SetInstance save instance
func (f *InstantiateFactory) SetInstance(name string, instance interface{}) (err error) {
	if !f.Initialized() {
		return ErrNotInitialized
	}

	// force to use camel case name
	name = str.ToLowerCamel(name)

	if _, ok := f.instanceMap.Get(name); ok {
		return fmt.Errorf("instance name %v is already taken", name)
	}

	f.instanceMap.Set(name, instance)
	return
}

// GetInstance get instance by name
func (f *InstantiateFactory) GetInstance(name string) (inst interface{}) {
	if !f.Initialized() {
		return nil
	}
	//items := f.Items()
	//log.Debug(items)
	var ok bool
	if inst, ok = f.instanceMap.Get(name); !ok {
		return nil
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
