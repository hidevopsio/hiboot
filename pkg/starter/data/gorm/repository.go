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
	"github.com/hidevopsio/hiboot/pkg/starter/data"
	"sync"
	"github.com/hidevopsio/hiboot/pkg/starter/data/gorm/adapter"
)

type Repository interface {
	data.Repository
}

type repository struct {
	data.BaseRepository
	dataSource DataSource
	db adapter.DB
}

var repo *repository
var once sync.Once

func GetRepository() *repository {
	once.Do(func() {
		repo = &repository{}
	})
	return repo
}


// Open bolt database
func (r *repository) SetDataSource(d interface{})  {
	if d != nil {
		r.dataSource = d.(DataSource)
		r.db = r.dataSource.DB()
	}
}

func (r *repository) DataSource() interface{}  {
	return r.dataSource
}

func (r *repository) CloseDataSource() error {
	if r.dataSource != nil {
		return r.dataSource.Close()
	}
	return data.InvalidDataSourceError
}