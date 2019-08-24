package annotation

import (
	"fmt"
	"github.com/fatih/structtag"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/utils/reflector"
	"hidevops.io/hiboot/pkg/utils/str"
	"reflect"
)

// GetField is a function that get specific annotation field of an object.
func GetField(object interface{}, att interface{}) (field reflect.StructField, ok bool) {
	field, ok = reflector.GetEmbeddedFieldType(object, att)
	if ok {
		if reflect.TypeOf(att) == reflect.TypeOf(at.Annotation{}) {
			return
		}
		typ := reflector.IndirectType(field.Type)
		_, ok = reflector.GetEmbeddedFieldByType(typ, at.Annotation{}, reflect.Struct)
	}
	return
}

// GetFields iterate annotations of a struct
func GetFields(object interface{}) []reflect.StructField {
	var fields []reflect.StructField
	reflectType := reflect.TypeOf(object)
	if reflectType = reflector.IndirectType(reflectType); reflectType.Kind() == reflect.Struct {
		for i := 0; i < reflectType.NumField(); i++ {
			v := reflectType.Field(i)
			if v.Anonymous {
				_, ok := reflector.GetEmbeddedFieldByType(v.Type, at.Annotation{}, reflect.Struct)
				if ok {
					fields = append(fields, v)
				}
			}
		}
	}
	return fields
}

// Has is a function that check if object is the implements of specific Annotation
func Contains(object interface{}, at interface{}) (ok bool) {
	_, ok = GetField(object, at)
	return
}

// InjectIntoObject inject annotations into object
func InjectIntoObject(object interface{}) (err error) {
	// convert to ptr if it is struct object
	ot := reflect.TypeOf(object)
	if ot == nil {
		err = fmt.Errorf("object must not be nil")
		return
	}
	if ot.Kind() != reflect.Ptr {
		err = fmt.Errorf("object must be the point of a struct")
		log.Error(err)
		return
	}

	annotationFields := GetFields(object)
	for _, annotationField := range annotationFields {
		tags, err := structtag.Parse(string(annotationField.Tag))
		if err != nil {
			log.Error(err)
			continue
		}
		// iterate over all tags
		objectValue := reflector.Indirect(reflect.ValueOf(object))
		fieldValue := objectValue.FieldByName(annotationField.Name)
		for _, tag := range tags.Tags() {
			tagObjectValue := fieldValue.FieldByName(str.ToCamel(tag.Key))
			v := str.Convert(tag.Name, tagObjectValue.Kind())
			if tagObjectValue.CanSet() {
				tagObjectValue.Set(reflect.ValueOf(v))
			}
		}
	}
	return
}