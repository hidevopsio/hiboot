package annotation_test

import (
	"github.com/stretchr/testify/assert"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/inject/annotation"
	"hidevops.io/hiboot/pkg/log"
	"reflect"
	"testing"
)

type AtBaz struct {
	at.Annotation
	Code int `value:"200" json:"code"`
}

type AtFoo struct {
	at.Annotation
	ID int `json:"fooId"`
	Age int `json:"age"`
}

type AtBar struct {
	at.Annotation
}

type AtFooBar struct {
	AtFoo
	Code int `value:"200" json:"code"`
}

type AtFooBaz struct {
	AtFoo
	Code int `value:"400" json:"code"`
}

type MyObj struct{
	Name string
	Value string
}

type foo struct {
	AtBaz `value:"baz"`
	AtFoo `value:"foo,option 1,option 2" age:"18" fooId:"123"`
	AtBar `value:"bar"`
	AtFooBar `value:"foobar" age:"12"`
	AtFooBaz `value:"foobaz" age:"22"`
	MyObj
}

type bar struct {
	AtFoo `value age:"25"`
	AtBar `value:"bar"`
}

type multipleBar struct {
	AtFoo `value:"foo"`
	Bar1 struct{
		AtBar `value:"bar1"`
	}
	Bar2 struct{
		AtBar `value:"bar2"`
	}
}

type AtIntMap struct {
	at.Annotation
	FieldName string `value:"codes"`
	Codes map[int]string
}

type AtStrMap struct {
	at.Annotation
	FieldName string `value:"messages"`
	Messages map[string]string
}

type foobar struct {
	AtIntMap `200:"success" 404:"not found" 403:"unauthorized"`
	AtStrMap `ok:"successful" failed:"failed"`
}

type atApiOperation struct {
	at.ApiOperation `value:"testApi" operationId:"getGreeting" description:"This is the Greeting api for demo"`
}

func TestImplementsAnnotation(t *testing.T) {
	log.SetLevel("debug")

	f := new(foo)
	f.Value = "my object value"
	t.Run("test api operation", func(t *testing.T) {
		ao := &atApiOperation{}
		err := annotation.InjectIntoFields(ao)
		assert.Equal(t, nil, err)
	})

	t.Run("test map injection", func(t *testing.T) {
		fb := &foobar{}
		err := annotation.InjectIntoFields(fb)
		assert.Equal(t, nil, err)
	})

	fields := annotation.GetFields(f)
	t.Run("should check if object contains at.Annotation", func(t *testing.T) {
		assert.NotEqual(t, nil, fields)

		ok := annotation.Contains(fields, AtFoo{})
		assert.Equal(t, true, ok)
	})

	t.Run("should find if contains child annotation", func(t *testing.T) {
		ok := annotation.ContainsChild(fields, AtFoo{})
		assert.Equal(t, true, ok)
	})

	t.Run("should check if object contains at.Annotation", func(t *testing.T) {
		ok := annotation.Contains(&f.AtFoo, at.Annotation{})
		assert.Equal(t, true, ok)
	})

	t.Run("should get annotation AtFoo", func(t *testing.T) {
		af, ok := annotation.GetField(f, AtFoo{})
		value, ok := af.StructField.Tag.Lookup("value")
		assert.Equal(t, "foo,option 1,option 2", value)
		assert.Equal(t, true, ok)

		age, ok := af.StructField.Tag.Lookup("age")
		assert.Equal(t, "18", age)
		assert.Equal(t, true, ok)
	})

	t.Run("should report error for invalid type that pass to GetFields", func(t *testing.T) {
		af := annotation.GetFields(123)
		assert.Equal(t, []*annotation.Field([]*annotation.Field(nil)), af)
	})

	t.Run("should report error for invalid type that pass to GetField", func(t *testing.T) {
		_, ok := annotation.GetField(123, AtFoo{})
		assert.Equal(t, false, ok)
	})

	t.Run("should inject all annotations", func(t *testing.T) {
		assert.Equal(t, "", f.AtFoo.Value)
		assert.Equal(t, 0, f.AtFoo.Age)
		err := annotation.InjectIntoFields(f)
		assert.Equal(t, nil, err)
		assert.Equal(t, "foo", f.AtFoo.Value)
		assert.Equal(t, 18, f.AtFoo.Age)

		assert.Equal(t, "bar", f.AtBar.Value)
		assert.Equal(t, "my object value", f.Value)
	})

	t.Run("should find annotation AtFoo", func(t *testing.T) {
		as := annotation.Find(f, AtFoo{})
		assert.NotEqual(t, 0, len(as))
	})

	t.Run("should find annotation AtBaz", func(t *testing.T) {
		atBazFileds := annotation.Find(f, AtBaz{})
		assert.Equal(t, 1, len(atBazFileds))
	})

	t.Run("should notify bad syntax for struct tag pair", func(t *testing.T) {
		// notify bad syntax for struct tag pair
		b := new(bar)
		err := annotation.InjectIntoFields(b)
		assert.NotEqual(t, nil, err)
		assert.Equal(t, "bad syntax for struct tag pair", err.Error())
	})

	t.Run("should inject to object", func(t *testing.T) {
		fo := foo{}
		err := annotation.InjectIntoFields(fo)
		assert.NotEqual(t, nil, err)
	})

	t.Run("should check if an object implements Annotation", func(t *testing.T) {
		ok := annotation.Contains(f, AtFoo{})
		assert.Equal(t, true, ok)
	})

	t.Run("should inject annotation into sub struct", func(t *testing.T) {
		var fb struct{at.GetMapping `value:"/path/to/api"`}
		err := annotation.InjectIntoFields(&fb)
		assert.Equal(t, nil, err)
	})

	t.Run("should report error when inject nil object", func(t *testing.T) {
		err := annotation.InjectIntoFields(nil)
		assert.NotEqual(t, nil, err)
	})

	t.Run("should report error when inject invalid tagged annotation", func(t *testing.T) {
		ff, ok := annotation.GetField(&bar{}, AtFoo{})
		assert.Equal(t, true, ok)
		err := annotation.InjectIntoField(ff)
		assert.NotEqual(t, nil, err)
	})

	t.Run("should get annotation by type", func(t *testing.T) {
		var fb struct{at.PostMapping `value:"/path/to/api"`}
		f, ok := annotation.GetField(reflect.TypeOf(&fb), at.PostMapping{})
		assert.Equal(t, true, ok)
		assert.Equal(t, "PostMapping", f.StructField.Name)
		assert.Equal(t, false, f.Value.IsValid())
	})

	t.Run("should inject into a field", func(t *testing.T) {
		ff := new(foo)
		field, ok := annotation.GetField(ff, AtBaz{})
		assert.Equal(t, true, ok)
		err := annotation.InjectIntoField(field)
		assert.Equal(t, nil, err)
	})

	t.Run("should get and inject into multiple sub annotations", func(t *testing.T) {
		ma := &multipleBar{}
		maf := annotation.GetFields(ma)
		assert.Equal(t, 3, len(maf))

		err := annotation.InjectIntoFields(ma)
		assert.Equal(t, nil, err)
	})
}

