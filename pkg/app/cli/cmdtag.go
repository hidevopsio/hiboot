package cli

import (
	"github.com/hidevopsio/hiboot/pkg/inject"
	"reflect"
)

type cmdTag struct {
	inject.BaseTag
}

func init() {
	inject.AddTag(new(cmdTag))
}

func (t *cmdTag) IsSingleton() bool {
	return true
}

func (t *cmdTag) Decode(object reflect.Value, field reflect.StructField, tag string) (retVal interface{}) {
	if field.Type.Kind() == reflect.Ptr {
		// parse tag and instantiate filed
		ft := field.Type
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		o := reflect.New(ft)
		retVal = o.Interface()

		cmd := object.Interface().(Command)
		child := retVal.(Command)

		// add child command
		cmd.Add(child)
	}

	return
}
