package cli

import (
	"github.com/hidevopsio/hiboot/pkg/inject"
	"reflect"
	"strconv"
	"strings"
	"github.com/hidevopsio/hiboot/pkg/log"
)

type flagTag struct {
	inject.BaseTag
}

func init() {
	inject.AddTag("flag", new(flagTag))
}

func (t *flagTag) IsSingleton() bool  {
	return false
}

func (t *flagTag) Decode(object reflect.Value, field reflect.StructField, tag string) (retVal interface{}) {
	// Profile string `flag:"shorthand=p,value=dev,usage=--profile=test"`
	if field.Type.Kind() == reflect.Ptr {
		// parse tag and instantiate filed
		ft := field.Type
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		o := reflect.New(ft)
		retVal = o.Interface()

		properties := t.ParseProperties(tag)

		cmd := object.Interface().(Command)

		pflags := cmd.PersistentFlags()

		shorthand := properties["shorthand"].(string)
		value := properties["value"].(string)
		usage := properties["usage"].(string)
		name := strings.ToLower(field.Name)
		switch ft.Kind() {
		case reflect.String:
			fv := retVal.(*string)
			log.Debugf("flag: %v, shorthand: %v, value: %v, usage: %v", name, shorthand, value, usage)
			pflags.StringVarP(fv, name, shorthand, value, usage)
		case reflect.Int:
			fv := retVal.(*int)
			intVal, err := strconv.Atoi(value)
			if err == nil {
				pflags.IntVarP(fv, name, shorthand, intVal, usage)
			}
		}
	}

	return
}