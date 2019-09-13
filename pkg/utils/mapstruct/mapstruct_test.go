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

package mapstruct

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Foo struct {
	Name string `json:"name"`
}

type Bar struct {
	Name string `json:"name"`
	Foo  Foo    `json:"foo"`
	Foos map[string]Foo
}

type Basic struct {
	Foo
	Vstring     string      `json:"vstring"`
	Vint        int         `json:"vint"`
	Vuint       uint        `json:"vuint"`
	Vbool       bool        `json:"vbool"`
	Vfloat      float64     `json:"vfloat"`
	Vextra      string      `json:"vextra"`
	Vsilent     bool        `json:"vsilent"`
	Vdata       interface{} `json:"vdata"`
	VjsonInt    int         `json:"vjson_int"`
	VjsonFloat  float64     `json:"vjson_float"`
	VjsonNumber json.Number `json:"vjson_number"`
	Bar         Bar         `json:"bar"`
}

type EmbeddedSquash struct {
	Basic
	Vunique string `json:"vunique"`
}

func TestDecode(t *testing.T) {
	var foo Foo

	src := map[string]string{
		"name": "foo",
	}

	t.Run("should decode map to struct", func(t *testing.T) {
		type Foo struct {
			Name string `json:"-" at:"name"`
		}

		foo := &Foo{}


		err := Decode(foo, src)
		assert.Equal(t, nil, err)
		assert.Equal(t, "", foo.Name)

		err = Decode(foo, src, WithAnnotation)
		assert.Equal(t, nil, err)
		assert.Equal(t, "foo", foo.Name)
	})

	t.Run("should decode map to struct", func(t *testing.T) {
		err := Decode(&foo, src, WithSquash, WithWeaklyTypedInput)
		assert.Equal(t, nil, err)
		assert.Equal(t, "foo", foo.Name)
	})

	t.Run("should return result must be a pointer", func(t *testing.T) {
		err := Decode(nil, src)
		assert.Equal(t, "result must be a pointer", err.Error())
	})

	t.Run("should return parameters of mapstruct.Decode must not be nill", func(t *testing.T) {
		err := Decode(&foo, nil)
		assert.Equal(t, "parameters of mapstruct.Decode must not be nil", err.Error())
	})

	t.Run("should convert struct to map", func(t *testing.T) {

		data := &EmbeddedSquash{
			Vunique: "test",
			Basic: Basic{
				Foo:         Foo{Name: "embedded child field"},
				Vstring:     "test vstr",
				Vint:        123,
				Vuint:       456,
				Vbool:       true,
				Vfloat:      0.123,
				Vextra:      "test extra str",
				Vsilent:     false,
				Vdata:       &Foo{Name: "test child field"},
				VjsonInt:    33,
				VjsonFloat:  3.14,
				VjsonNumber: "12345",
				Bar: Bar{
					Name: "bar name",
					Foo:  Foo{Name: "foo name"},
					Foos: map[string]Foo{
						"/a": {Name: "f1"},
						"/b": {Name: "f2"},
						"/c": {Name: "f3"},
					},
				},
			},
		}
		m, ok := DecodeStructToMap(data)
		assert.Equal(t, true, ok)
		assert.Equal(t, "embedded child field", m["name"])
	})
}
