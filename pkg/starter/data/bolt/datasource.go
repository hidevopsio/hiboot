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
	"github.com/boltdb/bolt"
	"github.com/hidevopsio/hiboot/pkg/log"
	"sync"
	"time"
)

type DataSource interface {
	Open(properties *properties) error
	Close() error
	IsOpened() bool
	DB() *bolt.DB
}

type dataSource struct {
	db *bolt.DB
}

var ds *dataSource
var do sync.Once

func GetDataSource() *dataSource {
	do.Do(func() {
		ds = new(dataSource)
	})
	return ds
}

func (d *dataSource) DB() *bolt.DB {
	return d.db
}

func (d *dataSource) Open(properties *properties) error {
	if properties == nil {
		return InvalidPropertiesError
	}

	var err error
	d.db, err = bolt.Open(properties.Database,
		properties.Mode,
		&bolt.Options{Timeout: time.Duration(properties.Timeout) * time.Second},
	)

	if err != nil {
		if d.db != nil {
			defer d.db.Close()
		}
	} else {
		log.Infof("dataSource %v connected", properties.Database)
	}

	return err
}

// Close database
func (d *dataSource) Close() error {
	err := d.db.Close()
	d.db = nil
	return err
}

// IsOpened check if db is opened
func (d *dataSource) IsOpened() bool {
	return d.db != nil
}
