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
)

func TestConfiguration(t *testing.T) {

	// TODO: should test with fake data source
	conf := new(configuration)
	conf.GormProperties = properties{
		Type:      "mysql",
		Host:      "mysql-dev",
		Port:      "3306",
		Username:  "test",
		Password:  "LcNxqoI4zZjAnpiTD7JQxLJR/IgL2iTiSZ2nd7KPEBgxMV+FVhPSzM+fgH93XqZJNpboN4F/buX22yLTXK38AcVGTfID3rmQAOAc9A2DIWNy5v9+3NOY00M8z4dR1XHojheK0681cY9QVjtlJ70jFFDXb7PjFc2fQ0GIyIjBQDY=",
		Database:  "test",
		ParseTime: "True",
		Charset:   "utf8",
		Loc:       "Asia%2FShanghai",
		Config: Config{
			Decrypt: true,
		},
	}

	repo := conf.GormRepository()
	assert.Equal(t, nil, repo)
}
