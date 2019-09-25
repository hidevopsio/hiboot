package swagger

import (
	"fmt"
	"github.com/go-openapi/spec"
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
	"time"
)

const refPrefix = "#/definitions/"

type apiPathsBuilder struct {
	apiInfoBuilder *apiInfoBuilder
	primitiveTypes map[string]string
}

func newApiPathsBuilder(builder *apiInfoBuilder) *apiPathsBuilder {
	if builder.SystemServer != nil {
		if builder.SystemServer.Host != "" {
			builder.SwaggerProps.Host = builder.SystemServer.Host
		}
		if builder.SystemServer.ContextPath != "" {
			builder.SwaggerProps.BasePath = builder.SystemServer.ContextPath
		}
		if len(builder.SystemServer.Schemes) > 0 {
			builder.SwaggerProps.Schemes = builder.SystemServer.Schemes
		}
	}
	if builder.AppVersion != "" {
		builder.Info.Version = builder.AppVersion
	}

	visit := fmt.Sprintf("%s://%s/swagger-ui", builder.SwaggerProps.Schemes[0], filepath.Join(builder.SwaggerProps.Host, builder.SwaggerProps.BasePath))
	log.Infof("visit %v to open api doc", visit)

	return &apiPathsBuilder{
		apiInfoBuilder: builder,
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
			"Time": "string",
		},
	}
}

func deepFields(reflectType reflect.Type) []reflect.StructField {
	var fields []reflect.StructField

	if reflectType = reflector.IndirectType(reflectType); reflectType.Kind() == reflect.Struct {
		for i := 0; i < reflectType.NumField(); i++ {
			v := reflectType.Field(i)

			if annotation.IsAnnotation(v.Type) {
				continue
			}

			if v.Anonymous {
				vk :=  reflector.IndirectType(v.Type).Kind()
				if vk == reflect.Struct || vk == reflect.Interface {
					fields = append(fields, deepFields(v.Type)...)
				}
			} else {
				fields = append(fields, v)
			}
		}
	}

	return fields
}

func (b *apiPathsBuilder) buildSchemaArray(definition *spec.Schema, typ reflect.Type)  {
	definition.Type = spec.StringOrArray{"array"}
	// array items
	arrSchema := spec.Schema{}
	arrType := typ.Elem()
	b.buildSchemaProperty(&arrSchema, arrType)
	definition.Items = &spec.SchemaOrArray{Schema: &arrSchema}
}

func (b *apiPathsBuilder) buildSchemaProperty(definition *spec.Schema, typ reflect.Type)  {
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

		for _, f := range deepFields(typ) {
			var jsonName string
			var fieldName string
			var descName string
			var descTag *structtag.Tag

			tags, _ := structtag.Parse(string(f.Tag))

			jsonTag, err := tags.Get("json")
			if err == nil {
				jsonName = jsonTag.Name
			}

			// first, check if field tag schema is presented
			schemaTag, err := tags.Get("schema")
			if err == nil {
				descTag = schemaTag
			} else {
				// assign field tag json to desc
				descTag = jsonTag
			}

			// assign schema properties
			typName := f.Type.Name()
			ps := spec.Schema{}
			if descTag != nil {
				// assign schema
				if jsonName != "" {
					fieldName = jsonName
				} else {
					fieldName = str.ToLowerCamel(f.Name)
				}

				if schemaTag == nil {
					descName = str.ToKebab(f.Name)
					descName = strings.Replace(descName, "-", " ", -1)
					descName = str.UpperFirst(descName)
				} else {
					descName = schemaTag.Name
				}

				ps.Title = f.Name
				ps.Description = descName
				ps.Format = typName
			}

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

			if descTag != nil {
				definition.Properties[fieldName] = ps
			}
		}
	}
}

func (b *apiPathsBuilder) buildSchemaObject(ps *spec.Schema, typ reflect.Type) {
	if typ == reflect.TypeOf(time.Time{}) {
		swgTypName := b.primitiveTypes[typ.Name()]
		ps.Type = spec.StringOrArray{swgTypName}
	} else {
		ps.Type = spec.StringOrArray{"object"}
		b.buildSchemaProperty(ps, typ)
	}
}

func (b *apiPathsBuilder) buildSchema(ann *annotation.Annotation, field *reflect.StructField) (schema *spec.Schema) {
	if field == nil {
		field = &ann.Field.StructField
	}

	atSchema := annotation.GetAnnotation(ann.Parent.Interface, at.Schema{})

	s := atSchema.Field.Value.Interface().(at.Schema)
	schemaType := s.AtType
	primitiveTypes := b.primitiveTypes[schemaType]

	schema = &spec.Schema{}
	if primitiveTypes == "" {
		err := annotation.Inject(atSchema)
		if err == nil {
			ref := refPrefix + field.Name
			// parse body schema and assign to definitions
			schema.Ref = spec.MustCreateRef(ref)

			if b.apiInfoBuilder.Definitions == nil {
				def := make(spec.Definitions)
				b.apiInfoBuilder.Definitions = def
			}

			definition, ok := b.apiInfoBuilder.Definitions[field.Name]
			if !ok {
				definition = spec.Schema{}
				b.buildSchemaProperty(&definition, field.Type)
				b.apiInfoBuilder.Definitions[field.Name] = definition
			}
		}
	} else {
		schema.Type = spec.StringOrArray{s.AtType}
		schema.Description = s.AtDescription
	}

	return
}

func (b *apiPathsBuilder) buildParameter(operation *spec.Operation, annotations *annotation.Annotations, a *annotation.Annotation) {
	ao := a.Field.Value.Interface()
	atParameter := ao.(at.Parameter)
	// copy values
	parameter := spec.Parameter{}
	parameter.Name = atParameter.AtName
	parameter.Type = atParameter.AtType
	parameter.In = atParameter.AtIn
	parameter.Description = atParameter.AtDescription

	if atParameter.AtIn == "body" || atParameter.AtIn == "array" {

		atSchema := annotation.Find(annotations, at.Schema{})

		if atSchema != nil {

			field := b.findArrayField(atSchema)

			parameter.Schema = b.buildSchema(atSchema, field)
		}

	}

	operation.Parameters = append(operation.Parameters, parameter)
	return
}

func (b *apiPathsBuilder) findArrayField(schema *annotation.Annotation) (field *reflect.StructField) {
	parentType := schema.Parent.Type
	numField := parentType.NumField()
	for i := 0; i < numField; i++ {
		f := parentType.Field(i)
		nextIndex := f.Index[0] + 1
		if f.Type == reflect.TypeOf(at.Schema{}) && nextIndex < numField {
			nextField := parentType.Field(f.Index[0] + 1)
			if nextField.Type.Kind() == reflect.Slice {
				field = &nextField
				break
			}
		}
	}
	return field
}

func (b *apiPathsBuilder) buildResponse(operation *spec.Operation, annotations *annotation.Annotations, a *annotation.Annotation) {
	ao := a.Field.Value.Interface()
	atResponse := ao.(at.Response)
	if operation.Responses == nil {
		operation.Responses = new(spec.Responses)
		operation.Responses.StatusCodeResponses = make(map[int]spec.Response)
	}
	atSchema := annotation.Find(annotations, at.Schema{})

	response := spec.Response{}
	response.Description = atResponse.AtDescription
	if atSchema != nil {
		field := b.findArrayField(atSchema)

		response.Schema = b.buildSchema(atSchema, field)
	}

	operation.Responses.StatusCodeResponses[atResponse.AtCode] = response
	return
}

func (b *apiPathsBuilder) buildOperation(operation *spec.Operation, annotations *annotation.Annotations)  {
	for _, a := range annotations.Items {
		ao := a.Field.Value.Interface()
		switch ao.(type) {
		case at.Consumes:
			ann := ao.(at.Consumes)
			operation.Consumes = append(operation.Consumes, ann.AtValues...)
		case at.Produces:
			ann := ao.(at.Produces)
			operation.Produces = append(operation.Produces, ann.AtValues...)
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

func (b *apiPathsBuilder) Build(atController *annotation.Annotations, atMethod *annotation.Annotations) {

	if !annotation.ContainsChild(atMethod, at.Operation{}) {
		//log.Debugf("does not found any swagger annotations in %v", atController.Items[0].Parent.Type)
		return
	}

	method, path := webutils.GetHttpMethod(atMethod)
	if method != "" {
		atRequestMapping := annotation.GetAnnotation(atController, at.RequestMapping{})
		if atRequestMapping != nil {
			ann := atRequestMapping.Field.Value.Interface().(at.RequestMapping)
			path = filepath.Join(ann.AtValue, path)
		}
		//log.Debugf("%v:%v", method, path)

		pathItem := b.apiInfoBuilder.Paths.Paths[path]

		ann :=  annotation.GetAnnotation(atMethod, at.Operation{})

		atOperationInterface := ann.Field.Value.Interface()
		atOperation := atOperationInterface.(at.Operation)

		// copy values
		operation := &spec.Operation{}
		operation.ID = atOperation.AtID
		operation.Description = atOperation.AtDescription

		method = strings.Title(strings.ToLower(method))
		err := reflector.SetFieldValue(&pathItem, method, operation)
		if err == nil {
			b.buildOperation(operation, atMethod)

			// add new path item
			b.apiInfoBuilder.Paths.Paths[path] = pathItem
			//log.Debug(b.openAPIDefinition.Paths.Paths[path])
		}
	}
}
