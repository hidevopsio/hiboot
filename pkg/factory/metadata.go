package factory

import (
	"github.com/hidevopsio/hiboot/pkg/system/types"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
	"reflect"
)

// MetaData is the injectable object meta data
type MetaData struct {
	Kind      string
	Name      string
	ShortName string
	TypeName  string
	PkgName   string
	Context   interface{}
	Object    interface{}
	Type      reflect.Type
	ExtDep    []*MetaData
}

// NewMetaData create new meta data
func NewMetaData(params ...interface{}) *MetaData {
	var name string
	var shortName string
	var object interface{}
	var context interface{}

	if len(params) == 2 {
		if reflect.TypeOf(params[0]).Kind() == reflect.String {
			name = params[0].(string)
			object = params[1]
		} else {
			context = params[0]
			object = params[1]
		}
	} else {
		object = params[0]
	}

	pkgName, typeName := reflector.GetPkgAndName(object)
	typ := reflect.TypeOf(object)
	kind := typ.Kind()
	kindName := kind.String()

	if pkgName != "" {
		if name == "" {
			shortName = str.ToLowerCamel(typeName)
			name = pkgName + "." + shortName
		} else {
			shortName = name
			name = pkgName + "." + name
		}
	}
	if kind == reflect.Struct && typ.Name() == types.Method {
		kindName = types.Method
	}
	if kindName == types.Method || kindName == types.Func {
		t, ok := reflector.GetFuncOutType(object)
		if ok {
			typ = t
		}
	}

	return &MetaData{
		Kind:      kindName,
		PkgName:   pkgName,
		TypeName:  typeName,
		Name:      name,
		ShortName: shortName,
		Context:   context,
		Object:    object,
		Type:      typ,
	}
}
