package factory

import (
	"github.com/hidevopsio/hiboot/pkg/system/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

type fooBarService struct {
}

func newFooBarService() *fooBarService {
	return new(fooBarService)
}

func TestUtils(t *testing.T) {
	t.Run("should parse instance name via object", func(t *testing.T) {
		md := NewMetaData("", new(fooBarService))
		assert.Equal(t, "fooBarService", md.TypeName)
		assert.Equal(t, "factory", md.PkgName)
		assert.NotEqual(t, nil, md.Object)
	})

	t.Run("should parse instance name via object with eliminator", func(t *testing.T) {
		md := NewMetaData(new(fooBarService))
		assert.Equal(t, "fooBarService", md.TypeName)
		assert.Equal(t, "factory.fooBarService", md.Name)
		assert.NotEqual(t, nil, md.Object)
	})

	t.Run("should parse object instance name via constructor", func(t *testing.T) {
		md := NewMetaData("", newFooBarService)
		assert.Equal(t, "fooBarService", md.TypeName)
		assert.Equal(t, types.Func, md.Kind)
	})

	t.Run("should parse object pkg name", func(t *testing.T) {
		type service struct{}
		svc := new(service)
		md := NewMetaData(svc)
		assert.Equal(t, "factory.service", md.Name)
		assert.Equal(t, svc, md.Object)
	})

	t.Run("should parse object instance name", func(t *testing.T) {
		type service struct{}
		svc := new(service)
		md := NewMetaData("foo", svc)
		assert.Equal(t, "service", md.TypeName)
		assert.Equal(t, "factory.foo", md.Name)
		assert.Equal(t, svc, md.Object)
	})
}
