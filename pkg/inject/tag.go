// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package inject

import (
	"github.com/hidevopsio/hiboot/pkg/system"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"github.com/hidevopsio/hiboot/pkg/utils/replacer"
	"reflect"
	"strings"
)

type Tag interface {
	Init(systemConfig *system.Configuration, configurations cmap.ConcurrentMap)
	Decode(object reflect.Value, field reflect.StructField, tag string) (retVal interface{})
	Properties() cmap.ConcurrentMap
	IsSingleton() bool
}

type BaseTag struct {
	properties     cmap.ConcurrentMap
	systemConfig   *system.Configuration
	configurations cmap.ConcurrentMap
}

func (t *BaseTag) IsSingleton() bool {
	return false
}

func (t *BaseTag) Init(systemConfig *system.Configuration, configurations cmap.ConcurrentMap) {
	t.systemConfig = systemConfig
	t.configurations = configurations
}

// TODO move to replacer ?
func (t *BaseTag) replaceReferences(val string) interface{} {
	var retVal interface{}
	retVal = val

	matches := replacer.GetMatches(val)
	if len(matches) != 0 {
		for _, m := range matches {
			//log.Debug(m[1])
			// default value

			vars := strings.SplitN(m[1], ".", -1)
			configName := vars[0]
			// trying to find config
			config, ok := t.configurations.Get(configName)
			sysConf, err := replacer.GetReferenceValue(t.systemConfig, configName)
			if !ok && err == nil && sysConf.IsValid() {
				config = t.systemConfig
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
			val := v[n+1:]
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

func (t *BaseTag) Decode(object reflect.Value, field reflect.StructField, tag string) (retVal interface{}) {
	return nil
}
