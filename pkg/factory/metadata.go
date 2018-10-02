package factory

import (
	"github.com/hidevopsio/hiboot/pkg/system/types"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
	"reflect"
	"strings"
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
	ExtDep    []*MetaData
}

// ParseParams parse object name and type
func ParseParams(postfix string, params ...interface{}) (metaData *MetaData) {

	hasTwoParams := len(params) == 2 && reflect.TypeOf(params[0]).Kind() == reflect.String

	var name string
	var object interface{}
	if hasTwoParams {
		object = params[1]
		name = params[0].(string)
	} else {
		object = params[0]
	}
	pkgName, typeName := reflector.GetPkgAndName(object)
	typ := reflect.TypeOf(object)
	kind := typ.Kind()
	kindName := kind.String()
	if kind == reflect.Struct && typ.Name() == types.Method {
		kindName = types.Method
	}
	if !hasTwoParams {
		name = strings.Replace(typeName, postfix, "", -1)
		name = str.ToLowerCamel(name)

		if name == "" || name == strings.ToLower(postfix) {
			name = pkgName
		}
	}

	metaData = &MetaData{
		Kind:      kindName,
		PkgName:   pkgName,
		TypeName:  typeName,
		Name:      pkgName + "." + name,
		ShortName: name,
		Object:    object,
	}

	return
}

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
	if kind == reflect.Struct && typ.Name() == types.Method {
		kindName = types.Method
	}
	if pkgName != "" {
		if name == "" {
			shortName = str.ToLowerCamel(typeName)
			name = pkgName + "." + shortName
		} else {
			shortName = name
			name = pkgName + "." + name
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
	}
}
