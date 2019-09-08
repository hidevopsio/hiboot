package swagger

import (
	"errors"
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
	if annotation.ContainsChild(atController.Fields, at.DisableSwagger{}) {
		log.Debugf("swagger is disabled by user on controller %v", atController.Value.Type())
		return
	}

	method, path, err := s.parseHttpMethod(atMethod)
	if err == nil {
		log.Debug("==================================================")
		log.Debugf("%v:%v", method, path)
		swAnn := annotation.Filter(atMethod.Fields, at.Swagger{})
		for _, a := range swAnn {
			ao := a.Value.Interface()
			log.Debug(ao)

			switch ao.(type) {
			case at.Operation:
				ann := ao.(at.Operation)
				key := ann.Key
				log.Debugf("key: %v", ann.Key)
				key = strings.Replace(key, "${at.http.method}", method, -1)
				key = strings.Replace(key, "${at.http.path}", path, -1)
				log.Debugf("key: %v", key)
				op := s.builder.GetProperty(key)
				log.Debug(op)
				s.builder.SetProperty(key, ann.OperationProps)
				op = s.builder.GetProperty(key)
				log.Debug(op)
			}

		}

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