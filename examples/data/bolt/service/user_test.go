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

package service

import (
	"testing"
	"github.com/hidevopsio/hiboot/examples/data/bolt/entity"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/starter/data"
)

var userService *UserService

type FakeRepository struct {
	data.BaseKVRepository
}

func (r *FakeRepository) Get(params ...interface{}) error  {
	if len(params) == 2 {
		key := params[0].(string)
		if key == "1" {
			u := params[1].(*entity.User)
			u.Name = "John Doe"
			u.Age = 18
		}
	}

	return nil
}

func init() {
	userService = new(UserService)
	userService.Init(&FakeRepository{})
}

func TestAddUser(t *testing.T) {
	user := &entity.User{Name: "John Doe", Age: 18}
	err := userService.AddUser(user)
	assert.Equal(t, nil, err)
}

func TestGetUser(t *testing.T) {
	u, err := userService.GetUser("1")
	assert.Equal(t, nil, err)
	assert.Equal(t, "John Doe", u.Name)
	assert.Equal(t, 18, u.Age)
}

func TestDeleteUser(t *testing.T) {
	err := userService.DeleteUser("")
	assert.Equal(t, nil, err)
}
