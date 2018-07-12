package reflector

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
)

type Foo struct{
	Name string
}

type Bar struct{
	Name string
}

type Baz struct {
	Foo
	Bar Bar
}

func TestNewReflectType(t *testing.T) {
	foo := NewReflectType(Foo{})
	assert.NotEqual(t, nil, foo)
}

func TestCopy(t *testing.T) {
	foo := &Foo{Name: "foo"}
	bar := &Bar{}
	f := reflect.ValueOf(foo)
	b := reflect.ValueOf(bar)
	ok := Copy(b, f)
	assert.Equal(t, true, ok)
	assert.Equal(t, foo.Name, bar.Name)
}

func TestValidate(t *testing.T) {
	foo := &Foo{Name: "foo"}
	f, err := Validate(foo)
	assert.Equal(t, nil, err)
	assert.Equal(t, reflect.Struct, f.Kind())
}

func TestDeepFields(t *testing.T) {
	baz := &Baz{Bar: Bar{Name: "bar"}}
	baz.Name = "foo"
	bt := reflect.TypeOf(baz)
	df := DeepFields(bt)
	assert.Equal(t, 2, len(df))
	assert.Equal(t, "Name", df[0].Name)
	assert.Equal(t, "Bar", df[1].Name)
}

func TestIndirect(t *testing.T) {
	foo := &Foo{Name: "foo"}
	f := reflect.ValueOf(foo)
	fv := Indirect(f)
	assert.Equal(t, reflect.Ptr, f.Kind())
	assert.Equal(t, reflect.Struct, fv.Kind())
}

func TestIndirectType(t *testing.T) {
	foo := &Foo{Name: "foo"}
	f := reflect.TypeOf(foo)
	ft := IndirectType(f)
	assert.Equal(t, reflect.Ptr, f.Kind())
	assert.Equal(t, reflect.Struct, ft.Kind())
}

func TestGetFieldValue(t *testing.T) {
	foo := &Foo{Name: "foo"}
	fv := GetFieldValue(foo, "Name")
	assert.Equal(t, "foo", fv.Interface())
}

func TestParseReferences(t *testing.T) {
	foo := &Foo{Name: "foo"}
	rv, err := ParseReferences(foo, []string{"name"})
	assert.Equal(t, nil, err)
	assert.Equal(t, "foo", rv)
}

func TestGetKind(t *testing.T) {
	var x uint64
	x = 1234

	k := GetKind(reflect.ValueOf(x))
	assert.Equal(t, reflect.Uint, k)
}

func TestValidateReflectType(t *testing.T) {
	baz := &Baz{Foo: Foo{Name:"foo"}, Bar: Bar{Name: "bar"}}
	ValidateReflectType(baz, func(value *reflect.Value, reflectType reflect.Type, fieldSize int, isSlice bool) error {
		assert.Equal(t, *baz, value.Interface())
		assert.Equal(t, 1, fieldSize)
		assert.Equal(t, "Baz", reflectType.Name())
		assert.Equal(t, false, isSlice)
		return nil
	})
}