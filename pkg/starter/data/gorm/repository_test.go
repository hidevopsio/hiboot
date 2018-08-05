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

package gorm

import (
	"testing"
	"github.com/hidevopsio/hiboot/pkg/starter/data/gorm/adapter"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/starter/data"
	"os"
)

func TestRepository(t *testing.T) {
	gorm := &properties{
		Type:      "mysql",
		Host:      "mysql-dev",
		Port:      "3306",
		Username:  os.Getenv("MYSQL_USERNAME"),
		Password:  os.Getenv("MYSQL_PASSWORD"),
		Database:  "test",
		ParseTime: "True",
		Charset:   "utf8",
		Loc:       "Asia%2FShanghai",
	}
	dataSource := new(dataSource)
	dataSource.gorm = new(adapter.FakeDataSource)
	repo := new(repository)

	t.Run("should report error if trying to close the unopened database", func(t *testing.T) {
		err := repo.CloseDataSource()
		assert.Equal(t, data.InvalidDataSourceError, err)
	})

	dataSource.Open(gorm)
	repo.SetDataSource(dataSource)

	t.Run("should close database properly", func(t *testing.T) {
		err := repo.CloseDataSource()
		assert.Equal(t, nil, err)
	})
}