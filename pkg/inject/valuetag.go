package inject

import (
	"reflect"
	"strings"
	"strconv"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
)

type valueTag struct {
	BaseTag
}


func init() {
	AddTag(new(valueTag))
}


func (t *valueTag) IsSingleton() bool  {
	return false
}

func (t *valueTag) Decode(object reflect.Value, field reflect.StructField, tag string) (retVal interface{}) {
	if tag != "" {
		//log.Debug(valueTag)

		// check if filed type is slice
		switch reflector.GetKindByType(field.Type) {
		case reflect.Slice:
			retVal = t.replaceReferences(tag)
			if retVal == tag {
				retVal = strings.SplitN(tag, ",", -1)
			}
		case reflect.String:
			retVal = t.replaceReferences(tag)
		case reflect.Int:
			val, err := strconv.ParseInt(tag, 10, 64)
			if err == nil {
				retVal = val
			}
		case reflect.Uint:
			val, err := strconv.ParseUint(tag, 10, 64)
			if err == nil {
				retVal = val
			}
		case reflect.Float32:
			val, err := strconv.ParseFloat(tag, 32)
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
