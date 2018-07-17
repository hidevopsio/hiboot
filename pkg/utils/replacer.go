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
	"os"
	"regexp"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
)

// ParseVariables parse reference and env variables
func ParseVariables(src string, re *regexp.Regexp) [][]string {
	matches := re.FindAllStringSubmatch(src, -1)
	if matches == nil {
		return nil
	}
	return matches
}

// ReplaceStringVariables replace reference and env variables
func ReplaceStringVariables(source string, t interface{}) string {
	re := regexp.MustCompile(`\$\{(.*?)\}`)
	matches := ParseVariables(source, re)

	for _, match := range matches {
		varFullName := match[0]
		// replace references
		varName := match[1]
		vars := strings.SplitN(varName, ".", -1)
		refValue := ParseReferences(t, vars)
		// replace env
		envValue := os.Getenv(varName)
		if refValue != "" {
			source = strings.Replace(source, varFullName, refValue, -1)
		}
		source = strings.Replace(source, varFullName, envValue, -1)
	}
	return source
}

// GetFieldValue get filed value in reflected format
func GetFieldValue(f interface{}, name string) reflect.Value {
	r := reflect.ValueOf(f)
	fv := reflect.Indirect(r).FieldByName(name)

	return fv
}

// ParseReferences parse the variable references
func ParseReferences(st interface{}, varName []string) string {
	var parent interface{}
	parent = st
	for _, vn := range varName {
		capitalizedVarName := strings.Title(vn)
		field := GetFieldValue(parent, capitalizedVarName)

		k := reflector.GetKind(field)
		switch k {
		case reflect.String, reflect.Int, reflect.Float32:
			fv := fmt.Sprintf("%v", field.Interface())
			return fv
		case reflect.Invalid:
			return EmptyString
		default:
			// check if field is ptr
			parent = field.Addr().Interface()
		}

	}

	return EmptyString
}

// ReplaceMap replace references and env variables
func ReplaceMap(m map[string]interface{}, root interface{}) error {
	for k, v := range m {
		// log.Println(k, ": ", v)
		vt := reflect.TypeOf(v)
		if vt.Kind() == reflect.String {
			newStr := ReplaceStringVariables(v.(string), root)
			m[k] = newStr
		} else if vt.Kind() == reflect.Map {
			mv := v.(map[string]interface{})
			ReplaceMap(mv, root)
		}
	}
	return nil
}

// Replace given env and reference variables inside specific struct
func Replace(to interface{}, root interface{}) error {

	return reflector.ValidateReflectType(to, func(value *reflect.Value, reflectType reflect.Type, fieldSize int, isSlice bool) error {
		for i := 0; i < fieldSize; i++ {
			var dst reflect.Value
			if isSlice {
				dst = reflector.Indirect(value.Index(i))
				//log.Debug(dst.Interface())

				// TODO: refactoring below code
				dstType := dst.Type().Name()
				dstValue := dst.Interface()
				dv := fmt.Sprintf("%v", dstValue)
				//log.Debug("dstValue", dstValue)

				if dst.Kind() != reflect.String {
					child := dst.Addr().Interface()
					Replace(child, root)
				} else {
					if dv != "" && dstType == "string" && dst.IsValid() && dst.CanSet() {
						newStr := ReplaceStringVariables(dv, root)
						dst.SetString(newStr)
					}
				}
			} else {
				dst = reflector.Indirect(*value)
			}

			for _, field := range reflector.DeepFields(reflectType) {
				fieldName := field.Name
				//log.Debug("fieldName: ", fieldName)
				if dstField := dst.FieldByName(fieldName); dstField.IsValid() && dstField.CanSet() {
					fieldValue := dstField.Interface()
					//log.Debug("fieldValue: ", fieldValue)

					kind := dstField.Kind()
					switch kind {
					case reflect.String:
						fv := fmt.Sprintf("%v", fieldValue)
						newStr := ReplaceStringVariables(fv, root)
						dstField.SetString(newStr)
					case reflect.Map:
						mi := dstField.Interface()
						ReplaceMap(mi.(map[string]interface{}), root)
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
