package types

import "reflect"

type ReflectObject struct {
	Interface interface{}
	Type reflect.Type
	Value reflect.Value
}
