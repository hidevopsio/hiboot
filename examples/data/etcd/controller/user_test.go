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

package controller

import (
	"errors"
	"github.com/hidevopsio/hiboot/examples/data/etcd/entity"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/starter/data/etcd/fake"
	"github.com/hidevopsio/hiboot/pkg/utils/idgen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

var testUser entity.User

type fakeService struct {
	// add UserService, it means that the instance of UserServiceImpl can be found by UserService
	mock.Mock
}

func (s *fakeService) AddUser(id string, user *entity.User) (err error) {
	return
}

func (s *fakeService) GetUser(id string) (user *entity.User, err error) {
	args := s.Called(id)
	return args[0].(*entity.User), args.Error(1)
}

func (s *fakeService) DeleteUser(id string) (err error) {
	return
}

func (s *fakeService) GetAll() (users *[]entity.User, err error) {
	users = &[]entity.User{}
	args := s.Called()
	return args[0].(*[]entity.User), args.Error(1)
}

func TestCrdRequest(t *testing.T) {
	app.Component(new(fake.Repository))
	svc := new(fakeService)
	userController := newUserController(svc)
	testApp := web.NewTestApplication(t, userController)

	// TODO: if os.Getenv("INTEGRATION_TEST") == false {

	id, err := idgen.NextString()
	assert.Equal(t, nil, err)

	testUser = entity.User{
		Id:       id,
		Name:     "Bill Gates",
		Username: "billg",
		Password: "3948tdaD",
		Email:    "bill.gates@microsoft.com",
		Age:      60,
		Gender:   1,
	}

	t.Run("should add user with POST request", func(t *testing.T) {
		// First, let's Post User
		testApp.Post("/user").
			WithJSON(testUser).
			Expect().Status(http.StatusOK)
	})

	svc.On("GetUser", id).Return(&testUser, nil)

	t.Run("should get user with GET request", func(t *testing.T) {
		// Then Get User
		// e.g. GET /user/id/123456
		testApp.Get("/user/id/{id}").
			WithPath("id", id).
			Expect().Status(http.StatusOK)
	})

	unknownId, err := idgen.NextString()
	assert.Equal(t, nil, err)
	svc.On("GetUser", unknownId).Return((*entity.User)(nil), errors.New("not found"))

	t.Run("should return 404 if trying to find a record that does not exist", func(t *testing.T) {
		// Then Get User
		testApp.Get("/user/id/{id}").
			WithPath("id", unknownId).
			Expect().Status(http.StatusNotFound)
	})

	// assert that the expectations were met
	svc.AssertExpectations(t)

	id, err = idgen.NextString()
	t.Run("should delete the record with DELETE request", func(t *testing.T) {
		// Finally Delete User
		testApp.Delete("/user/id/{id}").
			WithPath("id", id).
			Expect().Status(http.StatusOK)
	})
}
