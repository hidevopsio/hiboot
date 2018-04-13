package utils

import (
	"reflect"
)


func NewReflectType(st interface{}) interface{} {
	ct := reflect.TypeOf(st)
	co := reflect.New(ct)
	cp := co.Elem().Addr().Interface()
	return cp
}
