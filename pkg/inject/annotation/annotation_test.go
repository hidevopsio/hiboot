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
	AtBar `value:"bar0"`
	Bar1 struct{
		AtBar `value:"bar1"`
		AtBaz `value:"baz"`
	}
	Bar2 struct{
		AtBar `value:"bar2"`
		AtBaz `value:"baz"`
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
	at.Operation `value:"testApi" operationId:"getGreeting" description:"This is the Greeting api for demo"`
}

func TestImplementsAnnotation(t *testing.T) {
	log.SetLevel("debug")

	f := new(foo)
	f.Value = "my object value"
	t.Run("test api operation", func(t *testing.T) {
		ao := &atApiOperation{}
		err := annotation.InjectAll(ao)
		assert.Equal(t, nil, err)
	})

	t.Run("test map injection", func(t *testing.T) {
		fb := &foobar{}
		err := annotation.InjectAll(fb)
		assert.Equal(t, nil, err)
	})

	annotations := annotation.GetAnnotations(f)
	t.Run("should check if object contains at.Annotation", func(t *testing.T) {
		assert.NotEqual(t, nil, annotations)

		ok := annotation.Contains(annotations, AtFoo{})
		assert.Equal(t, true, ok)
	})

	t.Run("should find if contains child annotation", func(t *testing.T) {
		ok := annotation.ContainsChild(annotations, AtFoo{})
		assert.Equal(t, true, ok)
	})

	t.Run("should find child annotation", func(t *testing.T) {
		a := annotation.Find(annotations, AtFoo{})
		assert.Equal(t, "AtFoo", a.Field.StructField.Name)
	})

	t.Run("should filter in child annotation", func(t *testing.T) {
		mb := &multipleBar{}
		ma := annotation.GetAnnotations(mb)
		abs := annotation.FilterIn(ma, AtBar{})
		assert.Equal(t, 3, len(abs))
		assert.Equal(t, "AtBar", abs[0].Field.StructField.Name)
		assert.Equal(t, "AtBar", abs[1].Field.StructField.Name)
		assert.Equal(t, "AtBar", abs[2].Field.StructField.Name)
	})

	t.Run("should check if object contains at.Annotation", func(t *testing.T) {
		ok := annotation.Contains(&f.AtFoo, at.Annotation{})
		assert.Equal(t, true, ok)
	})

	t.Run("should get annotation AtFoo", func(t *testing.T) {
		af := annotation.GetAnnotation(f, AtFoo{})
		value, ok := af.Field.StructField.Tag.Lookup("value")
		assert.Equal(t, "foo,option 1,option 2", value)
		assert.Equal(t, true, ok)

		age, ok := af.Field.StructField.Tag.Lookup("age")
		assert.Equal(t, "18", age)
		assert.Equal(t, true, ok)
	})

	t.Run("should report error for invalid type that pass to GetFields", func(t *testing.T) {
		af := annotation.GetAnnotations(123)
		assert.True(t, af == nil)
	})

	t.Run("should report error for invalid type that pass to GetField", func(t *testing.T) {
		a := annotation.GetAnnotation(123, AtFoo{})
		assert.True(t, nil == a)
	})

	t.Run("should inject all annotations", func(t *testing.T) {
		assert.Equal(t, "", f.AtFoo.Value)
		assert.Equal(t, 0, f.AtFoo.Age)
		err := annotation.InjectAll(f)
		assert.Equal(t, nil, err)
		assert.Equal(t, "foo", f.AtFoo.Value)
		assert.Equal(t, 18, f.AtFoo.Age)

		assert.Equal(t, "bar", f.AtBar.Value)
		assert.Equal(t, "my object value", f.Value)
	})

	t.Run("should find annotation AtFoo", func(t *testing.T) {
		as := annotation.FindAll(f, AtFoo{})
		assert.NotEqual(t, 0, len(as))
	})

	t.Run("should find annotation AtBaz", func(t *testing.T) {
		atBazFileds := annotation.FindAll(f, AtBaz{})
		assert.Equal(t, 1, len(atBazFileds))
	})

	t.Run("should notify bad syntax for struct tag pair", func(t *testing.T) {
		// notify bad syntax for struct tag pair
		b := new(bar)
		err := annotation.InjectAll(b)
		assert.NotEqual(t, nil, err)
		assert.Equal(t, "bad syntax for struct tag pair", err.Error())
	})

	t.Run("should inject to object", func(t *testing.T) {
		fo := foo{}
		err := annotation.InjectAll(fo)
		assert.NotEqual(t, nil, err)
	})

	t.Run("should check if an object implements Annotation", func(t *testing.T) {
		ok := annotation.Contains(f, AtFoo{})
		assert.Equal(t, true, ok)
	})

	t.Run("should inject annotation into sub struct", func(t *testing.T) {
		var fb struct{at.GetMapping `value:"/path/to/api"`}
		err := annotation.InjectAll(&fb)
		assert.Equal(t, nil, err)
	})

	t.Run("should report error when inject nil object", func(t *testing.T) {
		err := annotation.InjectAll(nil)
		assert.NotEqual(t, nil, err)
	})

	t.Run("should report error when inject invalid tagged annotation", func(t *testing.T) {
		ff := annotation.GetAnnotation(&bar{}, AtFoo{})
		err := annotation.Inject(ff)
		assert.NotEqual(t, nil, err)
	})

	t.Run("should get annotation by type", func(t *testing.T) {
		var fb struct{at.PostMapping `value:"/path/to/api"`}
		f := annotation.GetAnnotation(reflect.TypeOf(&fb), at.PostMapping{})
		assert.NotEqual(t, nil, f)
		assert.Equal(t, "PostMapping", f.Field.StructField.Name)
		assert.Equal(t, false, f.Field.Value.IsValid())
	})

	t.Run("should inject into a field", func(t *testing.T) {
		ff := new(foo)
		ann := annotation.GetAnnotation(ff, AtBaz{})
		assert.NotEqual(t, nil, ann)
		err := annotation.Inject(ann)
		assert.Equal(t, nil, err)
	})

	t.Run("should get and inject into multiple sub annotations", func(t *testing.T) {
		ma := &multipleBar{}
		maf := annotation.GetAnnotations(ma)
		assert.Equal(t, 2, len(maf.Items))
		assert.Equal(t, 2, len(maf.Children))

		err := annotation.InjectAll(ma)
		assert.Equal(t, nil, err)
		assert.Equal(t, "bar1", ma.Bar1.AtBar.Value)
		assert.Equal(t, "bar2", ma.Bar2.AtBar.Value)
		assert.Equal(t, "baz", ma.Bar1.AtBaz.Value)

		fa := annotation.Find(maf, AtFoo{})
		assert.Equal(t, "foo", fa.Field.Value.Interface().(AtFoo).Value)
	})

	t.Run("should get from nil interface", func(t *testing.T) {
		a := annotation.GetAnnotation((*multipleBar)(nil), AtFoo{})
		assert.Equal(t, false, a.Field.Value.IsValid())
	})

	t.Run("should get from nil interface", func(t *testing.T) {
		a := annotation.GetAnnotation((*annotation.Annotations)(nil), AtFoo{})
		assert.True(t, nil == a)
	})

	t.Run("should get nil from nil", func(t *testing.T) {
		a := annotation.GetAnnotation(nil, AtFoo{})
		assert.True(t, nil == a)
	})
	t.Run("should get nil from nil", func(t *testing.T) {
		a := annotation.GetAnnotations(nil)
		assert.True(t, nil == a)
	})

	t.Run("should get annotation from function return object", func(t *testing.T) {
		a := annotation.GetAnnotation(func() *multipleBar {
			return &multipleBar{}
		}, AtFoo{})
		assert.Equal(t, "AtFoo", a.Field.StructField.Name)
	})
}

