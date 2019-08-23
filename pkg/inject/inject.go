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
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/factory"
	"hidevops.io/hiboot/pkg/inject/annotation"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/utils/reflector"
	"hidevops.io/hiboot/pkg/utils/str"
	"reflect"
)

const (
	value          = "value"
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

	// ErrInvalidMethod the function is invalid
	ErrInvalidMethod = errors.New("[inject] invalid method")

	// ErrFactoryIsNil factory is invalid
	ErrFactoryIsNil = errors.New("[inject] factory is nil")

	tagsContainer []Tag

	//instancesMap cmap.ConcurrentMap
	//appFactory factory.ConfigurableFactory
)

// Inject is the interface for inject tag
type Inject interface {
	DefaultValue(object interface{}) error
	IntoObject(object interface{}) error
	IntoObjectValue(object reflect.Value, property string, tags ...Tag) error
	IntoMethod(object interface{}, m interface{}) (retVal interface{}, err error)
	IntoFunc(object interface{}) (retVal interface{}, err error)
}

type inject struct {
	factory factory.InstantiateFactory
}

// NewInject is the constructor of inject
func NewInject(factory factory.InstantiateFactory) Inject {
	return &inject{factory: factory}
}

// InitTag init tag implements
func InitTag(tag Tag) (t Tag) {
	f, ok := reflector.GetEmbeddedFieldType(tag, new(at.Tag))
	if ok {
		tv := reflector.Indirect(reflect.ValueOf(tag))
		v, ok := f.Tag.Lookup(value)
		if ok {
			fo := tv.FieldByName(f.Name)
			fo.Set(reflect.ValueOf(at.Tag(v)))
			t = tag
		}
	}
	return
}

// AddTag add new tag
func AddTag(tag Tag) {
	if tag != nil {
		tagsContainer = append(tagsContainer, InitTag(tag))
	}
}

func (i *inject) getInstanceByName(name string, typ reflect.Type) (inst interface{}) {
	n := reflector.GetLowerCamelFullNameByType(typ)
	inst = i.factory.GetInstance(n)
	return
}

// DefaultValue injects instance into the tagged field with `inject:"instanceName"`
func (i *inject) DefaultValue(object interface{}) error {
	return i.IntoObjectValue(reflect.ValueOf(object), "", InitTag(new(defaultTag)))
}

// IntoObject injects instance into the tagged field with `inject:"instanceName"`
func (i *inject) IntoObject(object interface{}) (err error) {
	err = annotation.InjectIntoObject(object)
	if err != nil {
		log.Warn("annotation.InjectIntoObject() return failed")
	}
	err = i.IntoObjectValue(reflect.ValueOf(object), "")
	return
}

// IntoObjectValue injects instance into the tagged field with `inject:"instanceName"`
func (i *inject) IntoObjectValue(object reflect.Value, property string, tags ...Tag) error {
	var err error

	//// TODO refactor IntoObject
	//if appFactory == nil {
	//	return ErrSystemConfiguration
	//}

	obj := reflector.Indirect(object)
	if obj.Kind() != reflect.Struct {
		log.Debugf("[inject] ignore object: %v, kind: %v", object, obj.Kind())
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
		prop := str.ToLowerCamel(f.Name)
		if property != "" {
			prop = property + "." + prop
		}
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

		// TODO: should inject annotation
		//ok := reflector.HasEmbeddedFieldType(ft, at.Annotation{})
		//if ok {
		//	// parse annotation
		//	log.Debug("found annotation")
		//}

		// TODO: assume that the f.Name of value and inject tag is not the same
		injectedObject = i.getInstanceByName(f.Name, f.Type)
		if injectedObject == nil {
			for _, tagImpl := range targetTags {
				tagImpl.Init(i.factory)
				injectedObject = tagImpl.Decode(object, f, prop)
				if injectedObject != nil {
					break
				}
			}
		}

		if injectedObject != nil && fieldObj.CanSet() {
			fov := reflect.ValueOf(injectedObject)
			switch injectedObject.(type) {
			//case []string:
			case []interface{}:
				switch f.Type.Elem().Kind() {
				case reflect.String:
					var sv []string
					src := injectedObject.([]interface{})
					for _, elm := range src {
						sv = append(sv, elm.(string))
					}
					fov = reflect.ValueOf(sv)
					log.Debugf("Injected %v.(%v) into %v.%v", injectedObject, fov.Type(), obj.Type(), f.Name)
				}
			}
			log.Debugf("Injected %v.(%v) into %v.%v", injectedObject, fov.Type(), obj.Type(), f.Name)
			fieldObj.Set(fov)
		}

		//log.Debugf("- kind: %v, %v, %v, %v", obj.Kind(), object.Type(), fieldObj.Type(), f.Name)
		//log.Debugf("isValid: %v, canSet: %v", fieldObj.IsValid(), fieldObj.CanSet())
		filedObject := reflect.Indirect(fieldObj)
		filedKind := filedObject.Kind()
		canNested := filedKind == reflect.Struct
		if canNested && fieldObj.IsValid() && fieldObj.CanSet() && filedObject.Type() != obj.Type() {
			err = i.IntoObjectValue(fieldObj, prop, tags...)
		}
	}
	return err
}

func (i *inject) parseFuncOrMethodInput(inType reflect.Type) (paramValue reflect.Value, ok bool) {
	inType = reflector.IndirectType(inType)
	inTypeName := inType.Name()
	inst := i.getInstanceByName(inTypeName, inType)
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
			i.factory.SetInstance(inst)
		}
	}

	if inst != nil {
		paramValue = reflect.ValueOf(inst)
	}
	return
}

// IntoFunc inject object into func and return instance
func (i *inject) IntoFunc(object interface{}) (retVal interface{}, err error) {
	fn := reflect.ValueOf(object)
	if fn.Kind() == reflect.Func {
		numIn := fn.Type().NumIn()
		inputs := make([]reflect.Value, numIn)
		// TODO: should load function inputs when resolving dependencies to improve performance
		for n := 0; n < numIn; n++ {
			fnInType := fn.Type().In(n)
			//expectedTypName := reflector.GetLowerCamelFullNameByType(fnInType)
			//log.Debugf("expected: %v", expectedTypName)
			val, ok := i.parseFuncOrMethodInput(fnInType)
			if ok {

				inputs[n] = val
				//log.Debugf("Injected %v into func parameter %v", val, fnInType)
			} else {
				return nil, fmt.Errorf("%v is not injected", fnInType.Name())
			}

			paramValue := reflect.Indirect(val)
			if val.IsValid() && paramValue.IsValid() && paramValue.Kind() == reflect.Struct {
				err = i.IntoObjectValue(val, "")
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
//TODO: IntoMethod or IntoFunc should accept metaData, because it contains dependencies
func (i *inject) IntoMethod(object interface{}, m interface{}) (retVal interface{}, err error) {
	if object != nil && m != nil {
		switch m.(type) {
		case reflect.Method:
			method := m.(reflect.Method)
			numIn := method.Type.NumIn()
			inputs := make([]reflect.Value, numIn)
			inputs[0] = reflect.ValueOf(object)
			for n := 1; n < numIn; n++ {
				fnInType := method.Type.In(n)
				val, ok := i.parseFuncOrMethodInput(fnInType)
				if ok {
					inputs[n] = val
				} else {
					return nil, fmt.Errorf("%v is not injected", fnInType.Name())
				}

				paramObject := reflect.Indirect(val)
				if val.IsValid() && paramObject.IsValid() && paramObject.Kind() == reflect.Struct {
					err = i.IntoObjectValue(val, "")
				}
			}
			results := method.Func.Call(inputs)
			if len(results) != 0 {
				retVal = results[0].Interface()
				return
			}
		}
	}
	err = ErrInvalidMethod
	return
}
