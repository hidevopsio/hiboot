package cli

import (
	"github.com/hidevopsio/hiboot/pkg/inject"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
	"reflect"
)

type cmdTag struct {
	inject.BaseTag
}

func init() {
	inject.AddTag(new(cmdTag))
}

// TODO: hide this method
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
		// first, try to get instance from factory
		instName := str.ToLowerCamel(ft.Name())
		retVal = t.ConfigurableFactory.GetInstance(str.ToLowerCamel(instName))
		// if it does not exist, then create a new instance
		if retVal == nil {
			o := reflect.New(ft)
			retVal = o.Interface()
		}

		cmd := object.Interface().(Command)
		child := retVal.(Command)

		// add child command
		cmd.Add(child)
	}

	return
}
