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

package instantiate

import (
	"testing"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"github.com/stretchr/testify/assert"
)

func TestInstantiateFactory(t *testing.T) {
	type foo struct {Name string}
	f := new(foo)

	factory := new(InstantiateFactory)

	t.Run("should failed to set/get instance when factory is not initialized", func(t *testing.T) {
		inst := factory.GetInstance("not-exist-instance")
		assert.Equal(t, nil, inst)

		err := factory.SetInstance("foo", nil)
		assert.Equal(t, NotInitializedError, err)

		item := factory.Items()
		assert.Equal(t, 0, len(item))
	})

	ic := cmap.New()
	t.Run("should initialize factory", func(t *testing.T) {
		factory.Initialize(ic)
		assert.Equal(t, true, factory.Initialized())
	})

	t.Run("should set and get instance from factory", func(t *testing.T) {
		factory.SetInstance("foo", f)
		inst := factory.GetInstance("foo")
		assert.Equal(t, f, inst)
	})

	t.Run("should failed to get instance that does not exist", func(t *testing.T) {
		inst := factory.GetInstance("not-exist-instance")
		assert.Equal(t, nil, inst)
	})

	t.Run("should failed to set instance that it already exists in test mode", func(t *testing.T) {
		nf := new(foo)
		err := factory.SetInstance("foo", nf)
		assert.NotEqual(t, nil, err)
	})

	t.Run("should get factory items", func(t *testing.T) {
		items := factory.Items()
		assert.Equal(t, 1, len(items))
	})
}