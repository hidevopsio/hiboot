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

package data

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
)

type Foo struct {
	ID string
	Name string
}

type Bar struct {
	Id string
	Name string
}

func TestParse(t *testing.T) {

	repo := BaseKVRepository{}
	foo := &Foo{ID: "test"}

	t.Run("should parse ID", func(t *testing.T) {
		bucket, key, value, err := repo.Parse(foo)
		assert.Equal(t, nil, err)
		assert.Equal(t, []byte("foo"), bucket)
		assert.Equal(t, []byte("test"), key)
		assert.Equal(t, foo, value)
	})

	bar := &Bar{Id: "test"}
	t.Run("should parse Id", func(t *testing.T) {
		bucket, key, value, err := repo.Parse(bar)
		assert.Equal(t, nil, err)
		assert.Equal(t, []byte("bar"), bucket)
		assert.Equal(t, []byte("test"), key)
		assert.Equal(t, bar, value)
	})

	baz := &Bar{}
	t.Run("should parse input Id", func(t *testing.T) {
		bucket, key, value, err := repo.Parse("test", baz)
		assert.Equal(t, nil, err)
		assert.Equal(t, []byte("bar"), bucket)
		assert.Equal(t, []byte("test"), key)
		assert.Equal(t, baz, value)
	})

	foobar := &Bar{}
	t.Run("should not parse Id", func(t *testing.T) {
		_, _, _, err := repo.Parse(foobar)
		assert.Equal(t, InvalidDataModelError, err)

	})

	t.Run("should not parse Id with nil input", func(t *testing.T) {
		_, _, _, err := repo.Parse("a", (*Foo)(nil))
		assert.Equal(t, reflector.InvalidInputError, err)
	})

	t.Run("should pass test on unimplemented method", func(t *testing.T) {
		err := repo.Put(foo)
		assert.Equal(t, nil, err)

		err = repo.Get(foo)
		assert.Equal(t, nil, err)

		err = repo.Delete(foo)
		assert.Equal(t, nil, err)
	})
}