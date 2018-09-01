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
	"github.com/stretchr/testify/assert"
	"testing"
)

type Foo struct {
	Name string
}

func TestDecode(t *testing.T) {
	var foo Foo

	src := map[string]string{
		"name": "foo",
	}

	t.Run("should decode map to struct", func(t *testing.T) {
		err := Decode(&foo, src)
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
}
