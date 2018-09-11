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
	_ "github.com/erikstmartin/go-testdb"
	"github.com/hidevopsio/gorm"
	"github.com/hidevopsio/hiboot/examples/data/gorm/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

var fakeUser = entity.User{
	Id:       1,
	Name:     "Bill Gates",
	Username: "billg",
	Password: "3948tdaD",
	Email:    "bill.gates@microsoft.com",
	Age:      60,
	Gender:   1,
}

func TestUserCrud(t *testing.T) {
	fakeRepository := new(gorm.FakeRepository)
	userService := newUserService(fakeRepository)

	t.Run("should return error if user is nil", func(t *testing.T) {
		err := userService.AddUser((*entity.User)(nil))
		assert.NotEqual(t, nil, err)
	})

	t.Run("should add user", func(t *testing.T) {
		err := userService.AddUser(&fakeUser)
		assert.Equal(t, nil, err)
	})

	t.Run("should generate user id", func(t *testing.T) {
		u := &entity.User{}
		err := userService.AddUser(u)
		assert.Equal(t, nil, err)
		assert.NotEqual(t, 0, u.Id)
	})

	t.Run("should generate user id", func(t *testing.T) {
		fakeRepository.Mock("Find", &[]entity.User{fakeUser}).Expect(nil)
		users, err := userService.GetAll()
		assert.Equal(t, nil, err)
		assert.Equal(t, 1, len(*users))
		assert.Equal(t, "Bill Gates", (*users)[0].Name)
	})

	t.Run("should get user that added above", func(t *testing.T) {
		// call mock method mocker.First(fakeUser).Expected(nil)
		fakeRepository.Mock("First", &fakeUser).Expect(nil)

		u, err := userService.GetUser(1)
		assert.Equal(t, nil, err)
		assert.Equal(t, "Bill Gates", u.Name)
		assert.Equal(t, uint(60), u.Age)
	})

	t.Run("should delete user", func(t *testing.T) {
		err := userService.DeleteUser(1)
		assert.Equal(t, nil, err)
	})
}
