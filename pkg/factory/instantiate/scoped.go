package instantiate

import (
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/inject/annotation"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"reflect"
	"strings"
)

// ScopedInstanceFactory implements ScopedInstanceFactory
type ScopedInstanceFactory[T any] struct {
}

var (
	instFactory        factory.InstantiateFactory
	instanceContainers factory.InstanceContainer
)

func initScopedFactory(fct factory.InstantiateFactory) {
	instFactory = fct
	instanceContainers = newInstanceContainer(cmap.New())
}

func (f *ScopedInstanceFactory[T]) GetInstance(params ...interface{}) (retVal T) {
	var err error
	var t T
	retVal = reflect.Zero(reflect.TypeOf(t)).Interface().(T)

	typ := reflect.TypeOf(t)
	typName := reflector.GetLowerCamelFullNameByType(typ)

	instItf := instFactory.GetInstance(typ, factory.MetaData{})
	if instItf == nil {
		retVal = reflector.New[T]()
		ann := annotation.GetAnnotation(t, at.Scope{})
		if ann == nil {
			// default is singleton
			_ = instFactory.SetInstance(retVal)
		}
		return
	} else {
		// TODO: check if instance is prototype?
		instObj := instItf.(*factory.MetaData)
		if instObj.Instance != nil {
			retVal = instObj.Instance.(T)
			return
		}
	}

	inst := instItf.(*factory.MetaData)
	if inst.Scope == factory.ScopePrototype {
		conditionalKey := typName
		instanceContainer := newInstanceContainer(cmap.New())
		if len(params) > 0 {
			for _, param := range params {
				if param != nil {
					ann := annotation.GetAnnotation(param, at.ConditionalOnField{})
					if ann != nil {
						fieldNames, ok := ann.Field.StructField.Tag.Lookup("value")
						if ok {

							// Split the fieldNames string by comma to get individual field names
							fieldList := strings.Split(fieldNames, ",")

							for _, fieldName := range fieldList {
								fv := reflector.GetFieldValue(param, fieldName)
								switch fv.Interface().(type) {
								case string:
									conditionalKey = conditionalKey + "-" + fv.Interface().(string)
								}
							}
						}
					}
					err = instanceContainer.Set(param)
					log.Debugf("set instance %v error code: %v", param, err)
				}
			}
			ic := instanceContainers.Get(conditionalKey)
			if ic != nil {
				instanceContainer = ic.(factory.InstanceContainer)
				finalInst := instanceContainer.Get(typ)
				retVal = finalInst.(T)
				log.Infof("found prototype scoped instance: %v", retVal)
				return
			} else {
				err = instanceContainers.Set(conditionalKey, instanceContainer)
				log.Debugf("set instance %v error code: %v", conditionalKey, err)
			}
		}

		err = instFactory.InjectScopedDependencies(instanceContainer, []*factory.MetaData{inst})
		if err == nil {
			finalInst := instanceContainer.Get(typ)
			retVal = finalInst.(T)
		}
	}

	return
}
