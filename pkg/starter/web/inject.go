package web

import (
	"reflect"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
)

type Inject struct {

}

const (
	injectIdentifier = "component"

	injectTypeRepository = "repository"
	injectTypeService = "service"

	dataSourceType = "dataSourceType"
)

func (i *Inject) IntoObject(object reflect.Value, dataSources DataSources) error  {
	for _, f := range reflector.DeepFields(object.Type()) {
		//log.Debugf("parent: %v, name: %v, type: %v, tag: %v", object.Elem().Type(), f.Name, f.Type, f.Tag)
		component := f.Tag.Get(injectIdentifier)
		if component != "" {
			obj := reflector.Indirect(object)

			var fieldObj reflect.Value
			if obj.IsValid() {
				fieldObj = obj.FieldByName(f.Name)
			}
			ft := f.Type
			if f.Type.Kind() == reflect.Ptr {
				ft = ft.Elem()
			}
			//log.Debugf("+ %v, %v, %v", object.Type(), fieldObj.Type(), ft)

			if fieldObj.CanSet() {
				// parse tag and instantiate filed
				switch component {
				case injectTypeService:

					o := reflect.New(ft)
					//log.Debug(fieldObj, o)
					fieldObj.Set(o)
					//log.Debug(fieldObj, o)
					log.Debugf("Injected service %v into %v.%v", o, obj.Type(), f.Name)
				case injectTypeRepository:
					dataSourceType := f.Tag.Get(dataSourceType)
					o := dataSources[dataSourceType]
					if o != nil {
						ov := reflect.ValueOf(o)
						fieldObj.Set(ov)
						log.Debugf("Injected repository %v into %v.%v", ov, obj.Type(), f.Name)
					}
				}
			}

			if obj.Kind() == reflect.Struct && fieldObj.IsValid() && fieldObj.CanSet() {
				i.IntoObject(fieldObj, dataSources)
			}
			//log.Debugf("- %v, %v", object.Type(), fieldObj.Type())
		}
	}

	return nil
}


