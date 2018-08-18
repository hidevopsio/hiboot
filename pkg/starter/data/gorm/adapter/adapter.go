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

package adapter

import (
	"github.com/hidevopsio/gorm"
	"errors"
)

type Gorm interface {
	Open(dialect string, args ...interface{}) (db gorm.Repository, err error)
	Close() error
}

type GormDataSource struct {
	db gorm.Repository
}

var DatabaseIsNotOpenedError = errors.New("database is not opened")

func (d *GormDataSource) Open(dialect string, args ...interface{}) (db gorm.Repository, err error) {
	d.db, err = gorm.Open(dialect, args...)
	if err != nil {
		d.db = nil
	}
	return d.db, err
}

func (d *GormDataSource) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return DatabaseIsNotOpenedError
}
