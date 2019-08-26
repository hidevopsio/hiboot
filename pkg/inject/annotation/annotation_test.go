package annotation_test

import (
	"github.com/stretchr/testify/assert"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/inject/annotation"
	"reflect"
	"testing"
)


type AtFoo struct {
	at.Annotation
	Age int
}

type AtBar struct {
	at.Annotation
}

type AtFooBar struct {
	AtFoo
	Code int `value:"200"`
}

type AtFooBaz struct {
	AtFoo
	Code int `value:"400"`
}

type MyObj struct{
	Name string
	Value string
}

type foo struct {
	AtFoo `value:"foo,option 1,option 2" age:"18"`
	AtBar `value:"bar"`
	AtFooBar `value:"foobar" age:"12"`
	AtFooBaz `value:"baz" age:"22"`

	MyObj
}

type bar struct {
	AtFoo `value age:"25"`
	AtBar `value:"bar"`
}

func TestImplementsAnnotation(t *testing.T) {
	f := new(foo)
	f.Value = "my object value"

	t.Run("should check if object contains at.Annotation", func(t *testing.T) {
		ok := annotation.Contains(&f.AtFoo, at.Annotation{})
		assert.Equal(t, true, ok)
	})

	t.Run("should get annotation AtFoo", func(t *testing.T) {
		af, ok := annotation.GetField(f, AtFoo{})
		value, ok := af.Tag.Lookup("value")
		assert.Equal(t, "foo,option 1,option 2", value)
		assert.Equal(t, true, ok)

		age, ok := af.Tag.Lookup("age")
		assert.Equal(t, "18", age)
		assert.Equal(t, true, ok)
	})

	t.Run("should report error for invalid type that pass to GetFields", func(t *testing.T) {
		af := annotation.GetFields(123)
		assert.Equal(t, []reflect.StructField([]reflect.StructField(nil)), af)
	})

	t.Run("should report error for invalid type that pass to GetField", func(t *testing.T) {
		_, ok := annotation.GetField(123, AtFoo{})
		assert.Equal(t, false, ok)
	})

	t.Run("should inject all annotations", func(t *testing.T) {
		assert.Equal(t, "", f.AtFoo.Value)
		assert.Equal(t, 0, f.AtFoo.Age)
		err := annotation.InjectIntoObject(f)
		assert.Equal(t, nil, err)
		assert.Equal(t, "foo", f.AtFoo.Value)
		assert.Equal(t, 18, f.AtFoo.Age)

		assert.Equal(t, "bar", f.AtBar.Value)
		assert.Equal(t, "my object value", f.Value)
	})

	t.Run("should get annotation AtFoo", func(t *testing.T) {
		atFoo, ok := annotation.Get(f, AtFoo{})
		assert.Equal(t, true, ok)
		assert.Equal(t, "foo", atFoo.(AtFoo).Value)
	})

	t.Run("should get all annotations", func(t *testing.T) {
		as := annotation.GetAll(f)
		assert.NotEqual(t, 0, len(as))
	})

	t.Run("should find annotation AtFoo", func(t *testing.T) {
		as := annotation.Find(f, AtFoo{})
		assert.NotEqual(t, 0, len(as))
	})

	t.Run("should notify bad syntax for struct tag pair", func(t *testing.T) {
		// notify bad syntax for struct tag pair
		b := new(bar)
		err := annotation.InjectIntoObject(b)
		assert.NotEqual(t, nil, err)
		assert.Equal(t, "bad syntax for struct tag pair", err.Error())
	})

	t.Run("should inject to object", func(t *testing.T) {
		fo := foo{}
		err := annotation.InjectIntoObject(fo)
		assert.NotEqual(t, nil, err)
	})

	t.Run("should check if an object implements Annotation", func(t *testing.T) {
		ok := annotation.Contains(f, AtFoo{})
		assert.Equal(t, true, ok)
	})

	t.Run("should inject annotation into sub struct", func(t *testing.T) {
		var fb struct{at.GetMapping `value:"/path/to/api"`}
		err := annotation.InjectIntoObject(&fb)
		assert.Equal(t, nil, err)
	})

	t.Run("should get annotation by type", func(t *testing.T) {
		var fb struct{at.PostMapping `value:"/path/to/api"`}
		f := annotation.GetFields(reflect.TypeOf(&fb))
		assert.Equal(t, "PostMapping", f[0].Name)
	})

	t.Run("should report error for nil", func(t *testing.T) {
		err := annotation.InjectIntoObject(nil)
		assert.NotEqual(t, nil, err)
		assert.Equal(t, "object must not be nil", err.Error())
	})
}