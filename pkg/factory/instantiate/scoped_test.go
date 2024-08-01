package instantiate_test

import (
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/factory/instantiate"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"github.com/hidevopsio/hiboot/pkg/utils/idgen"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type bar struct {
	at.Scope `value:"prototype"`
	ID       string
}

type baz struct {
	at.Scope `value:"prototype"`
	Name     string
}

type fooObj struct {
}

func (f *fooObj) Bar() *bar {
	id, err := idgen.NextString()
	if err == nil {
		return &bar{ID: id}
	}
	return &bar{ID: "0"}
}

type myConfig struct {
	at.ConditionalOnField `value:"Namespace,Name"`
	Namespace             string `json:"namespace"`
	Name                  string `json:"name"`
}
type scopedFuncObj struct {
	at.Scope `value:"prototype"`

	config *myConfig
}

func newScopedFuncObj(config *myConfig) *scopedFuncObj {
	log.Infof("myConfig: %v", config.Name)
	return &scopedFuncObj{config: config}
}

func TestScopedInstanceFactory(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	//log.Debug("methods: ", numOfMethod)
	testComponents := make([]*factory.MetaData, 0)

	// method
	fo := new(fooObj)
	ft := reflect.TypeOf(fo)
	method, ok := ft.MethodByName("Bar")
	assert.Equal(t, true, ok)
	testComponents = append(testComponents,
		factory.NewMetaData(fo, method),
	)

	ic := cmap.New()
	customProps := cmap.New()
	customProps.Set("app.project", "runtime-test")
	instFactory := instantiate.NewInstantiateFactory(ic, testComponents, customProps)
	instFactory.AppendComponent(newScopedFuncObj)
	_ = instFactory.BuildComponents()

	t.Run("should get singleton instance by default", func(t *testing.T) {
		type TestFoo struct {
			Name string
		}
		f := &instantiate.ScopedInstanceFactory[*TestFoo]{}
		result1 := f.GetInstance()
		assert.NotEqual(t, nil, result1)
		result2 := f.GetInstance()
		assert.NotEqual(t, nil, result2)
	})

	t.Run("should get scoped instance by default", func(t *testing.T) {
		type TestBar struct {
			at.Scope `value:"prototype"`
			Name     string
		}
		f := &instantiate.ScopedInstanceFactory[*TestBar]{}
		result1 := f.GetInstance()
		assert.NotEqual(t, nil, result1)
		result2 := f.GetInstance()
		assert.NotEqual(t, nil, result2)
	})

	t.Run("should get scoped instance each time", func(t *testing.T) {
		type FooService struct {
			factory *instantiate.ScopedInstanceFactory[*bar]
		}
		svc := &FooService{}
		svc.factory = &instantiate.ScopedInstanceFactory[*bar]{}
		result1 := svc.factory.GetInstance()
		result2 := svc.factory.GetInstance()
		assert.NotEqual(t, result1.ID, result2.ID)
	})

	t.Run("should get scoped instance with conditional name", func(t *testing.T) {
		type FooService struct {
			factory *instantiate.ScopedInstanceFactory[*scopedFuncObj]
		}
		svc := &FooService{}
		svc.factory = &instantiate.ScopedInstanceFactory[*scopedFuncObj]{}
		err := instFactory.SetInstance(&myConfig{Name: "default"})
		assert.Equal(t, nil, err)
		result0 := svc.factory.GetInstance()
		assert.Equal(t, "default", result0.config.Name)

		result1 := svc.factory.GetInstance(&myConfig{Name: "test1"})
		assert.Equal(t, "test1", result1.config.Name)

		result2 := svc.factory.GetInstance(&myConfig{Name: "test2"})
		assert.Equal(t, "test2", result2.config.Name)

		result3 := svc.factory.GetInstance(&myConfig{Name: "test2"})
		assert.Equal(t, "test2", result3.config.Name)

		result4 := svc.factory.GetInstance(&myConfig{Namespace: "dev", Name: "test4"})
		assert.Equal(t, "test4", result4.config.Name)
		assert.Equal(t, "dev", result4.config.Namespace)
	})
}
