package copier

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"errors"
)

type Foo struct {
	Name string
}

func (f *Foo) get(in string) string  {
	return f.Name + in
}

type Bar struct {
	Name string
}


type Baz struct {
	Bar
}

type NestedFoo struct {
	Name string
	Foo Foo
}


type NestedBar struct {
	Name string
	Foo Foo
}

type MyString string

func TestCopier(t *testing.T) {

	testCases := []struct {
		name string
		from interface{}
		to interface{}
		err error
	} {
		{
			name: "Should copy from Foo to Bar",
			from: &Foo{ Name: "foo"},
			to: &Bar{},
			err: nil,
		},
		{
			name: "Should copy from Bar to Baz",
			from: &Bar{ Name: "bar"},
			to: &Baz{},
			err: nil,
		},
		{
			name: "copy to value is unaddressable",
			from: &Foo{ Name: "foo"},
			to: nil,
			err: errors.New("copy to value is unaddressable"),
		},
		{
			name: "copy source is invalid",
			from: nil,
			to: &Bar{},
			err: errors.New("copy source is invalid"),
		},
		{
			name: "should copy directly",
			from: &Foo{Name: "foo"},
			to: &Foo{},
			err: nil,
		},
		{
			name: "source or target type is not struct",
			from: []string{"a", "b"},
			to: &Foo{Name: "foo"},
			err: errors.New("source or target type is not struct"),
		},
		{
			name: "should copy struct slice",
			from: &[]Foo{{Name: "foo"}, {Name: "bar"}},
			to: &[]Bar{},
			err: nil,
		},
		{
			name: "should copy struct slice",
			from: &Foo{Name: "bar"},
			to: &[]Bar{},
			err: nil,
		},
		{
			name: "should copy struct slice",
			from: &NestedFoo{Name: "nested foo", Foo: Foo{Name: "foo"}},
			to: &NestedFoo{},
			err: nil,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := Copy(testCase.to, testCase.from)
			assert.Equal(t, testCase.err, err)
		})

	}

}
