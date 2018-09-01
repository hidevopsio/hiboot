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
	"github.com/hidevopsio/gorm"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

type User struct {
	ID       uint `gorm:"primary_key"`
	Username string
	Password string
	Age      int `gorm:"default:18"`
	Gender   int `gorm:"default:18"`
}

func (User) TableName() string {
	return "user"
}

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestDataSourceOpen(t *testing.T) {
	prop := &properties{
		Type:      "mysql",
		Host:      "mysql-dev",
		Port:      "3306",
		Username:  "test",
		Password:  "LcNxqoI4zZjAnpiTD7JQxLJR/IgL2iTiSZ2nd7KPEBgxMV+FVhPSzM+fgH93XqZJNpboN4F/buX22yLTXK38AcVGTfID3rmQAOAc9A2DIWNy5v9+3NOY00M8z4dR1XHojheK0681cY9QVjtlJ70jFFDXb7PjFc2fQ0GIyIjBQDY=",
		Database:  "test",
		ParseTime: true,
		Charset:   "utf8",
		Loc:       "Asia/Shanghai",
		Config: Config{
			Decrypt: true,
		},
	}
	dataSource := new(dataSource)

	t.Run("should report error when close database if it's not opened", func(t *testing.T) {
		err := dataSource.Close()
		assert.Equal(t, DatabaseIsNotOpenedError, err)
	})

	err := dataSource.Open(prop)
	t.Run("should open database", func(t *testing.T) {
		assert.NotEqual(t, nil, err)
	})

	t.Run("should close database", func(t *testing.T) {
		if dataSource.IsOpened() {
			err := dataSource.Close()
			assert.Equal(t, nil, err)
		}
	})

	t.Run("should Init dataSource", func(t *testing.T) {
		repo := new(gorm.FakeRepository)
		dataSource.Init(repo)

		err := dataSource.Close()
		assert.Equal(t, nil, err)
	})
}
