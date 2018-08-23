package inject

import (
	"strings"
	"reflect"
	"github.com/hidevopsio/hiboot/pkg/starter"
	"github.com/hidevopsio/hiboot/pkg/utils/replacer"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
)

type Tag interface {
	Decode(object reflect.Value, field reflect.StructField, tag string) (retVal interface{})
	Properties() cmap.ConcurrentMap
	IsSingleton() bool
}

type BaseTag struct {
	properties cmap.ConcurrentMap
}

func (t *BaseTag) IsSingleton() bool  {
	return false
}

func (t *BaseTag) replaceReferences(val string) interface{}  {
	var retVal interface{}
	retVal = val
	systemConfig := autoConfiguration.Configuration(starter.System)

	matches := replacer.GetMatches(val)
	if len(matches) != 0 {
		for _, m := range matches {
			//log.Debug(m[1])
			// default value

			vars := strings.SplitN(m[1], ".", -1)
			configName := vars[0]
			config := autoConfiguration.Configuration(configName)
			sysConf, err := replacer.GetReferenceValue(systemConfig, configName)
			if config == nil && err == nil && sysConf.IsValid() {
				config = systemConfig
			}
			if config != nil {
				retVal = replacer.ReplaceStringVariables(val, config)
				if retVal != val {
					break
				}
			}
		}
	}
	return retVal
}

func (t *BaseTag) ParseProperties(tag string) cmap.ConcurrentMap {
	t.properties = cmap.New()

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
				t.properties.Set(key, replacedVal)
			}
		}
	}
	return t.properties
}

func (t *BaseTag) Properties() cmap.ConcurrentMap {
	return t.properties
}

func (t *BaseTag) Decode(object reflect.Value, field reflect.StructField, tag string) (retVal interface{})  {
	return nil
}