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

// Package reflector provides utilities for reflection
package reflector

import (
	"errors"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
	"reflect"
	"strings"
)

var (
	// ErrInvalidInput means that the input is invalid
	ErrInvalidInput = errors.New("input is invalid")

	// ErrInvalidMethod means that the method is invalid
	ErrInvalidMethod = errors.New("method is invalid")

	// ErrInvalidFunc means that the func is invalid
	ErrInvalidFunc = errors.New("func is invalid")

	// ErrFieldCanNotBeSet means that the field can not be set
	ErrFieldCanNotBeSet = errors.New("field can not be set")
)

// NewReflectType create instance by tyep
func NewReflectType(st interface{}) interface{} {
	ct := reflect.TypeOf(st)
	co := reflect.New(ct)
	cp := co.Elem().Addr().Interface()
	return cp
}

// Validate validate value
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

// DeepFields iterate struct field
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

// Indirect get indirect value
func Indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

// IndirectType get indirect type
func IndirectType(reflectType reflect.Type) reflect.Type {
	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
	}
	return reflectType
}

// GetFieldValue get field value
func GetFieldValue(f interface{}, name string) reflect.Value {
	r := reflect.ValueOf(f)
	fv := reflect.Indirect(r).FieldByName(name)

	return fv
}

// SetFieldValue set field value
func SetFieldValue(object interface{}, name string, value interface{}) error {

	obj := Indirect(reflect.ValueOf(object))

	if !obj.IsValid() {
		return ErrInvalidInput
	}

	if obj.Kind() != reflect.Struct {
		return ErrInvalidInput
	}

	fieldObj := obj.FieldByName(name)

	if !fieldObj.CanSet() {
		return ErrFieldCanNotBeSet
	}

	fov := reflect.ValueOf(value)
	fieldObj.Set(fov)

	//log.Debugf("Set %v.(%v) into %v.%v", value, fov.Type(), obj.Type(), name)
	return nil
}

// GetKind get kind
func GetKind(kind reflect.Kind) reflect.Kind {

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

// GetKindByValue get kind by value
func GetKindByValue(val reflect.Value) reflect.Kind {
	return GetKind(val.Kind())
}

// GetKindByType get kind by type
func GetKindByType(typ reflect.Type) reflect.Kind {
	return GetKind(typ.Kind())
}

// ValidateReflectType validate reflect type and iterate all fields
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

// GetType get object data type
func GetType(data interface{}) (typ reflect.Type, err error) {
	dv := Indirect(reflect.ValueOf(data))

	// Return is from value is invalid
	if !dv.IsValid() {
		err = ErrInvalidInput
		return
	}
	typ = dv.Type()

	//log.Debugf("%v %v %v %v %v", dv, typ, typ.String(), typ.Name(), typ.PkgPath())
	return
}

// GetName get object name
func GetName(data interface{}) (name string, err error) {

	typ, err := GetType(data)
	if err == nil {
		name = typ.Name()
	}
	return
}

// GetLowerCaseObjectName get lower case object name
func GetLowerCamelName(object interface{}) (string, error) {
	name, err := GetName(object)
	name = str.ToLowerCamel(name)
	return name, err
}

// HasField check if has specific field
func HasField(object interface{}, name string) bool {
	r := reflect.ValueOf(object)
	fv := reflect.Indirect(r).FieldByName(name)

	return fv.IsValid()
}

// CallMethodByName call method
func CallMethodByName(object interface{}, name string, args ...interface{}) (interface{}, error) {
	objVal := reflect.ValueOf(object)
	method, ok := objVal.Type().MethodByName(name)
	if ok {
		numIn := method.Type.NumIn()
		inputs := make([]reflect.Value, numIn)
		inputs[0] = objVal
		for i, arg := range args {
			inputs[i+1] = reflect.ValueOf(arg)
		}
		results := method.Func.Call(inputs)
		if len(results) != 0 {
			return results[0].Interface(), nil
		}
		return nil, nil
	}
	return nil, ErrInvalidMethod
}

// CallFunc call function
func CallFunc(object interface{}, args ...interface{}) (interface{}, error) {
	fn := reflect.ValueOf(object)
	if fn.Kind() == reflect.Func {
		numIn := fn.Type().NumIn()
		inputs := make([]reflect.Value, numIn)
		for i, arg := range args {
			inputs[i] = reflect.ValueOf(arg)
		}
		results := fn.Call(inputs)
		if len(results) != 0 {
			return results[0].Interface(), nil
		} else {
			return nil, nil
		}
	}
	return nil, ErrInvalidFunc
}

// HasEmbeddedField check if has embedded fieled
func HasEmbeddedField(object interface{}, name string) bool {
	//log.Debugf("HasEmbeddedField: %v", name)
	typ := IndirectType(reflect.TypeOf(object))
	if typ.Kind() != reflect.Struct {
		return false
	}
	field, ok := typ.FieldByName(name)
	return field.Anonymous && ok
}

// GetEmbeddedInterfaceFieldByType get embedded interface field by type
func GetEmbeddedFieldByType(typ reflect.Type, kind ...reflect.Kind) (field reflect.StructField) {
	expectedKind := reflect.Interface
	if len(kind) > 0 {
		expectedKind = kind[0]
	}
	if typ.Kind() == reflect.Struct {
		for i := 0; i < typ.NumField(); i++ {
			v := typ.Field(i)
			if v.Anonymous {
				if v.Type.Kind() == expectedKind {
					return v
				} else {
					return GetEmbeddedFieldByType(v.Type)
				}
			}
		}
	}
	return
}

// GetEmbeddedInterfaceField get embedded interface field
func GetEmbeddedField(object interface{}, dataTypes ...reflect.Kind) (field reflect.StructField) {
	if object == nil {
		return
	}

	typ, ok := GetFuncOutType(object)
	if ok {
		return GetEmbeddedFieldByType(typ, dataTypes...)
	}

	typ = IndirectType(reflect.TypeOf(object))
	return GetEmbeddedFieldByType(typ, dataTypes...)
}

// FindEmbeddedFieldTag find embedded field tag
func FindEmbeddedFieldTag(object interface{}, name string) (t string, ok bool) {
	f := GetEmbeddedField(object)
	t, ok = f.Tag.Lookup(name)
	return
}

// ParseObjectName e.g. ExampleObject => example
func ParseObjectName(obj interface{}, eliminator string) string {
	name, err := GetName(obj)
	if err == nil {
		name = strings.Replace(name, eliminator, "", -1)
		name = str.LowerFirst(name)
	}
	return name
}

// ParseObjectName e.g. ExampleObject => example
func ParseObjectPkgName(obj interface{}) string {

	typ := IndirectType(reflect.TypeOf(obj))
	name := io.DirName(typ.PkgPath())

	return name
}

// GetPkgPath get the package patch
func GetPkgPath(object interface{}) string {
	objType := IndirectType(reflect.TypeOf(object))
	return objType.PkgPath()
}

// GetFuncOutType get the function output data type
func GetFuncOutType(object interface{}) (typ reflect.Type, ok bool) {

	obj := reflect.ValueOf(object)
	t := obj.Type()
	typName := t.Name()
	typKind := t.Kind()
	//log.Debugf("type: %v type name: %v, kind: %v", t, typName, typKind)
	if typKind == reflect.Func {
		numOut := obj.Type().NumOut()
		if numOut > 0 {
			typ = IndirectType(obj.Type().Out(0))
			ok = true
		}
	} else if typKind == reflect.Struct && typName == "Method" {
		method := object.(reflect.Method)
		methodTyp := method.Func.Type()
		numOut := methodTyp.NumOut()
		if numOut > 0 {
			typ = IndirectType(methodTyp.Out(0))
			ok = true
		}
	}

	return
}

// GetFullName get the object name with package name, e.g. pkg.Object
func GetFullName(object interface{}) (name string) {

	pn, n := GetPkgAndName(object)
	name = pn + "." + n

	return
}

// GetLowerCamelFullName get the object name with package name, e.g. pkg.objectName
func GetLowerCamelFullName(object interface{}) (name string) {

	pn, n := GetPkgAndName(object)
	name = pn + "." + str.ToLowerCamel(n)

	return
}

// GetPkgAndName get the package name and the object name with, e.g. pkg, Object
func GetPkgAndName(object interface{}) (pkgName, name string) {

	typ, ok := GetFuncOutType(object)
	if ok {
		pkgName = io.DirName(typ.PkgPath())
		name = typ.Name()
		return
	}

	name, err := GetName(object)
	if err == nil {
		pkgName = ParseObjectPkgName(object)
	}
	return
}

// GetFullNameByType get the object name with package name by type, e.g. pkg.Object
func GetLowerCamelFullNameByType(objType reflect.Type) (name string) {
	indTyp := IndirectType(objType)
	depPkgName := io.DirName(indTyp.PkgPath())
	name = depPkgName + "." + str.ToLowerCamel(indTyp.Name())
	return
}
