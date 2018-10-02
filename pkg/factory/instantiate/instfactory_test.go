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

package instantiate_test

import (
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/factory/instantiate"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"github.com/stretchr/testify/assert"
	"testing"
)

type FooBar struct {
	Name string
}

type fooBarService struct {
	Name   string
	fooBar *FooBar
}

type BarService interface {
	Bar() string
}

type BarServiceImpl struct {
	BarService
}

func (s *BarServiceImpl) Bar() string {
	return "bar"
}

func newFooBarService(fooBar *FooBar) *fooBarService {
	return &fooBarService{
		fooBar: fooBar,
	}
}

func TestInstantiateFactory(t *testing.T) {
	type foo struct{ Name string }
	f := new(foo)

	appFactory := new(instantiate.InstantiateFactory)
	testName := "foobar"
	t.Run("should failed to set/get instance when factory is not initialized", func(t *testing.T) {
		inst := appFactory.GetInstance("not-exist-instance")
		assert.Equal(t, nil, inst)

		err := appFactory.SetInstance("foo", nil)
		assert.Equal(t, instantiate.ErrNotInitialized, err)

		item := appFactory.Items()
		assert.Equal(t, 0, len(item))
	})

	ic := cmap.New()
	t.Run("should initialize factory", func(t *testing.T) {
		appFactory.Initialize(ic, []*factory.MetaData{
			{Object: f},
			{Object: &FooBar{Name: testName}},
			{Object: &BarServiceImpl{}},
		})
		assert.Equal(t, true, appFactory.Initialized())
	})

	t.Run("should build components", func(t *testing.T) {
		appFactory.BuildComponents()
	})

	t.Run("should set and get instance from factory", func(t *testing.T) {
		appFactory.SetInstance("foo", f)
		inst := appFactory.GetInstance("foo")
		assert.Equal(t, f, inst)
	})

	t.Run("should failed to get instance that does not exist", func(t *testing.T) {
		inst := appFactory.GetInstance("not-exist-instance")
		assert.Equal(t, nil, inst)
	})

	t.Run("should failed to set instance that it already exists in test mode", func(t *testing.T) {
		nf := new(foo)
		err := appFactory.SetInstance("foo", nf)
		assert.NotEqual(t, nil, err)
	})

	t.Run("should get factory items", func(t *testing.T) {
		items := appFactory.Items()
		assert.Equal(t, 4, len(items))
	})

	t.Run("should check valid object", func(t *testing.T) {
		assert.Equal(t, true, appFactory.IsValidObjectType(f))
	})

	t.Run("should check invalid object", func(t *testing.T) {
		assert.Equal(t, false, appFactory.IsValidObjectType(1))
	})

	t.Run("should append component", func(t *testing.T) {
		appFactory.AppendComponent(&struct {
			Name string
		}{Name: "foo"})
		assert.Equal(t, false, appFactory.IsValidObjectType(1))
	})
}
