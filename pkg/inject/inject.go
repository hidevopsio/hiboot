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
	table = "table"
)

var (
	autoConfiguration starter.AutoConfiguration
)


func init() {
	log.SetLevel(log.DebugLevel)
	autoConfiguration = starter.GetAutoConfiguration()
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

	configurations := autoConfiguration.Configurations()
	instances := autoConfiguration.Instances()

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
		if instanceName != "" && fieldObj.CanSet() {
			// first, find if object is already instantiated
			fo := instances[instanceName]
			if fo == nil {
				// parse tag and instantiate filed
				dst := tags[dataSourceType]
				switch {
				case dst != "":
					cfs := configurations[dst]
					if cfs != nil {
						name := tags[namespace]
						if name == "" {
							name = tags[table]
						}
						if name == "" {
							log.Warn("please specify namespace or table for the repository")
							break
						}
						fo = newRepository(cfs, tags[namespace])
						instances[instanceName] = fo
					}

				default:
					o := reflect.New(ft)
					fo = o.Interface()
					instances[instanceName] = fo
				}
			}
			// set field object
			if fo != nil {
				fov := reflect.ValueOf(fo)
				fieldObj.Set(fov)
				log.Debugf("Injected %v(%v) into %v.%v", fov, fov.Type(), obj.Type(), f.Name)
			}
		}
		//log.Debugf("- kind: %v, %v, %v", obj.Kind(), object.Type(), fieldObj.Type())
		//log.Debugf("isValid: %v, canSet: %v", fieldObj.IsValid(), fieldObj.CanSet())
		if obj.Kind() == reflect.Struct && fieldObj.IsValid() && fieldObj.CanSet() {
			IntoObject(fieldObj)
		}
	}
}


