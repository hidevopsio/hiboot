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

package str

import (
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

const EmptyString = ""

// UpperFirst upper case first character of specific string
func UpperFirst(str string) string {
	return strings.Title(str)
}

// LowerFirst lower case first character of specific string
func LowerFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return EmptyString
}

// StringInSlice check if specific string is in slice
func InSlice(a string, list []string) bool {

	var retVal bool

	for _, b := range list {
		if b == a {
			retVal = true
			break
		}
	}
	return retVal
}

func Convert(src string, kind reflect.Kind) (retVal interface{}) {
	switch kind {
	case reflect.Slice:
		retVal = strings.SplitN(src, ",", -1)
	case reflect.String:
		retVal = src
	case reflect.Int:
		val, err := strconv.ParseInt(src, 10, 32)
		if err == nil {
			retVal = int(val)
		} else {
			retVal = int(0)
		}
	case reflect.Int8:
		val, err := strconv.ParseInt(src, 10, 8)
		if err == nil {
			retVal = int8(val)
		} else {
			retVal = int8(0)
		}

	case reflect.Int16:
		val, err := strconv.ParseInt(src, 10, 16)
		if err == nil {
			retVal = int16(val)
		} else {
			retVal = int16(0)
		}

	case reflect.Int32:
		val, err := strconv.ParseInt(src, 10, 32)
		if err == nil {
			retVal = int32(val)
		} else {
			retVal = int32(0)
		}

	case reflect.Int64:
		val, err := strconv.ParseInt(src, 10, 64)
		if err == nil {
			retVal = int64(val)
		} else {
			retVal = int64(0)
		}

	case reflect.Uint:
		val, err := strconv.ParseInt(src, 10, 32)
		if err == nil {
			retVal = uint(val)
		} else {
			retVal = uint(0)
		}
	case reflect.Uint8:
		val, err := strconv.ParseInt(src, 10, 8)
		if err == nil {
			retVal = uint8(val)
		} else {
			retVal = uint8(0)
		}

	case reflect.Uint16:
		val, err := strconv.ParseInt(src, 10, 16)
		if err == nil {
			retVal = uint16(val)
		} else {
			retVal = uint16(0)
		}

	case reflect.Uint32:
		val, err := strconv.ParseInt(src, 10, 32)
		if err == nil {
			retVal = uint32(val)
		} else {
			retVal = uint32(0)
		}

	case reflect.Uint64:
		val, err := strconv.ParseInt(src, 10, 64)
		if err == nil {
			retVal = uint64(val)
		} else {
			retVal = uint64(0)
		}
	case reflect.Float32:
		val, err := strconv.ParseFloat(src, 32)
		if err == nil {
			retVal = float32(val)
		} else {
			retVal = float32(0.0)
		}
	case reflect.Float64:
		val, err := strconv.ParseFloat(src, 64)
		if err == nil {
			retVal = val
		} else {
			retVal = float64(0.0)
		}
	case reflect.Bool:
		val, err := strconv.ParseBool(src)
		if err == nil {
			retVal = val
		} else {
			retVal = false
		}
	}
	return
}
