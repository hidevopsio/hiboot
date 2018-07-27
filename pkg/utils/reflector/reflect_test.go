package reflector

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
)

type Foo struct{
	Name string
	Age int
}

type Bar struct{
	Name string
	Age int
}

type Baz struct {
	Foo
	Bar Bar
}

func TestNewReflectType(t *testing.T) {
	foo := NewReflectType(Foo{})
	assert.NotEqual(t, nil, foo)
}

func TestValidate(t *testing.T) {
	foo := &Foo{Name: "foo"}

	t.Run("should validate object foo", func(t *testing.T) {
		f, err := Validate(foo)
		assert.Equal(t, nil, err)
		assert.Equal(t, reflect.Struct, f.Kind())
	})

	t.Run("should return error that value is unaddressable", func(t *testing.T) {
		_, err := Validate(123)
		assert.Equal(t, "value is unaddressable", err.Error())
	})

	t.Run("should return error that value is not valid", func(t *testing.T) {
		_, err := Validate((*Foo)(nil))
		assert.Equal(t, "value is not valid", err.Error())
	})
}

func TestDeepFields(t *testing.T) {
	baz := &Baz{Bar: Bar{Name: "bar"}}
	baz.Name = "foo"
	bt := reflect.TypeOf(baz)
	df := DeepFields(bt)
	assert.Equal(t, 3, len(df))
	assert.Equal(t, "Name", df[0].Name)
	assert.Equal(t, "Age", df[1].Name)
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

func TestGetKind(t *testing.T) {

	t.Run("should return reflect.Unit for uint64", func(t *testing.T) {
		var x uint64
		x = 1234
		k := GetKind(reflect.ValueOf(x))
		assert.Equal(t, reflect.Uint, k)
	})

	t.Run("should return reflect.Int for int64", func(t *testing.T) {
		var x int64
		x = 1234
		k := GetKind(reflect.ValueOf(x))
		assert.Equal(t, reflect.Int, k)
	})

	t.Run("should return reflect.Float32 for float64", func(t *testing.T) {
		var x float64
		x = 1.234
		k := GetKind(reflect.ValueOf(x))
		assert.Equal(t, reflect.Float32, k)
	})


	t.Run("should return Ptr", func(t *testing.T) {
		k := GetKind(reflect.ValueOf((*Foo)(nil)))
		assert.Equal(t, reflect.Ptr, k)
	})
}

func TestValidateReflectType(t *testing.T) {
	baz := &Baz{Foo: Foo{Name:"foo"}, Bar: Bar{Name: "bar"}}

	t.Run("should validete reflect type", func(t *testing.T) {
		ValidateReflectType(baz, func(value *reflect.Value, reflectType reflect.Type, fieldSize int, isSlice bool) error {
			assert.Equal(t, *baz, value.Interface())
			assert.Equal(t, 1, fieldSize)
			assert.Equal(t, "Baz", reflectType.Name())
			assert.Equal(t, false, isSlice)
			return nil
		})
	})

	t.Run("should return value is not valid", func(t *testing.T) {
		err := ValidateReflectType((*Foo)(nil), nil)
		assert.Equal(t, "value is not valid", err.Error())
	})

	t.Run("should validate slice", func(t *testing.T) {
		s := []int {1, 2, 3}
		err := ValidateReflectType(&s, nil)
		assert.Equal(t, nil, err)
	})

}

func TestGetName(t *testing.T) {
	n, err := GetName(new(Foo))
	assert.Equal(t, nil, err)
	assert.Equal(t, "Foo", n)

	n, err = GetName(Foo{})
	assert.Equal(t, nil, err)
	assert.Equal(t, "Foo", n)

	n, err = GetLowerCaseObjectName(new(Foo))
	assert.Equal(t, nil, err)
	assert.Equal(t, "foo", n)
}