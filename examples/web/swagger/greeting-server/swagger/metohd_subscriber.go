package swagger

import (
	"errors"
	"github.com/go-openapi/spec"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/inject/annotation"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/system"
	"strings"
)

var (
	// ErrHttpMethodNotFound method is not found error
	ErrHttpMethodNotFound = errors.New("HTTP method is not found")
)

type HttpMethod interface {
	GetMethod() string
	GetPath() string
}

type httpMethodSubscriber struct {
	at.HttpMethodSubscriber `value:"swagger"`
	builder system.Builder
}

func newHttpMethodSubscriber(builder  system.Builder) *httpMethodSubscriber {
	return &httpMethodSubscriber{
		builder: builder,
	}
}

// TODO: use data instead of atController
func (s *httpMethodSubscriber) Subscribe(atController *web.Annotations, atMethod *web.Annotations) {
	log.Debug("==================================================")

	// TODO: no need to use builder anymore according to the performance ?
	if !annotation.ContainsChild(atMethod.Fields, at.Operation{}) {
		log.Debugf("swagger is disabled by user on controller %v", atController.Value.Type())
		return
	}

	method, path, err := s.parseHttpMethod(atMethod)
	if err == nil {
		log.Debugf("%v:%v", method, path)
		swAnn := annotation.Filter(atMethod.Fields, at.Swagger{})

		////for test only
		//sw := s.builder.GetProperty("swagger")
		//log.Debug(sw)

		opAnn :=  annotation.Filter(atMethod.Fields, at.Operation{})
		ao := opAnn[0].Value.Interface()
		ann := ao.(at.Operation)
		var key string
		key = ann.Key
		key = strings.Replace(key, "${at.http.method}", method, -1)
		key = strings.Replace(key, "${at.http.path}", path, -1)
		key = strings.ToLower(key)
		pathItem := new(spec.PathItem)
		ato := ao.(at.Operation)
		operation := &ato.Operation

		//TODO: to be optimized
		switch method {
		case "GET":
			pathItem.Get = operation
		case "POST":
			pathItem.Post = operation
		}

		for _, a := range swAnn {
			ao := a.Value.Interface()
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

				rss := annotation.Filter(atMethod.Fields, at.ResponseSchema{})
				for _, rs := range rss {
					rso := rs.Value.Interface().(at.ResponseSchema)
					if rso.Code == ann.Code {
						ann.Response.Schema = &rso.Schema
						break
					}
				}
				operation.Responses.StatusCodeResponses[ann.Code] = ann.Response
			}
		}
		s.builder.SetProperty(key, operation)
	}
}

func (s *httpMethodSubscriber) parseHttpMethod(atMethod *web.Annotations) (method string, path string, err error) {
	// parse http method
	if atMethod.Object != nil {
		hma := annotation.Filter(atMethod.Fields, at.HttpMethod{})
		if len(hma) > 0 {
			hm := hma[0].Value.Interface()
			switch hm.(type) {
			case at.GetMapping:
				httpMethod := hm.(at.GetMapping)
				method, path = httpMethod.Method, httpMethod.Value
			case at.PostMapping:
				httpMethod := hm.(at.PostMapping)
				method, path = httpMethod.Method, httpMethod.Value
			case at.PutMapping:
				httpMethod := hm.(at.PutMapping)
				method, path = httpMethod.Method, httpMethod.Value
			case at.DeleteMapping:
				httpMethod := hm.(at.DeleteMapping)
				method, path = httpMethod.Method, httpMethod.Value
			case at.PatchMapping:
				httpMethod := hm.(at.PatchMapping)
				method, path = httpMethod.Method, httpMethod.Value
			case at.OptionsMapping:
				httpMethod := hm.(at.OptionsMapping)
				method, path = httpMethod.Method, httpMethod.Value
			case at.AnyMapping:
				httpMethod := hm.(at.AnyMapping)
				method, path = httpMethod.Method, httpMethod.Value
			case at.TraceMapping:
				httpMethod := hm.(at.TraceMapping)
				method, path = httpMethod.Method, httpMethod.Value
			default:
				err = ErrHttpMethodNotFound
			}
		} else {
			err = ErrHttpMethodNotFound
		}
	} else {
		err = ErrHttpMethodNotFound
	}
	return
}

func init() {
	app.Register(newHttpMethodSubscriber)
}