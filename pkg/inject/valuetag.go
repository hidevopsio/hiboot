package inject

import (
	"reflect"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
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
		needConvert := true
		switch kind {
		case reflect.Slice:
			retVal = t.replaceReferences(tag)
			if retVal != tag {
				needConvert = false
			}
		case reflect.String:
			retVal = t.replaceReferences(tag)
			needConvert = false
		}

		if needConvert {
			retVal = str.Convert(tag, kind)
		}
	}
	return retVal
}
