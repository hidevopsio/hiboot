package factory

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type fooBarService struct {
}

func newFooBarService() *fooBarService {
	return new(fooBarService)
}

func TestUtils(t *testing.T) {
	t.Run("should parse instance name via object", func(t *testing.T) {
		md := ParseParams("", new(fooBarService))
		assert.Equal(t, "fooBarService", md.TypeName)
		assert.Equal(t, "factory", md.PkgName)
		assert.NotEqual(t, nil, md.Object)
	})

	t.Run("should parse instance name via object with eliminator", func(t *testing.T) {
		md := ParseParams("Service", new(fooBarService))
		assert.Equal(t, "fooBarService", md.TypeName)
		assert.Equal(t, "fooBar", md.Name)
		assert.NotEqual(t, nil, md.Object)
	})

	t.Run("should parse object instance name via constructor", func(t *testing.T) {
		md := ParseParams("", newFooBarService)
		assert.Equal(t, "fooBarService", md.TypeName)
		assert.Equal(t, reflect.Func, md.Kind)
	})

	t.Run("should parse object pkg name", func(t *testing.T) {
		type service struct{}
		svc := new(service)
		md := ParseParams("service", svc)
		assert.Equal(t, "factory", md.Name)
		assert.Equal(t, svc, md.Object)
	})

	t.Run("should parse object instance name", func(t *testing.T) {
		type service struct{}
		svc := new(service)
		md := ParseParams("service", "foo", svc)
		assert.Equal(t, "service", md.TypeName)
		assert.Equal(t, "foo", md.Name)
		assert.Equal(t, svc, md.Object)
	})
}
