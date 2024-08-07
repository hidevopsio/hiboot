package instantiate

import (
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
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

// GetInstanceContainer get the instance container
func (f *ScopedInstanceFactory[T]) GetInstanceContainer(params ...interface{}) (instanceContainer factory.InstanceContainer, err error) {
	var t T
	instanceContainer = newInstanceContainer(cmap.New())
	typ := reflect.TypeOf(t)
	typName := reflector.GetLowerCamelFullNameByType(typ)

	finalInstance := reflect.Zero(reflect.TypeOf(t)).Interface().(T)
	instItf := instFactory.GetInstance(typ, factory.MetaData{})
	if instItf == nil {
		return
	}
	var ctx context.Context
	inst := instItf.(*factory.MetaData)
	if inst.Scope == factory.ScopePrototype {
		conditionalKey := typName
		if len(params) > 0 {
			for _, param := range params {
				if param != nil {
					switch param.(type) {
					case context.Context:
						ctx = param.(context.Context)
						continue
					default:
						conditionalKey = f.parseConditionalField(param, conditionalKey)
					}
					err = instanceContainer.Set(param)
					log.Debugf("set instance %v: %v error code: %v", conditionalKey, param, err)
				}
			}
			// check if instanceContainer already exists
			ic := instanceContainers.Get(conditionalKey)
			if ic != nil {
				// cached instanceContainer
				instanceContainer = ic.(factory.InstanceContainer)
				finalInst := instanceContainer.Get(typ)
				if finalInst != nil {
					finalInstance = finalInst.(T)
					log.Infof("found prototype scoped instance[%v]: %v", conditionalKey, finalInstance)
					return
				}
			} else {
				err = instanceContainers.Set(conditionalKey, instanceContainer)
				log.Debugf("set instance %v error code: %v", conditionalKey, err)
			}
		}

		instanceContainer, err = instFactory.InjectScopedObjects(ctx, []*factory.MetaData{inst}, instanceContainer)
	}
	return
}

// GetInstance get instance container and get the target instance form the container
func (f *ScopedInstanceFactory[T]) GetInstance(params ...interface{}) (finalInstance T, err error) {
	var t T
	finalInstance = reflect.Zero(reflect.TypeOf(t)).Interface().(T)

	typ := reflect.TypeOf(t)
	instItf := instFactory.GetInstance(typ, factory.MetaData{})
	if instItf == nil {
		finalInstance = reflector.New[T]()
		ann := annotation.GetAnnotation(t, at.Scope{})
		if ann == nil {
			// default is singleton
			err = instFactory.SetInstance(finalInstance)
		}
		return
	} else {
		// TODO: check if instance is prototype?
		instObj := instItf.(*factory.MetaData)
		if instObj.Instance != nil {
			finalInstance = instObj.Instance.(T)
			return
		}
	}

	inst := instItf.(*factory.MetaData)
	if inst.Scope == factory.ScopePrototype {
		var instanceContainer factory.InstanceContainer
		instanceContainer, err = f.GetInstanceContainer(params...)
		if err == nil {
			finalInstance = f.GetInstanceFromContainer(instanceContainer)
		}
	}

	return
}

// GetInstanceFromContainer get instance from a giving instance container
func (f *ScopedInstanceFactory[T]) GetInstanceFromContainer(instanceContainer factory.InstanceContainer) (finalInstance T) {
	var t T
	typ := reflect.TypeOf(t)

	finalInst := instanceContainer.Get(typ)
	if finalInst != nil {
		finalInstance = finalInst.(T)
	}
	return
}

func (f *ScopedInstanceFactory[T]) parseConditionalField(param interface{}, conditionalKey string) string {
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
					condition := fv.Interface().(string)
					if condition != "" {
						conditionalKey = conditionalKey + "-" + condition
					}
				}
			}
		}
	}
	return conditionalKey
}
