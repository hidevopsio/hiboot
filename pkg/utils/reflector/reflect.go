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
	"database/sql"
	"errors"
	"strings"
	"fmt"
)

func NewReflectType(st interface{}) interface{} {
	ct := reflect.TypeOf(st)
	co := reflect.New(ct)
	cp := co.Elem().Addr().Interface()
	return cp
}


func Set(to, from reflect.Value) bool {
	if from.IsValid() {
		if to.Kind() == reflect.Ptr {
			//set `to` to nil if from is nil
			if from.Kind() == reflect.Ptr && from.IsNil() {
				to.Set(reflect.Zero(to.Type()))
				return true
			} else if to.IsNil() {
				to.Set(reflect.New(to.Type().Elem()))
			}
			to = to.Elem()
		}

		if from.Type().ConvertibleTo(to.Type()) {
			to.Set(from.Convert(to.Type()))
		} else if scanner, ok := to.Addr().Interface().(sql.Scanner); ok {
			err := scanner.Scan(from.Interface())
			if err != nil {
				return false
			}
		} else if from.Kind() == reflect.Ptr {
			return Set(to, from.Elem())
		} else {
			return false
		}
	}
	return true
}


func Validate(toValue interface{}) (*reflect.Value, error) {

	to := Indirect(reflect.ValueOf(toValue))

	if !to.CanAddr() {
		return nil, errors.New("value is unaddressable")
	}

	// Return is from value is invalid
	if !to.IsValid() {
		return nil, errors.New("value is not valid")
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

func ParseReferences(st interface{}, varName []string) (string, error) {
	var parent interface{}
	parent = st
	for _, vn := range varName {
		capitalizedVarName := strings.Title(vn)
		field := GetFieldValue(parent, capitalizedVarName)

		k := GetKind(field)
		switch k {
		case reflect.String:
			fv := fmt.Sprintf("%v", field.Interface())
			return fv, nil
		case reflect.Int:
			fv := fmt.Sprintf("%v", field.Interface())
			return fv, nil
		case reflect.Invalid:
			return "", nil
		default:
			// check if field is ptr
			parent = field.Addr().Interface()
		}

	}

	return "", nil
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