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

package inject

import (
	"errors"
	"fmt"
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
	"reflect"
)

const (
	initMethodName = "Init"
)

var (
	// ErrNotImplemented the interface is not implemented
	ErrNotImplemented = errors.New("[inject] interface is not implemented")

	// ErrInvalidObject the object is invalid
	ErrInvalidObject = errors.New("[inject] invalid object")

	// ErrInvalidTagName the tag name is invalid
	ErrInvalidTagName = errors.New("[inject] invalid tag name, e.g. exampleTag")

	// ErrSystemConfiguration system is not configured
	ErrSystemConfiguration = errors.New("[inject] system is not configured")

	// ErrInvalidFunc the function is invalid
	ErrInvalidFunc = errors.New("[inject] invalid func")

	// ErrFactoryIsNil factory is invalid
	ErrFactoryIsNil = errors.New("[inject] factory is nil")

	tagsContainer []Tag

	//instancesMap cmap.ConcurrentMap
	appFactory factory.ConfigurableFactory
)

// SetFactory set factory from app
func SetFactory(f factory.ConfigurableFactory) {
	//if fct == nil {
	appFactory = f
	//}
}

// AddTag add new tag
func AddTag(tag Tag) {
	if tag != nil {
		tagsContainer = append(tagsContainer, tag)
	}
}

func getInstanceByName(name string, typ reflect.Type) (inst interface{}) {
	typ = reflector.IndirectType(typ)

	if name == "" {
		name = typ.Name()
	}

	name = io.DirName(typ.PkgPath()) + "." + str.ToLowerCamel(name)
	if appFactory != nil {
		inst = appFactory.GetInstance(name)
	}
	return
}

func saveInstance(name string, inst interface{}) (err error) {
	name = str.LowerFirst(name)
	if appFactory != nil {
		err = appFactory.SetInstance(name, inst)
	}
	return
}

// DefaultValue injects instance into the tagged field with `inject:"instanceName"`
func DefaultValue(object interface{}) error {
	return IntoObjectValue(reflect.ValueOf(object), new(defaultTag))
}

// IntoObject injects instance into the tagged field with `inject:"instanceName"`
func IntoObject(object interface{}) error {
	return IntoObjectValue(reflect.ValueOf(object))
}

// IntoObjectValue injects instance into the tagged field with `inject:"instanceName"`
func IntoObjectValue(object reflect.Value, tags ...Tag) error {
	var err error

	// TODO refactor IntoObject
	if appFactory == nil {
		return ErrSystemConfiguration
	}

	obj := reflector.Indirect(object)
	if obj.Kind() != reflect.Struct {
		log.Warnf("[inject] ignore object: %v, kind: %v", object, obj.Kind())
		return ErrInvalidObject
	}

	var targetTags []Tag
	if len(tags) != 0 {
		targetTags = tags
	} else {
		targetTags = tagsContainer
	}

	// field injection
	for _, f := range reflector.DeepFields(object.Type()) {
		//log.Debugf("parent: %v, name: %v, type: %v, tag: %v", obj.Type(), f.Name, f.Type, f.Tag)
		// check if object has value field to be injected
		var injectedObject interface{}

		ft := f.Type
		if f.Type.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}

		// set field object
		var fieldObj reflect.Value
		if obj.IsValid() && obj.Kind() == reflect.Struct {
			fieldObj = obj.FieldByName(f.Name)
		}

		// TODO: assume that the f.Name of value and inject tag is not the same

		injectedObject = getInstanceByName(f.Name, f.Type)
		if injectedObject == nil {
			for _, tagImpl := range targetTags {
				tagName := reflector.ParseObjectName(tagImpl, "Tag")
				if tagName == "" {
					return ErrInvalidTagName
				}
				tag, ok := f.Tag.Lookup(tagName)
				if ok {
					tagImpl.Init(appFactory)
					injectedObject = tagImpl.Decode(object, f, tag)
					if injectedObject != nil {
						//if tagImpl.IsSingleton() {
						//	err := saveInstance(f.Name, injectedObject)
						//	if err != nil {
						//		log.Warnf("instance %v is already exist", f.Name)
						//	}
						//}
						// ONLY one tag should be used for dependency injection
						break
					}
				}
			}
		}

		if injectedObject != nil && fieldObj.CanSet() {
			fov := reflect.ValueOf(injectedObject)
			fieldObj.Set(fov)
			log.Debugf("Injected %v.(%v) into %v.%v", injectedObject, fov.Type(), obj.Type(), f.Name)
		}

		//log.Debugf("- kind: %v, %v, %v, %v", obj.Kind(), object.Type(), fieldObj.Type(), f.Name)
		//log.Debugf("isValid: %v, canSet: %v", fieldObj.IsValid(), fieldObj.CanSet())
		filedObject := reflect.Indirect(fieldObj)
		filedKind := filedObject.Kind()
		canNested := filedKind == reflect.Struct
		if canNested && fieldObj.IsValid() && fieldObj.CanSet() && filedObject.Type() != obj.Type() {
			err = IntoObjectValue(fieldObj, tags...)
		}
	}
	return err
}

func parseFuncOrMethodInput(inType reflect.Type) (paramValue reflect.Value, ok bool) {
	inType = reflector.IndirectType(inType)
	inTypeName := inType.Name()
	inst := getInstanceByName(inTypeName, inType)
	ok = true
	if inst == nil {
		log.Debug(inType.Kind())
		switch inType.Kind() {
		// interface and slice creation is not supported
		case reflect.Interface, reflect.Slice:
			ok = false
			break
		default:
			// should find instance in the component container first

			// if it is not found, then create new instance
			paramValue = reflect.New(inType)
			inst = paramValue.Interface()
			// TODO: inTypeName
			err := saveInstance(inTypeName, inst)
			if err != nil {
				log.Warnf("instance %v is already exist", inTypeName)
			}
		}
	}

	if inst != nil {
		paramValue = reflect.ValueOf(inst)
	}
	return
}

// IntoFunc inject object into func and return instance
func IntoFunc(object interface{}) (retVal interface{}, err error) {
	fn := reflect.ValueOf(object)
	if fn.Kind() == reflect.Func {
		numIn := fn.Type().NumIn()
		inputs := make([]reflect.Value, numIn)
		// TODO: should load function inputs when resolving dependencies to improve performance
		for i := 0; i < numIn; i++ {
			fnInType := fn.Type().In(i)
			val, ok := parseFuncOrMethodInput(fnInType)
			if ok {
				inputs[i] = val
				//log.Debugf("Injected %v into func parameter %v", val, fnInType)
			} else {
				return nil, fmt.Errorf("%v is not injected", fnInType.Name())
			}

			paramObject := reflect.Indirect(val)
			if val.IsValid() && paramObject.IsValid() && paramObject.Kind() == reflect.Struct {
				err = IntoObjectValue(val)
			}
		}
		results := fn.Call(inputs)
		if len(results) != 0 {
			retVal = results[0].Interface()
			return
		}
		return
	}
	err = ErrInvalidFunc
	return
}

//IntoMethod inject object into func and return instance
func IntoMethod(object interface{}, m interface{}) (retVal interface{}, err error) {
	if object != nil && m != nil {
		method := m.(reflect.Method)
		numIn := method.Type.NumIn()
		inputs := make([]reflect.Value, numIn)
		inputs[0] = reflect.ValueOf(object)
		for i := 1; i < numIn; i++ {
			fnInType := method.Type.In(i)
			val, ok := parseFuncOrMethodInput(fnInType)
			if ok {
				inputs[i] = val
			} else {
				return nil, fmt.Errorf("%v is not injected", fnInType.Name())
			}

			paramObject := reflect.Indirect(val)
			if val.IsValid() && paramObject.IsValid() && paramObject.Kind() == reflect.Struct {
				err = IntoObjectValue(val)
			}
		}
		results := method.Func.Call(inputs)
		if len(results) != 0 {
			retVal = results[0].Interface()
			return
		} else {
			return
		}
	}
	err = ErrInvalidFunc
	return
}
