package factory

import (
	"github.com/hidevopsio/hiboot/pkg/system/types"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
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
	Type      reflect.Type
	Depends   []string
	ExtDep    []*MetaData
}

func appendDep(deps, dep string) (retVal string) {
	if deps == "" {
		retVal = dep
	} else {
		retVal = deps + "," + dep
	}
	return
}

func findDep(objTyp, inTyp reflect.Type) (name string) {
	indInTyp := reflector.IndirectType(inTyp)
	for _, field := range reflector.DeepFields(objTyp) {
		indFieldTyp := reflector.IndirectType(field.Type)
		//log.Debugf("%v <> %v", indFieldTyp, indInTyp)
		if indFieldTyp == indInTyp {
			name = str.ToLowerCamel(field.Name)
			depPkgName := io.DirName(field.Type.PkgPath())
			if depPkgName != "" {
				name = depPkgName + "." + name
			}
			break
		}
	}
	if name == "" {
		name = reflector.GetLowerCamelFullNameByType(inTyp)
	}
	return
}

func parseDependencies(object interface{}, kind string, typ reflect.Type) (deps []string) {
	var depNames string
	switch kind {
	case types.Func:
		fn := reflect.ValueOf(object)
		numIn := fn.Type().NumIn()
		for i := 0; i < numIn; i++ {
			inTyp := fn.Type().In(i)
			depNames = appendDep(depNames, findDep(typ, inTyp))
		}
	case types.Method:
		// TODO: too many duplicated code, optimize it
		method := object.(reflect.Method)
		numIn := method.Type.NumIn()
		for i := 1; i < numIn; i++ {
			inTyp := method.Type.In(i)
			depNames = appendDep(depNames, findDep(typ, inTyp))
		}
	default:
		// find user specific inject tag
		for _, field := range reflector.DeepFields(typ) {
			tag, ok := field.Tag.Lookup("inject")
			if ok {
				name := tag
				if name == "" {
					name = str.ToLowerCamel(field.Type.Name())
				}
				depNames = appendDep(depNames, name)
			}
		}

		// find user specific depends tag
		var depTag string
		var hasDepTag bool
		f := reflector.GetEmbeddedField(object, reflect.Struct)
		depTag, hasDepTag = f.Tag.Lookup("depends")
		if hasDepTag {
			depNames = appendDep(depNames, depTag)
		}
	}

	if depNames != "" {
		deps = strings.Split(depNames, ",")
	}
	return
}

// ParseParams parse parameters
func ParseParams(params ...interface{}) (name string, object interface{}) {
	numParams := len(params)
	if numParams != 0 {
		kind := reflect.TypeOf(params[0]).Kind()
		if numParams == 1 {
			if kind == reflect.String {
				name = str.LowerFirst(params[0].(string))
				return
			} else {
				object = params[0]
			}
		} else {
			if kind == reflect.String {
				name = str.LowerFirst(params[0].(string))
				object = params[1]
			} else {
				object = params[1]
			}
		}

		pkgName, typeName := reflector.GetPkgAndName(object)
		if pkgName != "" {
			if name == "" {
				name = pkgName + "." + str.ToLowerCamel(typeName)
			} else if !strings.Contains(name, ".") {
				name = pkgName + "." + name
			}
		}
	}

	return
}

// NewMetaData create new meta data
func NewMetaData(params ...interface{}) *MetaData {
	var name string
	var shortName string
	var object interface{}
	var context interface{}
	var deps []string

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

	switch object.(type) {
	case *MetaData:
		md := object.(*MetaData)
		deps = append(deps, md.Depends...)
		object = md.Object
		name = md.Name
	}

	pkgName, typeName := reflector.GetPkgAndName(object)
	typ := reflect.TypeOf(object)
	kind := typ.Kind()
	kindName := kind.String()

	if pkgName != "" {
		shortName = str.ToLowerCamel(typeName)
		if name == "" {
			name = pkgName + "." + shortName
		} else if !strings.Contains(name, ".") {
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

	deps = append(deps, parseDependencies(object, kindName, typ)...)

	metaData := &MetaData{
		Kind:      kindName,
		PkgName:   pkgName,
		TypeName:  typeName,
		Name:      name,
		ShortName: shortName,
		Context:   context,
		Object:    object,
		Type:      typ,
		Depends:   deps,
	}

	return metaData
}
