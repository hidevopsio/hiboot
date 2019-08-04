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
	"fmt"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/utils/reflector"
	"hidevops.io/hiboot/pkg/utils/str"
	"reflect"
)

type valueTag struct {
	at.Tag `value:"value"`
	BaseTag
}

func init() {
	AddTag(new(valueTag))
}

func IsImplemented(object interface{}) {
	s := reflect.ValueOf(object).Elem()
	t := s.Type()
	modelType := reflect.TypeOf((*at.StringAnnotation)(nil)).Elem()

	for i := 0; i < s.NumField(); i++ {
		f := t.Field(i)
		fmt.Printf("%d: %s %s -> %t\n", i, f.Name, f.Type, f.Type.Implements(modelType))
	}
}

func (t *valueTag) Decode(object reflect.Value, field reflect.StructField, property string) (retVal interface{}) {
	tag, ok := field.Tag.Lookup(string(t.Tag))
	if ok {
		//log.Debug(valueTag)

		// check if filed type is slice
		kind := field.Type.Kind()
		needConvert := true
		retVal = t.instantiateFactory.Replace(tag)
		switch kind {
		case reflect.Slice:
			if retVal != tag {
				needConvert = false
			}
		case reflect.String:
			needConvert = false
			// check if it is an annotation
			indTyp := reflector.IndirectType(field.Type)
			//log.Debugf("type name: %v %v %v", field.Type.Name(), indTyp.Name(), indTyp.NumMethod())

			impl := indTyp.Implements(reflect.TypeOf((*at.StringAnnotation)(nil)).Elem())
			if impl {
				//log.Debug("found at.StringAnnotation implements")
				retVal = reflect.New(indTyp).Interface().(at.StringAnnotation).Value(retVal.(string))
			}

		}

		if needConvert {
			in := fmt.Sprintf("%v", retVal)
			retVal = str.Convert(in, kind)
		}
	}
	return retVal
}
