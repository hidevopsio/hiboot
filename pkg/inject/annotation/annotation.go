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
// package annotation provides annotation support for HiBoot
package annotation

import (
	"fmt"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/utils/mapstruct"
	"hidevops.io/hiboot/pkg/utils/reflector"
	"hidevops.io/hiboot/pkg/utils/str"
	"hidevops.io/hiboot/pkg/utils/structtag"
	"reflect"
)

// annotation field
type Field struct{
	StructField reflect.StructField
	Value reflect.Value
}

// GetField is a function that get specific annotation field of an object.
func GetField(object interface{}, att interface{}) (field *Field, ok bool) {
	field = new(Field)
	var structField reflect.StructField
	atType := reflect.TypeOf(att)

	switch object.(type) {
	case []*Field:
		fields := object.([]*Field)
		for _, f := range fields {
			if f.StructField.Type == atType {
				field = f
				ok = true
				return
			}
		}
	}

	structField, ok = reflector.GetEmbeddedFieldType(object, att)
	if ok {
		fieldType := reflector.IndirectType(structField.Type)
		if atType != reflect.TypeOf(at.Annotation{}) {
			_, ok = reflector.GetEmbeddedFieldByType(fieldType, at.Annotation{}, reflect.Struct)
		}
		if ok {
			ov := reflector.IndirectValue(object)
			if ov.IsValid() && ov.CanAddr() && ov.Type().Kind() == reflect.Struct {
				field.Value = ov.FieldByName(fieldType.Name())
			}
			field.StructField = structField
		}
	}
	return
}

// GetFields iterate annotations of a struct
// TODO: get field at...
func GetFields(object interface{}) []*Field {
	var fields []*Field

	reflectType, ok := reflector.GetObjectType(object)
	if ok {
		ov := reflector.IndirectValue(object)
		if reflectType.Kind() == reflect.Struct {
			for i := 0; i < reflectType.NumField(); i++ {
				field := new(Field)
				f := reflectType.Field(i)
				typ := f.Type
				if f.Anonymous {
					_, ok = reflector.GetEmbeddedFieldByType(typ, at.Annotation{}, reflect.Struct)
					if ok {
						if ov.IsValid() && ov.CanAddr() && ov.Type().Kind() == reflect.Struct{
							field.Value = ov.FieldByName(f.Name)
						}
						field.StructField = f
						fields = append(fields, field)
					}
				} else {
					iTyp := reflector.IndirectType(typ)
					if iTyp.Name() == "" && typ.Kind() == reflect.Struct {
						// more annotations from child struct
						fieldObjVal := ov.FieldByName(f.Name)
						childFields := GetFields(fieldObjVal)
						fields = append(fields, childFields...)
					}
				}
			}
		}
	}
	return fields
}

// Filter is a function that filter specific annotations.
func Filter(input []*Field, att interface{}) (fields []*Field) {
	for _, f := range input {
		if f.Value.IsValid() {
			ok := f.StructField.Type == reflect.TypeOf(att)
			ok =  ok || reflector.HasEmbeddedFieldType(f.Value.Interface(), att)
			if ok {
				fields = append(fields, f)
			}
		}
	}
	return
}


// ContainsChild is a function that find specific annotations.
func ContainsChild(input []*Field, att interface{}) (ok bool) {
	f := Filter(input, att)
	ok = len(f) > 0
	return
}

// Find is a function that find specific annotation.
func Find(object interface{}, att interface{}) (fields []*Field) {
	allFields := GetFields(object)
	fields = Filter(allFields, att)
	return
}

// Has is a function that check if object is the implements of specific Annotation
func Contains(object interface{}, at interface{}) (ok bool) {
	_, ok = GetField(object, at)
	return
}

// InjectIntoField inject annotations into object
func InjectIntoField(field *Field) (err error) {
	var tags *structtag.Tags
	if field.Value.IsValid() {
		err = injectIntoField(tags, field)
	}
	return
}


// InjectIntoFields inject annotations into object
func InjectIntoFields(object interface{}) (err error) {
	// convert to ptr if it is struct object
	ot := reflect.TypeOf(object)
	if ot == nil {
		err = fmt.Errorf("object must not be nil")
		return
	}
	if ot.Kind() != reflect.Ptr {
		err = fmt.Errorf("object %v is not the point of a struct", ot.Name())
		//log.Error(err)
		return
	}

	fields := GetFields(object)
	var tags *structtag.Tags
	for _, field := range fields {
		if field.Value.IsValid() {
			err = injectIntoField(tags, field)
			if err != nil {
				break
			}
		}
	}
	return
}

func injectIntoField(tags *structtag.Tags, field *Field) (err error) {
	tags, err = structtag.Parse(string(field.StructField.Tag))
	if err != nil {
		log.Errorf("%v of %v", err, field.StructField.Type)
		return
	}
    fieldValue := field.Value
	typeField, ok := fieldValue.Type().FieldByName("FieldName")
	if ok {
		valueFieldName := typeField.Tag.Get("value")
		if valueFieldName != "" {
			valueFieldValue := field.Value.FieldByName(str.ToCamel(valueFieldName))
			if valueFieldValue.CanSet() {
				switch valueFieldValue.Interface().(type) {
				case map[int]string:
					values := make( map[int]string)
					for _, tag := range tags.Tags() {
						k := str.Convert(tag.Key, reflect.Int).(int)
						values[k] = tag.Name
					}
					valueFieldValue.Set(reflect.ValueOf(values))
					return
				case map[string]string:
					values := make( map[string]string)
					for _, tag := range tags.Tags() {
						values[tag.Key] = tag.Name
					}
					valueFieldValue.Set(reflect.ValueOf(values))
					return
				}
			}
		}
	}

	// iterate over all tags
	if tags != nil {
		values := make( map[string]string)
		for _, tag := range tags.Tags() {
			values[tag.Key] = tag.Name
		}
		if len(values) != 0 {
			// use mapstruct.WithSquash to decode embedded sub field
			err = mapstruct.Decode(fieldValue.Addr().Interface(), values, mapstruct.WithSquash)
		}
	}
	return
}