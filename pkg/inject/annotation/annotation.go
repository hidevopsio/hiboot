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
	"hidevops.io/hiboot/pkg/types"
	"hidevops.io/hiboot/pkg/utils/mapstruct"
	"hidevops.io/hiboot/pkg/utils/reflector"
	"hidevops.io/hiboot/pkg/utils/str"
	"hidevops.io/hiboot/pkg/utils/structtag"
	"reflect"
	"strings"
)

// annotation field
type Field struct{
	StructField reflect.StructField
	Value reflect.Value
}

type Annotation struct{
	Field *Field
	Parent *types.ReflectObject
}

type Annotations struct{
	Items    []*Annotation
	Children []*Annotations
}

// GetAnnotation is a function that get specific annotation of an object.
func GetAnnotation(object interface{}, att interface{}) (annotation *Annotation) {
	if object == nil {
		return
	}
	atType := reflect.TypeOf(att)
	switch object.(type) {
	case *Annotations:
		ans := object.(*Annotations)
		if ans == nil {
			return
		}
		for _, item := range ans.Items {
			if item.Field.StructField.Type == atType {
				annotation = item
				return
			}
		}
	}

	structField, ok := reflector.GetEmbeddedFieldType(object, att)
	if ok {
		fieldType := reflector.IndirectType(structField.Type)
		if atType != reflect.TypeOf(at.Annotation{}) {
			_, ok = reflector.GetEmbeddedFieldByType(fieldType, at.Annotation{}, reflect.Struct)
		}
		if ok {
			ov := reflector.IndirectValue(object)

			reflectObject := &types.ReflectObject{
				Interface: object,
				Type:      fieldType,
				Value:     ov,
			}
			annotation = new(Annotation)
			annotation.Field = new(Field)
			if ov.IsValid() && ov.CanAddr() && ov.Type().Kind() == reflect.Struct {
				annotation.Field.Value = ov.FieldByName(fieldType.Name())
			}
			annotation.Field.StructField = structField
			annotation.Parent = reflectObject
		}
	}
	return
}

// GetFields iterate annotations of a struct
func GetAnnotations(object interface{}) (annotations *Annotations) {
	if object == nil {
		return
	}
	reflectType, ok := reflector.GetObjectType(object)
	if ok {
		ov := reflector.IndirectValue(object)
		reflectObject := &types.ReflectObject{
			Interface: object,
			Type:      reflectType,
			Value:     ov,
		}
		if reflectType.Kind() == reflect.Struct {
			annotations = new(Annotations)
			for i := 0; i < reflectType.NumField(); i++ {
				ann := new(Annotation)
				ann.Field = new(Field)
				f := reflectType.Field(i)
				typ := f.Type
				iTyp := reflector.IndirectType(typ)
				//log.Debugf("%v %v", f.Name, iTyp.Name() )
				if f.Anonymous {
					_, ok = reflector.GetEmbeddedFieldByType(typ, at.Annotation{}, reflect.Struct)
					if ok {
						//log.Debugf("%v %v %v", ov.IsValid(), ov.CanAddr(), ov.Type().Kind())
						if ov.IsValid() && ov.CanAddr() && ov.Type().Kind() == reflect.Struct{
							ann.Field.Value = ov.FieldByName(f.Name)
						}
						ann.Field.StructField = f
						ann.Parent = reflectObject
						annotations.Items = append(annotations.Items, ann)
					}
				} else {
					if iTyp.Name() == "" && typ.Kind() == reflect.Struct {
						// more annotations from child struct
						fieldObjVal := ov.FieldByName(f.Name)
						if f.Name == strings.Title(f.Name) {
							childAnnotations := GetAnnotations(fieldObjVal.Addr().Interface())
							annotations.Children = append(annotations.Children, childAnnotations)
						}
					}
				}
			}
		}
	}
	return
}

// Filter is a function that filter specific annotations.
func FilterIn(input *Annotations, att interface{}) (annotations []*Annotation) {
	var ok bool
	if input != nil {
		for _, item := range input.Items {
			if item.Field.Value.IsValid() {
				ok = item.Field.StructField.Type == reflect.TypeOf(att)
				ok =  ok || reflector.HasEmbeddedFieldType(item.Field.Value.Interface(), att)
				if ok {
					annotations = append(annotations, item)
				}
			}
		}

		for _, child := range input.Children {
			childAnnotations := FilterIn(child, att)
			if childAnnotations != nil {
				annotations = append(annotations, childAnnotations...)
			}
		}
	}

	return
}

// ContainsChild is a function that find specific annotations.
func ContainsChild(input *Annotations, att interface{}) (ok bool) {
	items := FilterIn(input, att)
	ok = len(items) > 0
	return
}

// Find is a function that find specific (child) annotation
func Find(input *Annotations, att interface{}) (annotation *Annotation) {
	items := FilterIn(input, att)
	if len(items) > 0 {
		annotation = items[0]
	}
	return
}

// Find is a function that find specific annotation.
func FindAll(object interface{}, att interface{}) (annotations []*Annotation) {
	ans := GetAnnotations(object)
	annotations = FilterIn(ans, att)
	return
}

// Contains Has is a function that check if object is the implements of specific Annotation
func Contains(object interface{}, at interface{}) (ok bool) {
	ok = GetAnnotation(object, at) != nil
	return
}

// Inject inject annotations into object
func Inject(ann *Annotation) (err error) {
	var tags *structtag.Tags
	if ann.Field.Value.IsValid() {
		err = injectIntoField(tags, ann.Field)
	}
	return
}

// InjectItems inject annotations into object
func InjectItems(annotations *Annotations) (err error) {
	var tags *structtag.Tags
	for _, item := range annotations.Items {
		if item.Field.Value.IsValid() {
			err = injectIntoField(tags, item.Field)
			if err != nil {
				break
			}
		}
	}

	for _, child := range annotations.Children {
		err = InjectItems(child)
	}
	return
}

// InjectAll inject annotations into object
func InjectAll(object interface{}) (err error) {
	// convert to ptr if it is struct object
	ot := reflect.TypeOf(object)
	if ot == nil {
		err = fmt.Errorf("object must not be nil")
		return
	}
	if ot.Kind() != reflect.Ptr {
		err = fmt.Errorf("object %v is not the point of a struct", ot.Name())
		return
	}

	annotations := GetAnnotations(object)

	err = InjectItems(annotations)
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