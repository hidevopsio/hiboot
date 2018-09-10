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

package instantiate

import (
	"errors"
	"fmt"
	"github.com/hidevopsio/hiboot/pkg/inject"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
	"reflect"
	"strings"
)

var (
	NotInitializedError    = errors.New("InstantiateFactory is not initialized")
	InvalidObjectTypeError = errors.New("[factory] invalid Component")
)

type InstantiateFactory struct {
	instanceMap cmap.ConcurrentMap
}

func (f *InstantiateFactory) Initialize(instanceMap cmap.ConcurrentMap) {
	f.instanceMap = instanceMap
}

func (f *InstantiateFactory) Initialized() bool {
	return f.instanceMap != nil
}

func (f *InstantiateFactory) IsValidObjectType(inst interface{}) bool {
	val := reflect.ValueOf(inst)
	//log.Println(val.Kind())
	//log.Println(reflect.Indirect(val).Kind())
	if val.Kind() == reflect.Ptr && reflect.Indirect(val).Kind() == reflect.Struct {
		return true
	}
	return false
}

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
			inst, err = inject.IntoFunc(inst)
			if err != nil {
				return "", nil
			}
		}
		name = reflector.ParseObjectName(inst, eliminator)
		if name == "" || strings.ToLower(name) == strings.ToLower(eliminator) {
			name = reflector.ParseObjectPkgName(inst)
		}
	}

	return
}

func (f *InstantiateFactory) BuildComponents(components [][]interface{}) (err error) {
	for _, item := range components {
		name, inst := f.ParseInstance("", item...)
		if inst == nil {
			return InvalidObjectTypeError
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

func (f *InstantiateFactory) SetInstance(name string, instance interface{}) (err error) {
	if !f.Initialized() {
		return NotInitializedError
	}

	if _, ok := f.instanceMap.Get(name); ok {
		return fmt.Errorf("instance name %v is already taken", name)
	}

	f.instanceMap.Set(name, instance)
	return
}

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

func (f *InstantiateFactory) Items() map[string]interface{} {
	if !f.Initialized() {
		return nil
	}
	return f.instanceMap.Items()
}
