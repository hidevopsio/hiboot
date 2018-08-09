package inject

import (
	"reflect"
	"github.com/hidevopsio/hiboot/pkg/utils/mapstruct"
)

type injectTag struct {
	BaseTag
}

func init() {
	AddTag("inject", new(injectTag))
}


func (t *injectTag) IsSingleton() bool  {
	return true
}

func (t *injectTag) Decode(object reflect.Value, field reflect.StructField, tag string) (retVal interface{})  {
	properties := t.ParseProperties(tag)

	// first, find if object is already instantiated
	if field.Type.Kind() == reflect.Ptr {
		// if object is not exist, then instantiate new object
		// parse tag and instantiate filed
		ft := field.Type
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		o := reflect.New(ft)
		retVal = o.Interface()
		// inject field value
		if len(properties) != 0 {
			mapstruct.Decode(retVal, properties)
		}
	} else {
		return
	}
	return retVal
}