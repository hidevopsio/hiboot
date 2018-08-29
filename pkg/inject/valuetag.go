package inject

import (
	"reflect"
	"strings"
	"strconv"
)

type valueTag struct {
	BaseTag
}

func init() {
	AddTag(new(valueTag))
}

func (t *valueTag) IsSingleton() bool {
	return false
}

func (t *valueTag) Decode(object reflect.Value, field reflect.StructField, tag string) (retVal interface{}) {
	if tag != "" {
		//log.Debug(valueTag)

		// check if filed type is slice
		kind := field.Type.Kind()
		switch kind {
		case reflect.Slice:
			retVal = t.replaceReferences(tag)
			if retVal == tag {
				retVal = strings.SplitN(tag, ",", -1)
			}
		case reflect.String:
			retVal = t.replaceReferences(tag)
		case reflect.Int:
			val, err := strconv.ParseInt(tag, 10, 32)
			if err == nil {
				retVal = int(val)
			}
		case reflect.Int8:
			val, err := strconv.ParseInt(tag, 10, 8)
			if err == nil {
				retVal = int8(val)
			}

		case reflect.Int16:
			val, err := strconv.ParseInt(tag, 10, 16)
			if err == nil {
				retVal = int16(val)
			}

		case reflect.Int32:
			val, err := strconv.ParseInt(tag, 10, 32)
			if err == nil {
				retVal = int32(val)
			}

		case reflect.Int64:
			val, err := strconv.ParseInt(tag, 10, 64)
			if err == nil {
				retVal = int64(val)
			}

		case reflect.Uint:
			val, err := strconv.ParseInt(tag, 10, 32)
			if err == nil {
				retVal = uint(val)
			}
		case reflect.Uint8:
			val, err := strconv.ParseInt(tag, 10, 8)
			if err == nil {
				retVal = uint8(val)
			}

		case reflect.Uint16:
			val, err := strconv.ParseInt(tag, 10, 16)
			if err == nil {
				retVal = uint16(val)
			}

		case reflect.Uint32:
			val, err := strconv.ParseInt(tag, 10, 32)
			if err == nil {
				retVal = uint32(val)
			}

		case reflect.Uint64:
			val, err := strconv.ParseInt(tag, 10, 64)
			if err == nil {
				retVal = uint64(val)
			}
		case reflect.Float32:
			val, err := strconv.ParseFloat(tag, 32)
			if err == nil {
				retVal = float32(val)
			}
		case reflect.Float64:
			val, err := strconv.ParseFloat(tag, 64)
			if err == nil {
				retVal = val
			}
		case reflect.Bool:
			val, err := strconv.ParseBool(tag)
			if err == nil {
				retVal = val
			}
		}

	}
	return retVal
}
