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
	"strings"
	"strconv"
)

type defaultTag struct {
	BaseTag
}

func init() {
	AddTag(new(defaultTag))
}

func (t *defaultTag) IsSingleton() bool  {
	return false
}

func (t *defaultTag) Decode(object reflect.Value, field reflect.StructField, tag string) (retVal interface{}) {
	if tag != "" {
		fieldVal := object.FieldByName(field.Name).Interface()
		//log.Debugf("field: %v, value: %v", field.Name, fieldVal)

		// check if filed type is slice
		kind := field.Type.Kind()
		switch kind {
		case reflect.Slice:
			if len(fieldVal.([]string)) == 0  {
				retVal = t.replaceReferences(tag)
				if retVal == tag {
					retVal = strings.SplitN(tag, ",", -1)
				}
			}
		case reflect.String:
			if fieldVal.(string) == "" {
				retVal = t.replaceReferences(tag)
			}
		case reflect.Int:
			if fieldVal.(int) == 0 {
				val, err := strconv.ParseInt(tag, 10, 32)
				if err == nil {
					retVal = int(val)
				}
			}
		case reflect.Int8:
			if fieldVal.(int8) == 0 {
				val, err := strconv.ParseInt(tag, 10, 8)
				if err == nil {
					retVal = int8(val)
				}
			}

		case reflect.Int16:
			if fieldVal.(int16) == 0 {
				val, err := strconv.ParseInt(tag, 10, 16)
				if err == nil {
					retVal = int16(val)
				}
			}

		case reflect.Int32:
			if fieldVal.(int32) == 0 {
				val, err := strconv.ParseInt(tag, 10, 32)
				if err == nil {
					retVal = int32(val)
				}
			}

		case reflect.Int64:
			if fieldVal.(int64) == 0 {
				val, err := strconv.ParseInt(tag, 10, 64)
				if err == nil {
					retVal = int64(val)
				}
			}

		case reflect.Uint:
			if fieldVal.(uint) == 0 {
				val, err := strconv.ParseInt(tag, 10, 32)
				if err == nil {
					retVal = uint(val)
				}
			}
		case reflect.Uint8:
			if fieldVal.(uint8) == 0 {
				val, err := strconv.ParseInt(tag, 10, 8)
				if err == nil {
					retVal = uint8(val)
				}
			}

		case reflect.Uint16:
			if fieldVal.(uint16) == 0 {
				val, err := strconv.ParseInt(tag, 10, 16)
				if err == nil {
					retVal = uint16(val)
				}
			}

		case reflect.Uint32:
			if fieldVal.(uint32) == 0 {
				val, err := strconv.ParseInt(tag, 10, 32)
				if err == nil {
					retVal = uint32(val)
				}
			}

		case reflect.Uint64:
			if fieldVal.(uint64) == 0 {
				val, err := strconv.ParseInt(tag, 10, 64)
				if err == nil {
					retVal = uint64(val)
				}
			}
		case reflect.Float32:
			if fieldVal.(float32) == 0.0 {
				val, err := strconv.ParseFloat(tag, 32)
				if err == nil {
					retVal = float32(val)
				}
			}
		case reflect.Float64:
			if fieldVal.(float64) == 0.0 {
				val, err := strconv.ParseFloat(tag, 64)
				if err == nil {
					retVal = val
				}
			}
		case reflect.Bool:
			if fieldVal.(bool) == false {
				val, err := strconv.ParseBool(tag)
				if err == nil {
					retVal = val
				}
			}
		}
	}
	return retVal
}
