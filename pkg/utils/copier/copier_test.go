package copier

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"errors"
)

type Foo struct {
	Name string
}

type Bar struct {
	Foo
}

func TestCopier(t *testing.T) {

	testCases := []struct {
		name string
		from interface{}
		to interface{}
		err error
	} {
		{
			name: "1. Should copy from Foo to Bar",
			from: &Foo{ Name: "foo"},
			to: &Bar{},
			err: nil,
		},
		{
			name: "2. copy to value is unaddressable",
			from: &Foo{ Name: "foo"},
			to: nil,
			err: errors.New("copy to value is unaddressable"),
		},
		{
			name: "3. copy source is invalid",
			from: nil,
			to: &Bar{},
			err: errors.New("copy source is invalid"),
		},
		{
			name: "4. should copy directly",
			from: &Foo{Name: "foo"},
			to: &Foo{},
			err: nil,
		},
		{
			name: "5. source or target type is not struct",
			from: []string{"a", "b"},
			to: &Foo{Name: "foo"},
			err: errors.New("source or target type is not struct"),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := Copy(testCase.to, testCase.from)
			assert.Equal(t, testCase.err, err)
		})

	}

}
