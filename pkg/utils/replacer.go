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

package utils

import (
	"reflect"
	"fmt"
	"strings"
	"errors"
	"database/sql"
	"os"
	"github.com/hidevopsio/hiboot/pkg/log"
	"regexp"
)

func validate(toValue interface{}) (*reflect.Value, error) {

	to := indirect(reflect.ValueOf(toValue))

	if !to.CanAddr() {
		return nil, errors.New("value is unaddressable")
	}

	// Return is from value is invalid
	if !to.IsValid() {
		return nil, errors.New("value is not valid")
	}

	return &to, nil
}

// Copy copy things
func Copy(toValue interface{}, fromValue interface{}) (err error) {
	var (
		isSlice bool
		amount  = 1
		from    = indirect(reflect.ValueOf(fromValue))
		to      = indirect(reflect.ValueOf(toValue))
	)

	if !to.CanAddr() {
		return errors.New("copy to value is unaddressable")
	}

	// Return is from value is invalid
	if !from.IsValid() {
		return
	}

	// Just set it if possible to assign
	if from.Type().AssignableTo(to.Type()) {
		to.Set(from)
		return
	}

	fromType := indirectType(from.Type())
	toType := indirectType(to.Type())

	if fromType.Kind() != reflect.Struct || toType.Kind() != reflect.Struct {
		return
	}

	if to.Kind() == reflect.Slice {
		isSlice = true
		if from.Kind() == reflect.Slice {
			amount = from.Len()
		}
	}

	for i := 0; i < amount; i++ {
		var dest, source reflect.Value

		if isSlice {
			// source
			if from.Kind() == reflect.Slice {
				source = indirect(from.Index(i))
			} else {
				source = indirect(from)
			}

			// dest
			dest = indirect(reflect.New(toType).Elem())
		} else {
			source = indirect(from)
			dest = indirect(to)
		}

		// Copy from field to field or method
		for _, field := range DeepFields(fromType) {
			name := field.Name

			if fromField := source.FieldByName(name); fromField.IsValid() {
				// has field
				if toField := dest.FieldByName(name); toField.IsValid() {
					if toField.CanSet() {
						if !set(toField, fromField) {
							if err := Copy(toField.Addr().Interface(), fromField.Interface()); err != nil {
								return err
							}
						}
					}
				} else {
					// try to set to method
					var toMethod reflect.Value
					if dest.CanAddr() {
						toMethod = dest.Addr().MethodByName(name)
					} else {
						toMethod = dest.MethodByName(name)
					}

					if toMethod.IsValid() && toMethod.Type().NumIn() == 1 && fromField.Type().AssignableTo(toMethod.Type().In(0)) {
						toMethod.Call([]reflect.Value{fromField})
					}
				}
			}
		}

		// Copy from method to field
		for _, field := range DeepFields(toType) {
			name := field.Name

			var fromMethod reflect.Value
			if source.CanAddr() {
				fromMethod = source.Addr().MethodByName(name)
			} else {
				fromMethod = source.MethodByName(name)
			}

			if fromMethod.IsValid() && fromMethod.Type().NumIn() == 0 && fromMethod.Type().NumOut() == 1 {
				if toField := dest.FieldByName(name); toField.IsValid() && toField.CanSet() {
					values := fromMethod.Call([]reflect.Value{})
					if len(values) >= 1 {
						set(toField, values[0])
					}
				}
			}
		}

		if isSlice {
			if dest.Addr().Type().AssignableTo(to.Type().Elem()) {
				to.Set(reflect.Append(to, dest.Addr()))
			} else if dest.Type().AssignableTo(to.Type().Elem()) {
				to.Set(reflect.Append(to, dest))
			}
		}
	}
	return
}

func DeepFields(reflectType reflect.Type) []reflect.StructField {
	var fields []reflect.StructField

	if reflectType = indirectType(reflectType); reflectType.Kind() == reflect.Struct {
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

func indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

func indirectType(reflectType reflect.Type) reflect.Type {
	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
	}
	return reflectType
}

func set(to, from reflect.Value) bool {
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
			return set(to, from.Elem())
		} else {
			return false
		}
	}
	return true
}

func ParseVariables(src string, re *regexp.Regexp) [][]string {
	matches := re.FindAllStringSubmatch(src, -1)
	if matches == nil {
		return nil
	}
	return matches
}

func ReplaceStringVariables(source string, t interface{}) (string, error) {
	re := regexp.MustCompile(`\$\{(.*?)\}`)
	matches := ParseVariables(source, re)

	for _, match := range matches {
		varFullName := match[0]
		// replace references
		varName := match[1]
		vars := strings.SplitN(varName, ".", -1)
		refValue, err := ParseReferences(t, vars)
		if err != nil {
			return source, err
		}
		// replace env
		envValue := os.Getenv(varName)
		if refValue != "" {
			source = strings.Replace(source, varFullName, refValue, -1)
		}
		source = strings.Replace(source, varFullName, envValue, -1)
	}
	return source, nil
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

		k := getKind(field)
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

func getKind(val reflect.Value) reflect.Kind {

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
	v, err := validate(obj)
	if err != nil {
		return err
	}

	t := indirectType(v.Type())

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

func Replace(to interface{}, root interface{}) error {

	return ValidateReflectType(to, func(value *reflect.Value, reflectType reflect.Type, fieldSize int, isSlice bool) error {
		for i := 0; i < fieldSize; i++ {
			var dst reflect.Value
			if isSlice {
				//dst = indirect(reflect.New(toType).Elem())
				if value.Kind() == reflect.Slice {
					dst = indirect(value.Index(i))
					//log.Debug(dst.Interface())

					// TODO: refactoring below code
					dstType := dst.Type().Name()
					dstValue := dst.Interface()
					//log.Debug(dstType)
					dv := fmt.Sprintf("%v", dstValue)

					if dst.Kind() != reflect.String {
						child := dst.Addr().Interface()
						Replace(child, root)
					} else {
						if dv != "" {

							if dstType == "string" && dst.IsValid() && dst.CanSet() {
								newStr, err := ReplaceStringVariables(dv, root)
								if err != nil {
									return err
								}
								dst.SetString(newStr)
							} else {
								log.Error("")
							}

						}
					}
				} else {
					dst = indirect(*value)
				}
			} else {
				dst = indirect(*value)
			}

			for _, field := range DeepFields(reflectType) {
				fieldName := field.Name
				//log.Debug("fieldName: ", fieldName)
				if dstField := dst.FieldByName(fieldName); dstField.IsValid() && dstField.CanSet() {
					fieldValue := dstField.Interface()
					//log.Debug("fieldValue: ", fieldValue)

					kind := dstField.Kind()
					switch kind {
					case reflect.String:
						fv := fmt.Sprintf("%v", fieldValue)
						newStr, err := ReplaceStringVariables(fv, root)
						if err != nil {
							return err
						}
						dstField.SetString(newStr)
					default:
						//log.Debug(fieldName, " is a ", kind)
						Replace(dstField.Addr().Interface(), root)
					}
				}
			}
		}
		return nil
	})
}
