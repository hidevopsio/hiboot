package inject

import (
	"strings"
	"github.com/hidevopsio/hiboot/pkg/starter"
	"github.com/hidevopsio/hiboot/pkg/utils"
	"reflect"
)

type Tag interface {
	Decode(object reflect.Value, field reflect.StructField, tag string) (retVal interface{})
	Properties() map[string]interface{}
	IsSingleton() bool
}

type BaseTag struct {
	properties map[string]interface{}
}

func (t *BaseTag) IsSingleton() bool  {
	return true
}

func (t *BaseTag) replaceReferences(val string) interface{}  {
	var retVal interface{}
	retVal = val
	systemConfig := autoConfiguration.Configuration(starter.System)

	matches := utils.GetMatches(val)
	if len(matches) != 0 {
		for _, m := range matches {
			//log.Debug(m[1])
			// default value

			vars := strings.SplitN(m[1], ".", -1)
			configName := vars[0]
			config := autoConfiguration.Configuration(configName)
			sysConf, err := utils.GetReferenceValue(systemConfig, configName)
			if config == nil && err == nil && sysConf.IsValid() {
				config = systemConfig
			}
			if config != nil {
				retVal = utils.ReplaceStringVariables(val, config)
				if retVal != val {
					break
				}
			}
		}
	}
	return retVal
}

func (t *BaseTag) ParseProperties(tag string) map[string]interface{} {
	t.properties = make(map[string]interface{}) // ? map[string]string

	args := strings.Split(tag, ",")
	for _, v := range args {
		//log.Debug(v)
		n := strings.Index(v, "=")
		if n > 0 {
			key := v[:n]
			val := v[n + 1:]
			if key != "" && val != "" {
				// check if val contains reference or env
				// TODO: should lookup certain config instead of for loop
				replacedVal := t.replaceReferences(val)
				t.properties[key] = replacedVal
			}
		}
	}
	return t.properties
}

func (t *BaseTag) Properties() map[string]interface{} {
	return t.properties
}

func (t *BaseTag) Decode(object reflect.Value, field reflect.StructField, tag string) (retVal interface{})  {
	return nil
}