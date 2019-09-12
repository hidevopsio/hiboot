package swagger

import (
	"fmt"
	"github.com/go-openapi/spec"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web/webutils"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/inject/annotation"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/utils/reflector"
	"hidevops.io/hiboot/pkg/utils/str"
	"hidevops.io/hiboot/pkg/utils/structtag"
	"path/filepath"
	"reflect"
	"strings"
)

type pathsBuilder struct {
	openAPIDefinition *openAPIDefinition
	primitiveTypes map[string]string
}

func newOpenAPIDefinitionBuilder(openAPIDefinition *openAPIDefinition) *pathsBuilder {
	if openAPIDefinition.SystemServer != nil {
		if openAPIDefinition.SwaggerProps.Host == "" && openAPIDefinition.SystemServer.Host != "" {
			openAPIDefinition.SwaggerProps.Host = openAPIDefinition.SystemServer.Host
		}
		if openAPIDefinition.SwaggerProps.BasePath == "" && openAPIDefinition.SystemServer.ContextPath != "" {
			openAPIDefinition.SwaggerProps.BasePath = openAPIDefinition.SystemServer.ContextPath
		}
		if len(openAPIDefinition.SwaggerProps.Schemes) == 0 && len(openAPIDefinition.SystemServer.Schemes) > 0 {
			openAPIDefinition.SwaggerProps.Schemes = openAPIDefinition.SystemServer.Schemes
		}
	}
	if openAPIDefinition.Info.Version == "" && openAPIDefinition.AppVersion != "" {
		openAPIDefinition.Info.Version = openAPIDefinition.AppVersion
	}

	visit := fmt.Sprintf("%s://%s/swagger-ui", openAPIDefinition.SwaggerProps.Schemes[0], filepath.Join(openAPIDefinition.SwaggerProps.Host, openAPIDefinition.SwaggerProps.BasePath))
	log.Infof("visit %v to open api doc", visit)

	return &pathsBuilder{
		openAPIDefinition: openAPIDefinition,
		primitiveTypes: map[string]string{
			// array, boolean, integer, number, object, string
			"string": "string",
			"int": "integer",
			"int8": "integer",
			"int16": "integer",
			"int32": "integer",
			"int64": "integer",
			"uint": "integer",
			"uint8": "integer",
			"uint16": "integer",
			"uint32": "integer",
			"uint64": "integer",
			"float32": "number",
			"float64": "number",
			"struct": "object",
			"slice": "array",
			"bool": "boolean",
		},
	}
}

func init() {
	app.Register(newOpenAPIDefinitionBuilder)
}

func (b *pathsBuilder) buildSchemaArray(definition *spec.Schema, typ reflect.Type)  {
	definition.Type = spec.StringOrArray{"array"}
	// array items
	arrSchema := spec.Schema{}
	arrType := typ.Elem()
	b.buildSchemaProperty(&arrSchema, arrType)
	definition.Items = &spec.SchemaOrArray{Schema: &arrSchema}
}

func (b *pathsBuilder) buildSchemaProperty(definition *spec.Schema, typ reflect.Type)  {
	kind := typ.Kind()
	if kind == reflect.Ptr {
		typ = reflector.IndirectType(typ)
		kind = typ.Kind()
	}

	if kind == reflect.Slice {
		b.buildSchemaArray(definition, typ)
	} else if kind == reflect.Struct {
		definition.Properties = make(map[string]spec.Schema)
		definition.Type = spec.StringOrArray{"object"}

		for _, f := range reflector.DeepFields(typ) {
			desc, ok := f.Tag.Lookup("schema")
			if ok {
				typName := f.Type.Name()
				ps := spec.Schema{}
				ps.Title = f.Name
				ps.Description = desc
				ps.Format = typName
				fieldKind := f.Type.Kind()

				switch fieldKind {
				case reflect.Slice:
					b.buildSchemaArray(&ps, f.Type)
				case reflect.Struct:
					b.buildSchemaObject(&ps, f.Type)
				case reflect.Ptr:
					iTyp := reflector.IndirectType(f.Type)
					if iTyp.Kind() == reflect.Struct {
						b.buildSchemaObject(&ps, iTyp)
					}
				default:
					// convert primitive types
					swgTypName := b.primitiveTypes[typName]
					ps.Type = spec.StringOrArray{swgTypName}
				}

				// assign schema
				fieldName := b.getFieldName(f)
				definition.Properties[fieldName] = ps
			}
		}
	}
}

func (b *pathsBuilder) buildSchemaObject(ps *spec.Schema, typ reflect.Type) {
	ps.Type = spec.StringOrArray{"object"}
	childSchema := annotation.GetAnnotation(typ, at.Schema{})
	if childSchema != nil {
		b.buildSchemaProperty(ps, typ)
	}
}

func (b *pathsBuilder) getFieldName(f reflect.StructField) string {
	tags, err := structtag.Parse(string(f.Tag))
	var fieldName string
	if err == nil {
		tag, err := tags.Get("json")
		if err == nil {
			fieldName = tag.Name
		} else {
			fieldName = str.ToLowerCamel(f.Name)
		}
	}
	return fieldName
}

func (b *pathsBuilder) buildSchema(ann *annotation.Annotation, field *reflect.StructField) (schema *spec.Schema) {
	if field == nil {
		field = &ann.Field.StructField
	}

	atSchema := annotation.GetAnnotation(ann.Parent.Interface, at.Schema{})
	err := annotation.Inject(atSchema)
	if err == nil {
		s := atSchema.Field.Value.Interface().(at.Schema)
		ref := "#/definitions/" + field.Name
		s.Ref = spec.MustCreateRef(ref)

		// parse body schema and assign to definitions
		if b.openAPIDefinition.Definitions == nil {
			def := make(spec.Definitions)
			b.openAPIDefinition.Definitions = def
		}

		definition := spec.Schema{}
		b.buildSchemaProperty(&definition, field.Type)
		b.openAPIDefinition.Definitions[field.Name] = definition

		schema = &s.Schema
	}
	return
}

func (b *pathsBuilder) buildParameter(operation *spec.Operation, annotations *annotation.Annotations, a *annotation.Annotation) {
	ao := a.Field.Value.Interface()
	atParam := ao.(at.Parameter)
	if atParam.In == "body" || atParam.In == "array" {

		schema := annotation.Find(annotations, at.Schema{})

		if schema != nil {

			field := b.findArrayField(schema)

			atParam.Parameter.Schema = b.buildSchema(schema, field)
		}

	}

	operation.Parameters = append(operation.Parameters, atParam.Parameter)
	return
}

func (b *pathsBuilder) findArrayField(schema *annotation.Annotation) (field *reflect.StructField) {
	parentType := schema.Parent.Type
	var foundSchema bool
	for i := 0; i < parentType.NumField(); i++ {
		f := parentType.Field(i)
		if f.Type == reflect.TypeOf(at.Schema{}) {
			foundSchema = true
			schemaField := schema.Parent.Value.FieldByName(f.Name)
			atSchema := schemaField.Interface().(at.Schema)
			foundSchema = atSchema.Value == "array"
			continue
		}

		if foundSchema && f.Type.Kind() == reflect.Slice {
			field = &f
			break
		}
	}
	return field
}

func (b *pathsBuilder) buildResponse(operation *spec.Operation, annotations *annotation.Annotations, a *annotation.Annotation) {
	ao := a.Field.Value.Interface()
	atResp := ao.(at.Response)
	if operation.Responses == nil {
		operation.Responses = new(spec.Responses)
		operation.Responses.StatusCodeResponses = make(map[int]spec.Response)
	}
	schema := annotation.Find(annotations, at.Schema{})

	if schema != nil {
		field := b.findArrayField(schema)

		atResp.Response.Schema = b.buildSchema(schema, field)
	}

	operation.Responses.StatusCodeResponses[atResp.Code] = atResp.Response
	return
}

func (b *pathsBuilder) buildOperation(operation *spec.Operation, annotations *annotation.Annotations)  {
	for _, a := range annotations.Items {
		ao := a.Field.Value.Interface()
		switch ao.(type) {
		case at.Consumes:
			ann := ao.(at.Consumes)
			operation.Consumes = append(operation.Consumes, ann.Values...)
		case at.Produces:
			ann := ao.(at.Produces)
			operation.Produces = append(operation.Produces, ann.Values...)
		case at.Parameter:
			b.buildParameter(operation, annotations, a)
		case at.Response:
			b.buildResponse(operation, annotations, a)
		}
	}

	for _, child := range annotations.Children {
		b.buildOperation(operation, child)
	}
}


func (b *pathsBuilder) Build(atController *annotation.Annotations, atMethod *annotation.Annotations) {

	if !annotation.ContainsChild(atMethod, at.Operation{}) {
		//log.Debugf("does not found any swagger annotations in %v", atController.Items[0].Parent.Type)
		return
	}

	method, path := webutils.GetHttpMethod(atMethod)
	if method != "" {
		atRequestMapping := annotation.GetAnnotation(atController, at.RequestMapping{})
		if atRequestMapping != nil {
			ann := atRequestMapping.Field.Value.Interface().(at.RequestMapping)
			path = filepath.Join(ann.Value, path)
		}
		//log.Debugf("%v:%v", method, path)

		pathItem := b.openAPIDefinition.Paths.Paths[path]

		atOperation :=  annotation.GetAnnotation(atMethod, at.Operation{})

		atOperationInterface := atOperation.Field.Value.Interface()
		atOperationObject := atOperationInterface.(at.Operation)
		operation := &atOperationObject.Operation

		method = strings.Title(strings.ToLower(method))
		err := reflector.SetFieldValue(&pathItem, method, operation)
		if err == nil {
			b.buildOperation(operation, atMethod)

			// add new path item
			//path = strings.ToLower(path)
			b.openAPIDefinition.Paths.Paths[path] = pathItem
			//log.Debug(b.openAPIDefinition.Paths.Paths[path])
		}
	}
}
