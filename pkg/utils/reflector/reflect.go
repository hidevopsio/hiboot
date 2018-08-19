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

package reflector

import (
	"reflect"
	"errors"
	"strings"
	"unicode"
)

var InvalidInputError = errors.New("input is invalid")
var InvalidMethodError = errors.New("method is invalid")
var FieldCanNotBeSetError = errors.New("field can not be set")

func NewReflectType(st interface{}) interface{} {
	ct := reflect.TypeOf(st)
	co := reflect.New(ct)
	cp := co.Elem().Addr().Interface()
	return cp
}

func Validate(toValue interface{}) (*reflect.Value, error) {

	to := Indirect(reflect.ValueOf(toValue))

	// Return is from value is invalid
	if !to.IsValid() {
		return nil, errors.New("value is not valid")
	}

	if !to.CanAddr() {
		return nil, errors.New("value is unaddressable")
	}

	return &to, nil
}

func DeepFields(reflectType reflect.Type) []reflect.StructField {
	var fields []reflect.StructField

	if reflectType = IndirectType(reflectType); reflectType.Kind() == reflect.Struct {
		for i := 0; i < reflectType.NumField(); i++ {
			v := reflectType.Field(i)
			if v.Anonymous {
				fields = append(fields, DeepFields(v.Type)...)
			} else {
				fields = append(fields, v)
			}
		}
	}

	return fields
}

func Indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

func IndirectType(reflectType reflect.Type) reflect.Type {
	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
	}
	return reflectType
}


func GetFieldValue(f interface{}, name string) reflect.Value {
	r := reflect.ValueOf(f)
	fv := reflect.Indirect(r).FieldByName(name)

	return fv
}


func SetFieldValue(object interface{}, name string, value interface{}) error  {

	obj := Indirect(reflect.ValueOf(object))

	if ! obj.IsValid()  {
		return InvalidInputError
	}

	if obj.Kind() != reflect.Struct {
		return InvalidInputError
	}

	fieldObj := obj.FieldByName(name)

	if ! fieldObj.CanSet() {
		return FieldCanNotBeSetError
	}

	fov := reflect.ValueOf(value)
	fieldObj.Set(fov)

	//log.Debugf("Set %v.(%v) into %v.%v", value, fov.Type(), obj.Type(), name)
	return nil
}

func GetKind(val reflect.Value) reflect.Kind {

	// Capture the value's Kind.
	kind := val.Kind()

	// Check each condition until a case is true.
	switch {

	case kind >= reflect.Int && kind <= reflect.Int64:
		return reflect.Int

	case kind >= reflect.Uint && kind <= reflect.Uint64:
		return reflect.Uint

	case kind >= reflect.Float32 && kind <= reflect.Float64:
		return reflect.Float32

	default:
		return kind
	}
}

func ValidateReflectType(obj interface{}, callback func(value *reflect.Value, reflectType reflect.Type, fieldSize int, isSlice bool) error) error {
	v, err := Validate(obj)
	if err != nil {
		return err
	}

	t := IndirectType(v.Type())

	isSlice := false
	fieldSize := 1
	if v.Kind() == reflect.Slice {
		isSlice = true
		fieldSize = v.Len()
	}

	if callback != nil {
		return callback(v, t, fieldSize, isSlice)
	}

	return err
}

func GetName(data interface{}) (string, error)  {
	dv := Indirect(reflect.ValueOf(data))

	// Return is from value is invalid
	if !dv.IsValid() {
		return "", InvalidInputError
	}
	name := dv.Type().Name()
	return name, nil
}

func GetLowerCaseObjectName(data interface{}) (string, error) {
	name, err := GetName(data)
	name = strings.ToLower(name)
	return name, err
}

func HasField(object interface{}, name string) bool  {
	r := reflect.ValueOf(object)
	fv := reflect.Indirect(r).FieldByName(name)

	return fv.IsValid()
}

func CallMethodByName(object interface{}, name string, args ...interface{}) (interface{}, error)  {
	objVal := reflect.ValueOf(object)
	method, ok := objVal.Type().MethodByName(name)
	if ok {
		numIn := method.Type.NumIn()
		inputs := make([]reflect.Value, numIn)
		inputs[0] = objVal
		for i, arg := range args {
			inputs[i + 1] = reflect.ValueOf(arg)
		}
		results := method.Func.Call(inputs)
		if len(results) != 0 {
			return results[0].Interface(), nil
		} else {
			return nil, nil
		}
	}
	return nil, InvalidMethodError
}

func HasEmbeddedField(object interface{}, name string) bool {
	typ := IndirectType(reflect.TypeOf(object))
	field, ok := typ.FieldByName(name)
	return field.Anonymous && ok
}


// LowerFirst lower case first character of specific string
func lowerFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

// ParseObjectName e.g. ExampleObject => example
func ParseObjectName(cmd interface{}, eliminator string) string {
	name, err := GetName(cmd)
	if err == nil {
		name = strings.Replace(name, eliminator, "", -1)
		name = lowerFirst(name)
	}
	return name
}

func GetPkgPath(object interface{}) string {
	objType := IndirectType(reflect.TypeOf(object))
	return objType.PkgPath()
}
