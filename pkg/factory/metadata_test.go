package factory_test

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/system/types"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type fooBarService struct {
	HelloWorld HelloWorld `inject:""`
}

type Hello string
type HelloWorld string
type HelloHiboot string

type helloConfiguration struct {
	app.Configuration
}

func (c *helloConfiguration) Hello() Hello {
	return Hello("Hello World")
}

func (c *helloConfiguration) HelloWorld(h Hello) HelloWorld {
	return HelloWorld(h + "World")
}

func (c *helloConfiguration) HelloHiboot(h Hello) HelloHiboot {
	return HelloHiboot(h + "Hello Hiboot")
}

func newFooBarService() *fooBarService {
	return new(fooBarService)
}

func TestUtils(t *testing.T) {
	helloConfig := new(helloConfiguration)
	helloTyp := reflect.TypeOf(helloConfig)
	numOfMethod := helloTyp.NumMethod()
	//log.Debug("methods: ", numOfMethod)
	methodTestData := make([]*factory.MetaData, 0)
	for mi := 0; mi < numOfMethod; mi++ {
		method := helloTyp.Method(mi)
		// append inst to f.components
		methodTestData = append(methodTestData, factory.NewMetaData(helloConfig, method))
	}

	t.Run("should parse instance name via object", func(t *testing.T) {
		md := factory.NewMetaData("", new(fooBarService))
		assert.Equal(t, "fooBarService", md.TypeName)
		assert.Equal(t, "factory_test", md.PkgName)
		assert.NotEqual(t, nil, md.Object)
	})

	t.Run("should parse instance name via object", func(t *testing.T) {
		md := factory.NewMetaData("", new(fooBarService))
		assert.Equal(t, "fooBarService", md.TypeName)
		assert.Equal(t, "factory_test", md.PkgName)
		assert.NotEqual(t, nil, md.Object)
	})

	t.Run("should parse instance name via object with eliminator", func(t *testing.T) {
		md := factory.NewMetaData(new(fooBarService))
		assert.Equal(t, "fooBarService", md.TypeName)
		assert.Equal(t, "factory_test.fooBarService", md.Name)
		assert.NotEqual(t, nil, md.Object)
	})

	t.Run("should parse object instance name via constructor", func(t *testing.T) {
		md := factory.NewMetaData("", newFooBarService)
		assert.Equal(t, "fooBarService", md.TypeName)
		assert.Equal(t, types.Func, md.Kind)
	})

	t.Run("should parse object pkg name", func(t *testing.T) {
		type service struct{}
		svc := new(service)
		md := factory.NewMetaData(svc)
		assert.Equal(t, "factory_test.service", md.Name)
		assert.Equal(t, svc, md.Object)
	})

	t.Run("should parse object instance name", func(t *testing.T) {
		type service struct{}
		svc := new(service)
		md := factory.NewMetaData("foo", svc)
		assert.Equal(t, "service", md.TypeName)
		assert.Equal(t, "factory_test.foo", md.Name)
		assert.Equal(t, svc, md.Object)
	})

	t.Run("should parse object pkg name", func(t *testing.T) {
		type service struct{}
		svc := new(service)
		md := factory.NewMetaData(&factory.MetaData{Object: new(service)})
		assert.Equal(t, "factory_test.service", md.Name)
		assert.Equal(t, svc, md.Object)
	})

	t.Run("should parse object pkg name", func(t *testing.T) {
		type service struct{}
		svc := new(service)
		name, object := factory.ParseParams(svc)
		assert.Equal(t, "factory_test.service", name)
		assert.Equal(t, svc, object)
	})

	t.Run("should parse object pkg name", func(t *testing.T) {
		type service struct{}
		svc := new(service)
		name, object := factory.ParseParams("factory_test.fooService", svc)
		assert.Equal(t, "factory_test.fooService", name)
		assert.Equal(t, svc, object)
	})

	t.Run("should parse object pkg name", func(t *testing.T) {
		name, _ := factory.ParseParams("factory_test.fooService")
		assert.Equal(t, "factory_test.fooService", name)
	})
}
