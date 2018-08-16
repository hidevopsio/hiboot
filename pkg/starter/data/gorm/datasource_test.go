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
	"github.com/stretchr/testify/assert"
	"os"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/starter/data/gorm/adapter"
)

type User struct {
	ID   uint `gorm:"primary_key"`
	Username string
	Password string
	Age  int `gorm:"default:18"`
	Gender  int `gorm:"default:18"`
}

func (User) TableName() string {
	return "user"
}

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestDataSourceOpen(t *testing.T) {
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

	t.Run("should report error when close database if it's not opened", func(t *testing.T) {
		err := dataSource.Close()
		assert.Equal(t, DatabaseIsNotOpenedError, err)
	})

	t.Run("should open database", func(t *testing.T) {
		err := dataSource.Open(gorm)
		assert.Equal(t, nil, err)
	})

	t.Run("should close database", func(t *testing.T) {
		if dataSource.IsOpened() {
			err := dataSource.Close()
			assert.Equal(t, nil, err)
		}
	})

}
