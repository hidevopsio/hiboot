package inject

import (
	"reflect"
	"strings"
	"strconv"
	"github.com/hidevopsio/hiboot/pkg/log"
)

type defaultTag struct {
	BaseTag
}

var bitNumLut = map[string]int{
	"int8": 8,
	"int16": 16,
	"int32": 32,
	"int64": 64,
	"int": 32,
	"uint8": 8,
	"uint16": 16,
	"uint32": 32,
	"uint64": 64,
	"uint": 32,
	"float32": 32,
	"float64": 64,
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
		log.Debugf("field: %v, value: %v", field.Name, fieldVal)

		// check if filed type is slice
		kind := field.Type.Kind()
		switch kind {
		case reflect.Slice:
			if fieldVal == nil {
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
