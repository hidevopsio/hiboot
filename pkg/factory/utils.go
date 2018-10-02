package factory

import (
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
	"reflect"
	"strings"
)

// ParseParams parse object name and type
func ParseParams(eliminator string, params ...interface{}) (metaData *MetaData) {

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
	kind := reflect.TypeOf(object).Kind()
	if !hasTwoParams {
		name = strings.Replace(typeName, eliminator, "", -1)
		name = str.ToLowerCamel(name)

		if name == "" || name == strings.ToLower(eliminator) {
			name = pkgName
		}
	}

	metaData = &MetaData{
		Kind:     kind,
		PkgName:  pkgName,
		TypeName: typeName,
		//Name:     pkgName + "." + name,
		Name:     name,
		Object:   object,
	}

	return
}
