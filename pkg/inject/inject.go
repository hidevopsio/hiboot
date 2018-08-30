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
	"reflect"
	"errors"
	"strings"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"github.com/hidevopsio/hiboot/pkg/system"
)


const (
	injectIdentifier = "inject"
	valueIdentifier = "value"
	initMethodName = "Init"
)

var (
	InstanceContainerIsNilError = errors.New("[inject] instance container is nil")
	NotImplementedError = errors.New("[inject] interface is not implemented")
	InvalidObjectError      = errors.New("[inject] invalid object")
	UnsupportedInjectionTypeError      = errors.New("[inject] unsupported injection type")
	IllegalArgumentError = errors.New("[inject] input argument type can not be the same as receiver")
	TagIsAlreadyExistError = errors.New("[inject] tag is already exist")
	TagIsNilError = errors.New("[inject] tag is nil")
	InvalidTagNameError = errors.New("[inject] invalid tag name, e.g. exampleTag")
	SystemConfigurationError = errors.New("[inject] system is not configured")

	// TODO use cmap.ConcurrentMap for tagsContainer
	tagsContainer []Tag

	instancesMap cmap.ConcurrentMap
)

func SetInstance(i cmap.ConcurrentMap)  {
	instancesMap = i
}

// AddTag
func AddTag(tag Tag) {
	tagsContainer = append(tagsContainer, tag)
}

func getInstanceByName(name string, instType reflect.Type) (inst interface{}) {
	name = str.LowerFirst(name)
	var ok bool
	inst, ok = instancesMap.Get(name)
	// TODO: we should pro load all candidates into instances for improving performance.
	// if inst is nil, and the object type is an interface
	// then try to find the instance that embedded with the interface
	if !ok && instType.Kind() == reflect.Interface {
		for item := range instancesMap.IterBuffered() {
			ist := item.Val
			//log.Debug(n)
			if ist != nil && reflector.HasEmbeddedField(ist, instType.Name()) {
				inst = ist
				break
			}
		}
	}
	return
}

func saveInstance(name string, inst interface{}) {
	name = str.LowerFirst(name)

	if _, ok := instancesMap.Get(name); ok {
		log.Fatalf("%v is already taken", name)
	}

	instancesMap.Set(name, inst)
}


// IntoObject injects instance into the tagged field with `inject:"instanceName"`
func IntoObject(object interface{}) error {
	// TODO: save injected object to map to avoid re-injection
	return IntoObjectValue(reflect.ValueOf(object))
}

// IntoObject injects instance into the tagged field with `inject:"instanceName"`
func IntoObjectValue(object reflect.Value) error {
    var err error

    // TODO refactor IntoObject
    if instancesMap == nil {
    	return InstanceContainerIsNilError
	}

	obj := reflector.Indirect(object)
	if obj.Kind() != reflect.Struct {
		log.Errorf("[inject] object: %v, kind: %v", object, obj.Kind())
		return InvalidObjectError
	}

	sc, ok := instancesMap.Get("systemConfiguration")
	if !ok {
		return SystemConfigurationError
	}
	systemConfig := sc.(*system.Configuration)

	cs, ok := instancesMap.Get("configurations")
	if !ok {
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
							saveInstance(f.Name, injectedObject)
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
		injectByMethod := true
		for i := 1; i < numIn; i++ {
			inType := reflector.IndirectType(method.Type.In(i))
			var paramValue reflect.Value
			inTypeName := inType.Name()
			pkgName := io.DirName(inType.PkgPath())
			//log.Debugf("pkg: %v", pkgName)
			inst := getInstanceByName(inTypeName, inType)
			if inst == nil {
				alternativeName := strings.Title(pkgName) + inTypeName
				inst = getInstanceByName(alternativeName, inType)
			}
			if inst == nil {
				//log.Debug(inType.Kind())
				switch inType.Kind() {
				case reflect.Interface, reflect.Slice:
					injectByMethod = false
					break
				default:
					paramValue = reflect.New(inType)
					inst = paramValue.Interface()
					saveInstance(inTypeName, inst)
				}
			}

			if inst != nil {
				paramValue = reflect.ValueOf(inst)
			}
			inputs[i] = paramValue

			//log.Debugf("inType: %v, name: %v, instance: %v", inType, inTypeName, inst)
			//log.Debugf("kind: %v == %v, %v, %v ", obj.Kind(), reflect.Struct, paramValue.IsValid(), paramValue.CanSet())
			paramObject := reflect.Indirect(paramValue)
			if paramValue.IsValid() && paramObject.IsValid() && paramObject.Type() != obj.Type() && paramObject.Kind() == reflect.Struct {
				err = IntoObjectValue(paramValue)
			}
		}
		// finally call Init method to inject
		if injectByMethod {
			method.Func.Call(inputs)
		}
	}

	return err
}


