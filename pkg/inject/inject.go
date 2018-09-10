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
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/system"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
	"reflect"
	"strings"
)

const (
	injectIdentifier = "inject"
	valueIdentifier  = "value"
	initMethodName   = "Init"
)

var (
	// NotImplementedError: the interface is not implemented
	NotImplementedError = errors.New("[inject] interface is not implemented")
	// InvalidObjectError: the object is invalid
	InvalidObjectError = errors.New("[inject] invalid object")
	// InvalidTagNameError the tag name is invalid
	InvalidTagNameError = errors.New("[inject] invalid tag name, e.g. exampleTag")
	// SystemConfigurationError system is not configured
	SystemConfigurationError = errors.New("[inject] system is not configured")

	// InvalidInputError
	InvalidInputError = errors.New("[inject] invalid input")
	// InvalidFuncError
	InvalidFuncError = errors.New("[inject] invalid func")

	// FactoryIsNilError
	FactoryIsNilError = errors.New("[inject] factory is nil")

	// TODO use cmap.ConcurrentMap for tagsContainer
	tagsContainer []Tag

	//instancesMap cmap.ConcurrentMap
	fct factory.ConfigurableFactory
)

// SetFactory set factory from app
func SetFactory(f factory.ConfigurableFactory) {
	if fct == nil {
		fct = f
	}
}

// AddTag add new tag
func AddTag(tag Tag) {
	tagsContainer = append(tagsContainer, tag)
}

func getInstanceByName(name string, instType reflect.Type) (inst interface{}) {
	name = str.ToLowerCamel(name)
	if fct != nil {
		inst = fct.GetInstance(name)
		// TODO: we should pro load all candidates into instances for improving performance.
		// if inst is nil, and the object type is an interface
		// then try to find the instance that embedded with the interface
		//if !ok && instType.Kind() == reflect.Interface {
		//	for _, ist := range fct.Items() {
		//		//log.Debug(n)
		//		if ist != nil && reflector.HasEmbeddedField(ist, instType.Name()) {
		//			inst = ist
		//			break
		//		}
		//	}
		//}
	}
	return
}

func saveInstance(name string, inst interface{}) error {
	name = str.LowerFirst(name)
	if fct == nil {
		return FactoryIsNilError
	}
	return fct.SetInstance(name, inst)
}

// IntoObject injects instance into the tagged field with `inject:"instanceName"`
func IntoObject(object interface{}) error {
	// TODO: save injected object to map to avoid re-injection
	return IntoObjectValue(reflect.ValueOf(object))
}

// IntoObjectValue injects instance into the tagged field with `inject:"instanceName"`
func IntoObjectValue(object reflect.Value) error {
	var err error

	// TODO refactor IntoObject
	if fct == nil {
		return SystemConfigurationError
	}

	obj := reflector.Indirect(object)
	if obj.Kind() != reflect.Struct {
		log.Errorf("[inject] object: %v, kind: %v", object, obj.Kind())
		return InvalidObjectError
	}

	sc := fct.GetInstance("systemConfiguration")
	if sc == nil {
		return SystemConfigurationError
	}
	systemConfig := sc.(*system.Configuration)

	cs := fct.GetInstance("configurations")
	if cs == nil {
		return SystemConfigurationError
	}
	configurations := cs.(cmap.ConcurrentMap)

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
			for _, tagImpl := range tagsContainer {
				tagName := reflector.ParseObjectName(tagImpl, "Tag")
				if tagName == "" {
					return InvalidTagNameError
				}
				tag, ok := f.Tag.Lookup(tagName)
				if ok {
					tagImpl.Init(systemConfig, configurations)
					injectedObject = tagImpl.Decode(object, f, tag)
					if injectedObject != nil {
						if tagImpl.IsSingleton() {
							err := saveInstance(f.Name, injectedObject)
							if err != nil {
								log.Warnf("instance %v is already exist", f.Name)
							}
						}
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
			err = IntoObjectValue(fieldObj)
		}
	}

	// method injection
	// Init, Setter
	method, ok := object.Type().MethodByName(initMethodName)
	if ok {
		numIn := method.Type.NumIn()
		inputs := make([]reflect.Value, numIn)
		inputs[0] = obj.Addr()
		var val reflect.Value
		for i := 1; i < numIn; i++ {
			val, ok = parseMethodInput(method.Type.In(i))
			if ok {
				inputs[i] = val
				//log.Debugf("inType: %v, name: %v, instance: %v", inType, inTypeName, inst)
				//log.Debugf("kind: %v == %v, %v, %v ", obj.Kind(), reflect.Struct, paramValue.IsValid(), paramValue.CanSet())
				paramObject := reflect.Indirect(val)
				if val.IsValid() && paramObject.IsValid() && paramObject.Type() != obj.Type() && paramObject.Kind() == reflect.Struct {
					err = IntoObjectValue(val)
				}
			} else {
				break
			}
		}
		// finally call Init method to inject
		if ok {
			method.Func.Call(inputs)
		}
	}

	return err
}

func parseMethodInput(inType reflect.Type) (paramValue reflect.Value, ok bool) {
	inType = reflector.IndirectType(inType)
	inTypeName := inType.Name()
	pkgName := io.DirName(inType.PkgPath())
	//log.Debugf("pkg: %v", pkgName)
	inst := getInstanceByName(inTypeName, inType)
	if inst == nil {
		alternativeName := strings.Title(pkgName) + inTypeName
		inst = getInstanceByName(alternativeName, inType)
	}
	ok = true
	if inst == nil {
		//log.Debug(inType.Kind())
		switch inType.Kind() {
		// interface and slice creation is not supported
		case reflect.Interface, reflect.Slice:
			ok = false
			break
		default:
			paramValue = reflect.New(inType)
			inst = paramValue.Interface()
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

// IntoMethod inject object into method and return instance
func IntoFunc(object interface{}) (retVal interface{}, err error) {
	fn := reflect.ValueOf(object)
	if fn.Kind() == reflect.Func {
		numIn := fn.Type().NumIn()
		inputs := make([]reflect.Value, numIn)
		for i := 0; i < numIn; i++ {
			val, ok := parseMethodInput(fn.Type().In(i))
			if ok {
				inputs[i] = val
			} else {
				return nil, InvalidInputError
			}

			paramObject := reflect.Indirect(val)
			if val.IsValid() && paramObject.IsValid() && paramObject.Kind() == reflect.Struct {
				err = IntoObjectValue(val)
			}
		}
		results := fn.Call(inputs)
		if len(results) != 0 {
			return results[0].Interface(), nil
		} else {
			return nil, nil
		}
	}
	return nil, InvalidFuncError
}