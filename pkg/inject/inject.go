package inject

import (
	"reflect"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"strings"
	"github.com/hidevopsio/hiboot/pkg/starter"
)


const (
	injectIdentifier = "inject"
	dataSourceType = "dataSourceType"
	namespace = "namespace"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func newRepository(c interface{}, repoName string) interface{}  {
	cv := reflect.ValueOf(c)

	configType := cv.Type()
	log.Debug("type: ", configType)
	name := configType.Elem().Name()
	log.Debug("fieldName: ", name)

	// call Init
	numOfMethod := cv.NumMethod()
	log.Debug("methods: ", numOfMethod)

	method := cv.MethodByName("NewRepository")

	numIn := method.Type().NumIn()
	if numIn == 1 {
		argv := make([]reflect.Value, numIn)
		argv[0] = reflect.ValueOf(repoName)
		retVal := method.Call(argv)
		instance := retVal[0].Interface()
		log.Debugf("instantiated: %v", instance)
		return instance
	}

	return nil
}


// IntoObject injects instance into the tagged field with `inject:"instanceName"`
func IntoObject(object reflect.Value) {
	ac := starter.GetInstance()
	configurations := ac.Configurations()
	instances := ac.Instances()

	for _, f := range reflector.DeepFields(object.Type()) {
		//log.Debugf("parent: %v, name: %v, type: %v, tag: %v", object.Elem().Type(), f.Name, f.Type, f.Tag)
		injectTag := f.Tag.Get(injectIdentifier)
		args := strings.Split(injectTag, ",")
		tags := make(map[string]string) // ? map[string]string
		for _, v := range args[1:] {
			//log.Debug(v)
			kv := strings.Split(v, "=")
			if len(kv) == 2 {
				tags[kv[0]] = kv[1]
			}
		}
		instanceName := args[0]

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
		if instanceName != "" {
			if fieldObj.CanSet() {
				// parse tag and instantiate filed
				dst := tags[dataSourceType]
				switch  {
				case dst != "":
					cfs := configurations[dst]
					if cfs != nil {
						repo := newRepository(cfs, tags[namespace])
						val := reflect.ValueOf(repo)
						fieldObj.Set(val)
						log.Debugf("Injected repository %v into %v.%v", val, obj.Type(), f.Name)
					}

				default:
					o := reflect.New(ft)
					instances[instanceName] = o.Interface()
					//log.Debug(fieldObj, o)
					fieldObj.Set(o)
					//log.Debug(fieldObj, o)
					log.Debugf("Injected service %v into %v.%v", o, obj.Type(), f.Name)
				}
			}
		}
		//log.Debugf("- kind: %v, %v, %v", obj.Kind(), object.Type(), fieldObj.Type())
		//log.Debugf("isValid: %v, canSet: %v", fieldObj.IsValid(), fieldObj.CanSet())
		if obj.Kind() == reflect.Struct && fieldObj.IsValid() && fieldObj.CanSet() {
			IntoObject(fieldObj)
		}
	}
}


