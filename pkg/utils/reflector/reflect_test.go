package reflector

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
	"github.com/hidevopsio/hiboot/pkg/log"
)

type Foo struct{
	Name string
	Age int
	nickname string
}

func (f *Foo) Init(name string)  {
	f.Name = name
}

type Bar struct{
	Name string
	Age int
}

type Baz struct {
	Foo
	Bar Bar
}

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestConstructor(t *testing.T) {
	objVal := reflect.ValueOf(&Foo{})

	objType := objVal.Type()
	log.Debug("type: ", objType)
	objName := objType.Elem().Name()
	log.Debug("name: ", objName)

	object := objVal.Interface()
	log.Debug("object: ", object)

	// call Init
	method, ok := objType.MethodByName("Init")
	if ok {
		methodType := method.Type
		numIn := methodType.NumIn()
		inputs := make([]reflect.Value, numIn)
		for i := 0; i < numIn; i++ {
			t := methodType.In(i)


			log.Debugf("%v: %v %v", i, t.Name(), t)
		}
		inputs[0] = reflect.ValueOf(object)
		//method.Func.Call(inputs)
	}
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
	assert.Equal(t, 4, len(df))
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

func TestSetFieldValue(t *testing.T) {
	foo := &Foo{}
	value := "foo"
	t.Run("should set field value", func(t *testing.T) {
		// set field object
		err := SetFieldValue(foo, "Name", value)
		assert.Equal(t, nil, err)
		assert.Equal(t, "foo", value)
	})

	t.Run("should not set invalid object", func(t *testing.T) {
		x := 123
		err := SetFieldValue(x, "Name", 321)
		assert.Equal(t, InvalidInputError, err)
	})

	t.Run("should not set invalid object", func(t *testing.T) {
		err := SetFieldValue((*Foo)(nil), "Name", value)
		assert.Equal(t, InvalidInputError, err)
	})

	t.Run("should not set invalid object", func(t *testing.T) {
		err := SetFieldValue(foo, "nickname", value)
		assert.Equal(t, FieldCanNotBeSetError, err)
	})
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

	t.Run("should get struct name by pointer", func(t *testing.T) {
		n, err := GetName(new(Foo))
		assert.Equal(t, nil, err)
		assert.Equal(t, "Foo", n)
	})

	t.Run("should get struct name", func(t *testing.T) {
		n, err := GetName(Foo{})
		assert.Equal(t, nil, err)
		assert.Equal(t, "Foo", n)
	})

	t.Run("should get struct name by lower case", func(t *testing.T) {
		n, err := GetLowerCaseObjectName(new(Foo))
		assert.Equal(t, nil, err)
		assert.Equal(t, "foo", n)
	})

	t.Run("should return InvalidInputError if input is nil", func(t *testing.T) {
		_, err := GetLowerCaseObjectName((*Foo)(nil))
		assert.Equal(t, InvalidInputError, err)
	})
}