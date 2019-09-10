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
	"path/filepath"
	"strings"
)


type pathsBuilder struct {
	openAPIDefinition *openAPIDefinition
}

func newOpenAPIDefinitionBuilder(openAPIDefinition *openAPIDefinition) *pathsBuilder {
	if openAPIDefinition == nil {
		log.Fatal(`

// Please register swagger.OpenAPIDefinitionBuilder(), see below code snippet for reference

func init() {
	app.Register(swagger.OpenAPIDefinitionBuilder().
		Version("1.0.0").
		Title("HiBoot Swagger Demo Application - Greeting Server").
		Description("Greeting Server is an application that demonstrate the usage of Swagger Annotations").
		Schemes("http", "https").
		Host("apps.hidevops.io").
		BasePath("/api/v1/greeting-server"),
	)
}

`)
		return nil
	}

	swgProp := openAPIDefinition.SwaggerProps
	visit := fmt.Sprintf("%s://%s%s/swagger-ui", swgProp.Schemes[0], swgProp.Host, swgProp.BasePath)
	log.Infof("visit open api doc: %v", visit)

	return &pathsBuilder{openAPIDefinition: openAPIDefinition}
}

func init() {
	app.Register(newOpenAPIDefinitionBuilder)
}

func (b *pathsBuilder) buildOperation(operation *spec.Operation, annotations *annotation.Annotations)  {
	for _, a := range annotations.Items {
		ao := a.Field.Value.Interface()
		switch ao.(type) {
		case at.Parameter:
			ann := ao.(at.Parameter)
			operation.Parameters = append(operation.Parameters, ann.Parameter)
		case at.Produces:
			ann := ao.(at.Produces)
			operation.Produces = append(operation.Produces, ann.Values...)
		case at.Response:
			ann := ao.(at.Response)
			if operation.Responses == nil {
				operation.Responses = new(spec.Responses)
				operation.Responses.StatusCodeResponses = make(map[int]spec.Response)
			}

			atSchema := annotation.GetAnnotation(annotations, at.Schema{})
			if atSchema != nil {
				atSchemaObj := atSchema.Field.Value.Interface().(at.Schema)
				ann.Response.Schema = &atSchemaObj.Schema
			}

			operation.Responses.StatusCodeResponses[ann.Code] = ann.Response
		}
	}

	for _, child := range annotations.Children {
		b.buildOperation(operation, child)
	}
}


func (b *pathsBuilder) Build(atController *annotation.Annotations, atMethod *annotation.Annotations) {

	if !annotation.ContainsChild(atMethod, at.Operation{}) {
		log.Debugf("does not found any swagger annotations in %v", atController.Items[0].Parent.Type)
		return
	}

	method, path := webutils.GetHttpMethod(atMethod)
	if method != "" {
		pathItem := spec.PathItem{}

		atRequestMapping := annotation.GetAnnotation(atController, at.RequestMapping{})
		if atRequestMapping != nil {
			ann := atRequestMapping.Field.Value.Interface().(at.RequestMapping)
			path = filepath.Join(ann.Value, path)
		}
		log.Debugf("%v:%v", method, path)

		atOperation :=  annotation.GetAnnotation(atMethod, at.Operation{})

		atOperationInterface := atOperation.Field.Value.Interface()
		atOperationObject := atOperationInterface.(at.Operation)
		operation := &atOperationObject.Operation

		err := reflector.SetFieldValue(&pathItem, strings.Title(strings.ToLower(method)), operation)
		if err == nil {
			b.buildOperation(operation, atMethod)

			// add new path item
			b.openAPIDefinition.Paths.Paths[strings.ToLower(path)] = pathItem
		}
	}
}
