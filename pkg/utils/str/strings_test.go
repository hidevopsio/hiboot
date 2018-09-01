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

package str

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
)

func TestLowerFirst(t *testing.T) {
	s := "Foo"

	ns := LowerFirst(s)

	assert.Equal(t, "foo", ns)

	es := LowerFirst("")

	assert.Equal(t, "", es)
}

func TestUpperFirst(t *testing.T) {
	s := "foo"

	ns := UpperFirst(s)

	assert.Equal(t, "Foo", ns)
}

func TestStringInSlice(t *testing.T) {
	s := []string{
		"foo",
		"bar",
		"baz",
	}

	assert.Equal(t, true, InSlice("bar", s))
}

func TestConvert(t *testing.T) {
	testData := []struct {
		src  string
		kind reflect.Kind
		dst interface{}
	}{
		{
			src:  "a,b,c,d",
			kind: reflect.Slice,
			dst: []string{"a", "b", "c", "d"},
		},
		{
			src:  "this is a string",
			kind: reflect.String,
			dst: "this is a string",
		},
		{
			src:  "1234",
			kind: reflect.Int,
			dst: int(1234),
		},
		{
			src:  "12",
			kind: reflect.Int8,
			dst: int8(12),
		},
		{
			src:  "123",
			kind: reflect.Int16,
			dst: int16(123),
		},
		{
			src:  "2345",
			kind: reflect.Int32,
			dst: int32(2345),
		},
		{
			src:  "12345678",
			kind: reflect.Int64,
			dst: int64(12345678),
		},
		{
			src:  "4321",
			kind: reflect.Uint,
			dst: uint(4321),
		},
		{
			src:  "21",
			kind: reflect.Uint8,
			dst: uint8(21),
		},
		{
			src:  "321",
			kind: reflect.Uint16,
			dst: uint16(321),
		},
		{
			src:  "5432",
			kind: reflect.Uint32,
			dst: uint32(5432),
		},
		{
			src:  "87654321",
			kind: reflect.Uint64,
			dst: uint64(87654321),
		},
		{
			src:  "0.1",
			kind: reflect.Float32,
			dst: float32(0.1),
		},
		{
			src:  "0.1234",
			kind: reflect.Float64,
			dst: float64(0.1234),
		},
		{
			src:  "true",
			kind: reflect.Bool,
			dst: true,
		},
		{
			src:  "",
			kind: reflect.Int,
			dst: int(0),
		},
		{
			src:  "",
			kind: reflect.Int8,
			dst: int8(0),
		},
		{
			src:  "",
			kind: reflect.Int16,
			dst: int16(0),
		},
		{
			src:  "",
			kind: reflect.Int32,
			dst: int32(0),
		},
		{
			src:  "",
			kind: reflect.Int64,
			dst: int64(0),
		},
		{
			src:  "",
			kind: reflect.Uint,
			dst: uint(0),
		},
		{
			src:  "",
			kind: reflect.Uint8,
			dst: uint8(0),
		},
		{
			src:  "",
			kind: reflect.Uint16,
			dst: uint16(0),
		},
		{
			src:  "",
			kind: reflect.Uint32,
			dst: uint32(0),
		},
		{
			src:  "",
			kind: reflect.Uint64,
			dst: uint64(0),
		},
		{
			src:  "",
			kind: reflect.Float32,
			dst: float32(0.0),
		},
		{
			src:  "",
			kind: reflect.Float64,
			dst: float64(0.0),
		},
		{
			src:  "",
			kind: reflect.Bool,
			dst: false,
		},
		{
			src:  " ",
			kind: reflect.Int,
			dst: int(0),
		},
		{
			src:  " ",
			kind: reflect.Int8,
			dst: int8(0),
		},
		{
			src:  " ",
			kind: reflect.Int16,
			dst: int16(0),
		},
		{
			src:  " ",
			kind: reflect.Int32,
			dst: int32(0),
		},
		{
			src:  " ",
			kind: reflect.Int64,
			dst: int64(0),
		},
		{
			src:  " ",
			kind: reflect.Uint,
			dst: uint(0),
		},
		{
			src:  " ",
			kind: reflect.Uint8,
			dst: uint8(0),
		},
		{
			src:  "",
			kind: reflect.Uint16,
			dst: uint16(0),
		},
		{
			src:  " ",
			kind: reflect.Uint32,
			dst: uint32(0),
		},
		{
			src:  " ",
			kind: reflect.Uint64,
			dst: uint64(0),
		},
		{
			src:  " ",
			kind: reflect.Float32,
			dst: float32(0.0),
		},
		{
			src:  " ",
			kind: reflect.Float64,
			dst: float64(0.0),
		},
		{
			src:  " ",
			kind: reflect.Bool,
			dst: false,
		},
	}

	t.Run("should convert all test data to specific type", func(t *testing.T) {
		for _, data := range testData {
			dst := Convert(data.src, data.kind)
			assert.Equal(t, data.dst, dst)
		}
	})
}
