// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package reflector

import (
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type Foo struct {
	Name     string
	Age      int
	nickname string
}

func (f *Foo) Init(name string) {
	f.Name = name
}

func (f *Foo) Bar() {

}

func (f *Foo) Nickname() string {
	return f.nickname
}

func (f *Foo) SetNickname(nickname string) *Foo {
	f.nickname = nickname
	return f
}

type Bar struct {
	Name string
	Age  int
}

type Baz struct {
	Foo
	Bar Bar
}

type FooBar struct{}

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestLookupField(t *testing.T) {
	t.Run("should find anonymous field on Baz", func(t *testing.T) {
		assert.Equal(t, true, HasField(&Baz{}, "Foo"))
	})

	t.Run("should not find anonymous field on Foo", func(t *testing.T) {
		assert.Equal(t, false, HasField(&Bar{}, "Foo"))
	})
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
		k := GetKindByValue(reflect.ValueOf(x))
		assert.Equal(t, reflect.Uint, k)
	})

	t.Run("should return reflect.Int for int64", func(t *testing.T) {
		var x int64
		x = 1234
		k := GetKindByValue(reflect.ValueOf(x))
		assert.Equal(t, reflect.Int, k)
	})

	t.Run("should return reflect.Float32 for float64", func(t *testing.T) {
		var x float64
		x = 1.234
		k := GetKindByValue(reflect.ValueOf(x))
		assert.Equal(t, reflect.Float32, k)
	})

	t.Run("should return Ptr", func(t *testing.T) {
		k := GetKindByValue(reflect.ValueOf((*Foo)(nil)))
		assert.Equal(t, reflect.Ptr, k)
	})

	t.Run("should return Int", func(t *testing.T) {
		k := GetKindByType(reflect.ValueOf(int(1)).Type())
		assert.Equal(t, reflect.Int, k)
	})

	t.Run("should return Uint", func(t *testing.T) {
		k := GetKindByType(reflect.ValueOf(uint(1)).Type())
		assert.Equal(t, reflect.Uint, k)
	})

	t.Run("should return Bool", func(t *testing.T) {
		k := GetKindByType(reflect.ValueOf(true).Type())
		assert.Equal(t, reflect.Bool, k)
	})

	t.Run("should return Float32", func(t *testing.T) {
		k := GetKindByType(reflect.ValueOf(0.01).Type())
		assert.Equal(t, reflect.Float32, k)
	})

	t.Run("should return String", func(t *testing.T) {
		k := GetKindByType(reflect.ValueOf("abc").Type())
		assert.Equal(t, reflect.String, k)
	})
}

func TestValidateReflectType(t *testing.T) {
	baz := &Baz{Foo: Foo{Name: "foo"}, Bar: Bar{Name: "bar"}}

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
		s := []int{1, 2, 3}
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

func TestCallMethodByName(t *testing.T) {
	foo := new(Foo)
	t.Run("should call method SetNickname", func(t *testing.T) {
		res, err := CallMethodByName(foo, "SetNickname", "foobar")
		assert.Equal(t, nil, err)
		assert.Equal(t, foo, res)
		assert.Equal(t, "foobar", foo.nickname)
	})

	t.Run("should call method Nickname", func(t *testing.T) {
		res, err := CallMethodByName(foo, "Nickname")
		assert.Equal(t, nil, err)
		assert.Equal(t, "foobar", res)
	})

	t.Run("should call method Bar", func(t *testing.T) {
		res, err := CallMethodByName(foo, "Bar")
		assert.Equal(t, nil, err)
		assert.Equal(t, nil, res)
	})

	t.Run("should call method by name", func(t *testing.T) {
		_, err := CallMethodByName(foo, "NotExistMethod")
		assert.Equal(t, InvalidMethodError, err)
	})

}

type Embedded struct{}
type PtrEmbedded struct{}

type hidden struct{}

type Test struct {
	Embedded
	hidden
	*PtrEmbedded
	Exportedvalue int
	unexported    string
	A             struct{ A int }
}

func TestEmbedded(t *testing.T) {
	type fooInterface interface {

	}

	testCases := []struct {
		name      string
		source    interface{}
		fieldName string
		ok        bool
		private   bool
		anonymous bool
	}{
		{
			name:      "Should find Embedded field by name",
			source:    Test{},
			fieldName: "Embedded",
			ok:        true,
			private:   false,
			anonymous: true,
		},
		{
			name:      "Should find Embedded field by name",
			source:    &Test{},
			fieldName: "Embedded",
			ok:        true,
			private:   false,
			anonymous: true,
		},
		{
			name:      "Should find hidden field by name",
			source:    Test{},
			fieldName: "hidden",
			ok:        true,
			private:   true,
			anonymous: true,
		},
		{
			name:      "Should find PtrEmbedded field by name",
			source:    Test{},
			fieldName: "PtrEmbedded",
			ok:        true,
			private:   false,
			anonymous: true,
		},
		{
			name:      "Should find Exportedvalue field by name",
			source:    Test{},
			fieldName: "Exportedvalue",
			ok:        true,
			private:   false,
			anonymous: false,
		},
		{
			name:      "Should find unexported field by name",
			source:    Test{},
			fieldName: "unexported",
			ok:        true,
			private:   true,
			anonymous: false,
		},
		{
			name:      "Should find A field by name",
			source:    Test{},
			fieldName: "A",
			ok:        true,
			private:   false,
			anonymous: false,
		},
	}

	for _, testCase := range testCases {
		typ := IndirectType(reflect.TypeOf(testCase.source))
		t.Run(testCase.name, func(t *testing.T) {
			field, ok := typ.FieldByName(testCase.fieldName)
			exported := field.PkgPath != ""
			assert.Equal(t, testCase.ok, ok)
			assert.Equal(t, testCase.private, exported)
			assert.Equal(t, testCase.anonymous, field.Anonymous)
			assert.Equal(t, testCase.fieldName, field.Name)
		})

		t.Run(testCase.name, func(t *testing.T) {
			ok := HasEmbeddedField(testCase.source, testCase.fieldName)
			assert.Equal(t, testCase.anonymous, ok)
		})
	}
}

func TestParseObjectName(t *testing.T) {
	t.Run("should parse object name", func(t *testing.T) {
		name := ParseObjectName(new(FooBar), "Bar")
		assert.Equal(t, "foo", name)
	})
}

func TestGetPkgPath(t *testing.T) {
	t.Run("should get object pkg path", func(t *testing.T) {
		pkgPath := GetPkgPath(Foo{})
		assert.Contains(t, "github.com/hidevopsio/hiboot/pkg/utils/reflector", pkgPath)
	})
}

func TestParseObjectPkgName(t *testing.T) {
	pkgName := ParseObjectPkgName(Foo{})
	assert.Equal(t, "reflector", pkgName)

	pkgName = ParseObjectPkgName(&Foo{})
	assert.Equal(t, "reflector", pkgName)
}

func SayHello(name string) string {
	return "Hello " + name
}

func Dummy() {
	// for test only, do nothing
}

func TestCallFunc(t *testing.T) {
	t.Run("should call func", func(t *testing.T) {
		res, err := CallFunc(SayHello, "Steve")
		assert.Equal(t, nil, err)
		assert.NotEqual(t, nil, res)
		assert.Equal(t, "Hello Steve", res.(string))
	})

	t.Run("should call func", func(t *testing.T) {
		res, err := CallFunc(Dummy)
		assert.Equal(t, nil, err)
		assert.Equal(t, nil, res)
	})

	t.Run("should call func", func(t *testing.T) {
		_, err := CallFunc(int(1))
		assert.Equal(t, InvalidFuncError, err)
	})
}

func TestGetEmbeddedInterfaceField(t *testing.T) {
	type fakeInterface interface{}
	type fakeService struct{ fakeInterface }
	type fakeChildService struct{ fakeService }

	t.Run("should get embedded fakeInterface", func(t *testing.T) {
		field := GetEmbeddedInterfaceField(new(fakeService))
		assert.Equal(t, "fakeInterface", field.Name)
	})

	t.Run("should get embedded fakeInterface", func(t *testing.T) {
		field := GetEmbeddedInterfaceField(new(fakeChildService))
		assert.Equal(t, "fakeInterface", field.Name)
	})

	t.Run("should not get embedded interface that it's not exist", func(t *testing.T) {
		type fooService struct{}
		type barService struct{ fooService }
		field := GetEmbeddedInterfaceField(new(barService))
		assert.Equal(t, false, field.Anonymous)
	})
}
