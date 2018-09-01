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

package bolt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDataSource(t *testing.T) {

	properties := &properties{
		Database: "test.db",
		Mode:     0600,
		Timeout:  2,
	}

	d := GetDataSource()

	t.Run("should open bolt database", func(t *testing.T) {
		err := d.Open(nil)
		assert.Equal(t, InvalidPropertiesError, err)
	})

	t.Run("should open bolt database", func(t *testing.T) {
		err := d.Open(properties)
		assert.Equal(t, nil, err)
	})

	// close bolt database
	d.Close()
}

func TestDataSourceWithEmptyFile(t *testing.T) {

	properties := &properties{
		Timeout: 2,
	}
	d := GetDataSource()
	t.Run("should open bolt database", func(t *testing.T) {
		err := d.Open(properties)
		assert.NotEqual(t, nil, err)
	})
}
