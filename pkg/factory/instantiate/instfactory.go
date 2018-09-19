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
	"github.com/hidevopsio/hiboot/pkg/inject"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
	"reflect"
	"strings"
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
}

// Initialize init the factory
func (f *InstantiateFactory) Initialize(instanceMap cmap.ConcurrentMap) {
	f.instanceMap = instanceMap
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

// ParseInstance parse object name and type
func (f *InstantiateFactory) ParseInstance(eliminator string, params ...interface{}) (name string, inst interface{}) {

	hasTwoParams := len(params) == 2 && reflect.TypeOf(params[0]).Kind() == reflect.String

	if hasTwoParams {
		inst = params[1]
		name = params[0].(string)
	} else {
		inst = params[0]
	}

	if !hasTwoParams {
		if reflect.TypeOf(inst).Kind() == reflect.Func {
			// call func
			var err error
			fn := inst
			inst, err = inject.IntoFunc(fn)
			if err != nil {
				return "", nil
			}
			// name should get from fn out
			fnVal := reflect.ValueOf(fn)
			if fnVal.Type().NumOut() == 1 {
				retTyp := fnVal.Type().Out(0)
				log.Debugf("[factory] constructor return type: %v, kind: %v", retTyp.Name(), retTyp.Kind())
				if retTyp.Kind() == reflect.Interface {
					name = retTyp.Name()
					return name, inst
				}
			}
		}
		name = reflector.ParseObjectName(inst, eliminator)
		if name == "" || strings.ToLower(name) == strings.ToLower(eliminator) {
			name = reflector.ParseObjectPkgName(inst)
		}
	}

	return
}

// BuildComponents build all registered components
func (f *InstantiateFactory) BuildComponents(components [][]interface{}) (err error) {
	for _, item := range components {
		name, inst := f.ParseInstance("", item...)
		if inst == nil {
			return ErrInvalidObjectType
		}
		// use interface name if it's available as use does not specify its name
		field := reflector.GetEmbeddedInterfaceField(inst)
		if field.Anonymous {
			name = str.ToLowerCamel(field.Name)
			//log.Debugf("component %v has embedded field: %v", inst, name)
		}
		if name == "" {
			continue
		}

		if f.IsValidObjectType(inst) {
			err = f.SetInstance(name, inst)
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
