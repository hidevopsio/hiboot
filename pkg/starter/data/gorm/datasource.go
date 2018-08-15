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
	"fmt"
	"errors"
	"sync"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"github.com/hidevopsio/hiboot/pkg/starter/data/gorm/adapter"
	"github.com/hidevopsio/hiboot/pkg/utils/crypto/rsa"
)

type DataSource interface {
	Open(p *properties) error
	IsOpened() bool
	Close() error
	DB() adapter.DB
}

type dataSource struct {
	gorm adapter.Gorm
	db adapter.DB
}

var DatabaseIsNotOpenedError = errors.New("database is not opened")
var ds *dataSource
var dsOnce sync.Once

func GetDataSource() DataSource {
	dsOnce.Do(func() {
		ds = new(dataSource)
		ds.gorm = new(adapter.GormDataSource)
	})
	return ds
}

func (d *dataSource) Open(p *properties) error {
	var err error
	password := p.Password
	if p.Config.Decrypt {
		pwd, err := rsa.DecryptBase64([]byte(password), []byte(p.Config.DecryptKey))
		if err == nil {
			password = string(pwd)
		}
	}

	source := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=%v&loc=%v",
		p.Username, p.Password, p.Host, p.Port,  p.Database, p.Charset, p.ParseTime, p.Loc)

	d.db, err = d.gorm.Open(p.Type, source)

	if err != nil {
		d.db = nil
		defer d.gorm.Close()
		return err
	}

	return nil
}

func (d *dataSource) IsOpened() bool {
	return d.db != nil
}

func (d *dataSource) Close() error {
	if d.db != nil {
		err := d.gorm.Close()
		d.db = nil
		return err
	}
	return DatabaseIsNotOpenedError
}

func (d *dataSource) DB() adapter.DB {
	return d.db
}

