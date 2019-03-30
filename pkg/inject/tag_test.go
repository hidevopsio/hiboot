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

package inject

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestTag(t *testing.T) {
	tag := new(BaseTag)

	t.Run("should get properties", func(t *testing.T) {
		p := tag.Properties()
		assert.NotEqual(t, nil, p)
	})

	t.Run("should get check if it's singleton", func(t *testing.T) {
		s := tag.IsSingleton()
		assert.Equal(t, false, s)
	})

	t.Run("should get properties", func(t *testing.T) {
		fakeObj := struct{ Name string }{}
		objVal := reflect.ValueOf(fakeObj)
		field := objVal.Type().Field(0)
		f := tag.Decode(objVal, field, "", "fake")
		assert.Equal(t, nil, f)
	})
}
