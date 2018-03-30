package utils

import (
	"reflect"
	"fmt"
	"github.com/hidevopsio/hi/boot/pkg/log"
)

func Merge(to interface{}, from interface{}) {
	f := reflect.ValueOf(from).Elem()
	t := reflect.ValueOf(to).Elem()

	for i := 0; i < f.NumField(); i++ {
		varName := f.Type().Field(i).Name
		//varType := f.Type().Field(i).Type
		ff := f.Field(i)
		varValue := ff.Interface()
		//log.Debugf("%v %v %v\n", varName, varType, varValue)
		tf := t.FieldByName(varName)

		if tf.IsValid() && tf.CanSet() {
			kind := tf.Kind()
			switch kind {
			case reflect.String:
				fv := fmt.Sprintf("%v", varValue)
				if fv != "" {
					tf.SetString(fmt.Sprintf("%v", varValue))
				}
				break
			case reflect.Slice:
				if ff.Len() != 0 {
					// TODO merge slice recursively?
					log.Println(varName, varValue)
					tf.Set(reflect.ValueOf(varValue))
				}
				break
			case reflect.Map:
				if ff.Len() != 0 {
					// TODO merge map recursively?
					log.Println(varName, varValue)
					tf.Set(reflect.ValueOf(varValue))
				}
				break
			default:
				break
			}

		}
	}
}
