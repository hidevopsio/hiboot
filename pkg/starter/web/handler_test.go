package web

import (
	"testing"
	"reflect"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestParse(t *testing.T) {

	hdl := new(handler)

	controller := new(FooController)
	ctrlVal := reflect.ValueOf(controller)

	t.Run("should parse method with path params", func(t *testing.T) {
		method, ok := ctrlVal.Type().MethodByName("PutByIdNameAge")
		assert.Equal(t, true, ok)
		hdl.parse(method, controller, "/foo/{id}/{name}/{age}")
		log.Debug(hdl)
		assert.Equal(t, 3, len(hdl.pathParams))
		assert.Equal(t, "FooController", hdl.requests[0].typeName)
		assert.Equal(t, "int", hdl.requests[1].typeName)
		assert.Equal(t, "string", hdl.requests[2].typeName)
		assert.Equal(t, "int", hdl.requests[3].typeName)
	})
}
