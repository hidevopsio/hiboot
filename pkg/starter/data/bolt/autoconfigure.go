package bolt

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

import (
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/app"
)

type configuration struct {
	// the properties member name must be Bolt if the mapstructure is bolt,
	// so that the reference can be parsed
	BoltProperties properties `mapstructure:"bolt"`
}

func init() {
	app.AddConfig("bolt", configuration{})
}

func (c *configuration) dataSource() DataSource {
	dataSource := GetDataSource()
	if ! dataSource.IsOpened() {
		err := dataSource.Open(&c.BoltProperties)
		if err != nil {
			log.Error(err.Error())
		}
	}
	return dataSource
}

func (c *configuration) BoltRepository() Repository {
	repository := GetRepository()
	repository.SetDataSource(c.dataSource())
	return repository
}

