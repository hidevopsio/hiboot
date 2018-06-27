package web

import (
	"reflect"
	"github.com/hidevopsio/hiboot/pkg/utils"
	"github.com/hidevopsio/hiboot/pkg/log"
)

type Inject struct {

}

func (i *Inject) IntoObject(instance reflect.Value) error  {
	for _, f := range utils.DeepFields(instance.Type()) {
		log.Debugf("name: %v, type: %v, tag: %v", f.Name, f.Type, f.Tag)
		if f.Tag != "" {
			inst := instance
			if instance.Kind() == reflect.Ptr {
				inst = instance.Elem()
			}
			log.Debugf("+ %v, %v", instance.Kind(), inst.Kind())
			if inst.Kind() == reflect.Struct {
				if child := inst.FieldByName(f.Name); child.IsValid() && child.CanSet() {
					i.IntoObject(child)
				}
			}
			log.Debugf("- %v, %v", instance.Kind(), inst.Kind())

			// parse tag and instantiate filed
		}
	}

	return nil
}

