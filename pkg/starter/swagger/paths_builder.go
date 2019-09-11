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
}

func newOpenAPIDefinitionBuilder(openAPIDefinition *openAPIDefinition) *pathsBuilder {
	if openAPIDefinition.SystemServer != nil {
		if openAPIDefinition.SwaggerProps.Host == "" {
			openAPIDefinition.SwaggerProps.Host = openAPIDefinition.SystemServer.Host
		}
		if openAPIDefinition.SwaggerProps.BasePath == "" {
			openAPIDefinition.SwaggerProps.BasePath = openAPIDefinition.SystemServer.ContextPath
		}
		if openAPIDefinition.SwaggerProps.Schemes == nil {
			openAPIDefinition.SwaggerProps.Schemes = openAPIDefinition.SystemServer.Schemes
		}
	}
	if openAPIDefinition.Info.Version == "" &&  openAPIDefinition.AppVersion != "" {
		openAPIDefinition.Info.Version = openAPIDefinition.AppVersion
	}

	visit := fmt.Sprintf("%s://%s/swagger-ui", openAPIDefinition.SwaggerProps.Schemes[0], filepath.Join(openAPIDefinition.SwaggerProps.Host, openAPIDefinition.SwaggerProps.BasePath))
	log.Infof("visit %v to open api doc", visit)

	return &pathsBuilder{openAPIDefinition: openAPIDefinition}
}

func init() {
	app.Register(newOpenAPIDefinitionBuilder)
}

func (b *pathsBuilder) buildSchemaArray(definition *spec.Schema, typ reflect.Type)  {
	definition.Type = spec.StringOrArray{"array"}
	// array items
	arrSchema := spec.Schema{}
	arrType := typ.Elem()
	b.buildSchema(&arrSchema, arrType)
	definition.Items = &spec.SchemaOrArray{Schema: &arrSchema}
}

func (b *pathsBuilder) buildSchema(definition *spec.Schema, typ reflect.Type)  {
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
				ps := spec.Schema{}
				ps.Title = f.Name
				ps.Description = desc
				if f.Type.Kind() == reflect.Slice {
					b.buildSchemaArray(&ps, f.Type)
				} else if f.Type.Kind() == reflect.Struct {
					ps.Type = spec.StringOrArray{"object"}
					childSchema := annotation.GetAnnotation(f.Type, at.Schema{})
					if childSchema != nil {
						b.buildSchema(&ps, f.Type)
					}
				} else {
					ps.Type = spec.StringOrArray{f.Type.Name()}
				}

				// assign schema
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
				definition.Properties[fieldName] = ps
			}
		}
	}
}

func (b *pathsBuilder) buildSchemaBody(body *annotation.Annotation,) (schema *spec.Schema) {
	atSchema := annotation.GetAnnotation(body.Parent.Interface, at.Schema{})
	err := annotation.Inject(atSchema)
	if err == nil {
		s := atSchema.Field.Value.Interface().(at.Schema)
		ref := "#/definitions/" + body.Field.StructField.Name
		s.Ref = spec.MustCreateRef(ref)

		// parse body schema and assign to definitions
		if b.openAPIDefinition.Definitions == nil {
			def := make(spec.Definitions)
			b.openAPIDefinition.Definitions = def
		}

		definition := spec.Schema{}
		b.buildSchema(&definition, body.Field.StructField.Type)
		b.openAPIDefinition.Definitions[body.Field.StructField.Name] = definition

		schema = &s.Schema
	}
	return
}

func (b *pathsBuilder) buildParameter(operation *spec.Operation, annotations *annotation.Annotations, a *annotation.Annotation) {
	ao := a.Field.Value.Interface()
	atParam := ao.(at.Parameter)
	switch atParam.In {
	case "body":
		log.Debug("body")
		body := annotation.Find(annotations, at.Schema{})

		atParam.Parameter.Schema = b.buildSchemaBody(body)
	}

	operation.Parameters = append(operation.Parameters, atParam.Parameter)
	return
}

func (b *pathsBuilder) buildResponse(operation *spec.Operation, annotations *annotation.Annotations, a *annotation.Annotation) {
	ao := a.Field.Value.Interface()
	atResp := ao.(at.Response)
	if operation.Responses == nil {
		operation.Responses = new(spec.Responses)
		operation.Responses.StatusCodeResponses = make(map[int]spec.Response)
	}
	body := annotation.Find(annotations, at.Schema{})
	if body != nil {
		atResp.Response.Schema = b.buildSchemaBody(body)
	}

	operation.Responses.StatusCodeResponses[atResp.Code] = atResp.Response
	return
}

func (b *pathsBuilder) buildOperation(operation *spec.Operation, annotations *annotation.Annotations)  {
	for _, a := range annotations.Items {
		ao := a.Field.Value.Interface()
		switch ao.(type) {
		case at.Parameter:
			b.buildParameter(operation, annotations, a)
		case at.Consumes:
			ann := ao.(at.Consumes)
			operation.Consumes = append(operation.Consumes, ann.Values...)
		case at.Produces:
			ann := ao.(at.Produces)
			operation.Produces = append(operation.Produces, ann.Values...)
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
