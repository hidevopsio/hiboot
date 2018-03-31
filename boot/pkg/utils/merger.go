package utils

import (
	"reflect"
	"fmt"
	"github.com/hidevopsio/hi/boot/pkg/log"
	"strings"
	"errors"
	"database/sql"
	"strconv"
)

func validate(toValue interface{}) (*reflect.Value, error)  {

	to := indirect(reflect.ValueOf(toValue))

	if !to.CanAddr() {
		return nil, errors.New("value is unaddressable")
	}

	// Return is from value is invalid
	if !to.IsValid() {
		return  nil, errors.New("value is not valid")
	}

	return &to, nil
}


// Copy copy things
func Copy(toValue interface{}, fromValue interface{}) (err error) {
	var (
		isSlice bool
		amount  = 1
		from    = indirect(reflect.ValueOf(fromValue))
		to      = indirect(reflect.ValueOf(toValue))
	)

	if !to.CanAddr() {
		return errors.New("copy to value is unaddressable")
	}

	// Return is from value is invalid
	if !from.IsValid() {
		return
	}

	// Just set it if possible to assign
	if from.Type().AssignableTo(to.Type()) {
		to.Set(from)
		return
	}

	fromType := indirectType(from.Type())
	toType := indirectType(to.Type())

	if fromType.Kind() != reflect.Struct || toType.Kind() != reflect.Struct {
		return
	}

	if to.Kind() == reflect.Slice {
		isSlice = true
		if from.Kind() == reflect.Slice {
			amount = from.Len()
		}
	}

	for i := 0; i < amount; i++ {
		var dest, source reflect.Value

		if isSlice {
			// source
			if from.Kind() == reflect.Slice {
				source = indirect(from.Index(i))
			} else {
				source = indirect(from)
			}

			// dest
			dest = indirect(reflect.New(toType).Elem())
		} else {
			source = indirect(from)
			dest = indirect(to)
		}

		// Copy from field to field or method
		for _, field := range deepFields(fromType) {
			name := field.Name

			if fromField := source.FieldByName(name); fromField.IsValid() {
				// has field
				if toField := dest.FieldByName(name); toField.IsValid() {
					if toField.CanSet() {
						if !set(toField, fromField) {
							if err := Copy(toField.Addr().Interface(), fromField.Interface()); err != nil {
								return err
							}
						}
					}
				} else {
					// try to set to method
					var toMethod reflect.Value
					if dest.CanAddr() {
						toMethod = dest.Addr().MethodByName(name)
					} else {
						toMethod = dest.MethodByName(name)
					}

					if toMethod.IsValid() && toMethod.Type().NumIn() == 1 && fromField.Type().AssignableTo(toMethod.Type().In(0)) {
						toMethod.Call([]reflect.Value{fromField})
					}
				}
			}
		}

		// Copy from method to field
		for _, field := range deepFields(toType) {
			name := field.Name

			var fromMethod reflect.Value
			if source.CanAddr() {
				fromMethod = source.Addr().MethodByName(name)
			} else {
				fromMethod = source.MethodByName(name)
			}

			if fromMethod.IsValid() && fromMethod.Type().NumIn() == 0 && fromMethod.Type().NumOut() == 1 {
				if toField := dest.FieldByName(name); toField.IsValid() && toField.CanSet() {
					values := fromMethod.Call([]reflect.Value{})
					if len(values) >= 1 {
						set(toField, values[0])
					}
				}
			}
		}

		if isSlice {
			if dest.Addr().Type().AssignableTo(to.Type().Elem()) {
				to.Set(reflect.Append(to, dest.Addr()))
			} else if dest.Type().AssignableTo(to.Type().Elem()) {
				to.Set(reflect.Append(to, dest))
			}
		}
	}
	return
}

func deepFields(reflectType reflect.Type) []reflect.StructField {
	var fields []reflect.StructField

	if reflectType = indirectType(reflectType); reflectType.Kind() == reflect.Struct {
		for i := 0; i < reflectType.NumField(); i++ {
			v := reflectType.Field(i)
			if v.Anonymous {
				fields = append(fields, deepFields(v.Type)...)
			} else {
				fields = append(fields, v)
			}
		}
	}

	return fields
}

func indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

func indirectType(reflectType reflect.Type) reflect.Type {
	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
	}
	return reflectType
}

func set(to, from reflect.Value) bool {
	if from.IsValid() {
		if to.Kind() == reflect.Ptr {
			//set `to` to nil if from is nil
			if from.Kind() == reflect.Ptr && from.IsNil() {
				to.Set(reflect.Zero(to.Type()))
				return true
			} else if to.IsNil() {
				to.Set(reflect.New(to.Type().Elem()))
			}
			to = to.Elem()
		}

		if from.Type().ConvertibleTo(to.Type()) {
			to.Set(from.Convert(to.Type()))
		} else if scanner, ok := to.Addr().Interface().(sql.Scanner); ok {
			err := scanner.Scan(from.Interface())
			if err != nil {
				return false
			}
		} else if from.Kind() == reflect.Ptr {
			return set(to, from.Elem())
		} else {
			return false
		}
	}
	return true
}

func Merge(to interface{}, from interface{}) {

	f := reflect.ValueOf(from).Elem()
	t := reflect.ValueOf(to).Elem()

	for i := 0; i < f.NumField(); i++ {
		varName := f.Type().Field(i).Name
		//varType := f.Type().Field(i).Type
		ff := f.Field(i)
		varValue := ff.Interface()
		//log.Debugf("%v %v %v\n", varName, varType, varValue)
		tf := t.FieldByName(varName)

		if tf.IsValid() && tf.CanSet() {
			kind := tf.Kind()
			switch kind {
			case reflect.String:
				fv := fmt.Sprintf("%v", varValue)
				if fv != "" {
					tf.SetString(fmt.Sprintf("%v", varValue))
				}

			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				fv := fmt.Sprintf("%v", varValue)
				if fv != "" {
					fv, _ := strconv.Atoi(fv)
					tf.SetInt(int64(fv))
				}
			default:
				Merge(t.Addr().Interface(), f.Addr().Interface())
				break
			}

		}
	}
}


func Replace(to interface{}, name, value string) {

	//t := reflect.ValueOf(to).Elem()
	t, err := validate(to)
	if err != nil {
		return
	}

	toType := indirectType(t.Type())
	//k := toType.Kind()
	//if k != reflect.Struct {
	//	return
	//}

	isSlice := false
	fieldSize := 1
	if t.Kind() == reflect.Slice {
		isSlice = true
		fieldSize = t.Len()
	}

	dstVarName := "${" + name + "}"

	for i := 0; i < fieldSize; i++ {
		var dst reflect.Value
		if isSlice {
			//dst = indirect(reflect.New(toType).Elem())
			if t.Kind() == reflect.Slice {
				dst = indirect(t.Index(i))
				log.Debug(dst.Interface())

				// TODO: refactoring below code
				dstType := dst.Type().Name()
				dstValue := dst.Interface()
				log.Debug(dstType)
				dv := fmt.Sprintf("%v", dstValue)
				if strings.Contains(dv, dstVarName) {
					if dstType == "string" && dst.IsValid() && dst.CanSet() {
						dv = strings.Replace(dv, dstVarName, value, -1)
						dst.SetString(dv)
					}
				}
			} else {
				dst = indirect(*t)
			}
		} else {
			dst = indirect(*t)
		}


		for _, field := range deepFields(toType) {
			fieldName := field.Name
			log.Debug("fieldName: ", fieldName)
			if dstField := dst.FieldByName(fieldName); dstField.IsValid() && dstField.CanSet() {
				fieldValue := dstField.Interface()
				log.Debug("fieldValue: ", fieldValue)

				kind := dstField.Kind()
				switch kind {
				case reflect.String:
					fv := fmt.Sprintf("%v", fieldValue)
					log.Debug(fieldName, ": ", fieldValue)
					if strings.Contains(fv, dstVarName) {
						newVal := strings.Replace(fv, dstVarName, value, -1)
						log.Debug("newVal: " + newVal)
						dstField.SetString(newVal)
					}
				default:
					log.Debug(fieldName, " is a ", kind)
					Replace(dstField.Addr().Interface(), name, value)
				}
			}
		}
	}
}
