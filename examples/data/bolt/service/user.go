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
	"github.com/hidevopsio/hiboot/examples/data/bolt/entity"
	"github.com/hidevopsio/hiboot/pkg/starter/data/bolt"
)

type BoltRepository bolt.Repository

type UserService struct {
	repository BoltRepository
}

// will inject BoltRepository that configured in github.com/hidevopsio/hiboot/pkg/starter/data/bolt
func (s *UserService) Init(repository BoltRepository)  {
	s.repository = repository
}

func (s *UserService) AddUser(user *entity.User) error {
	return s.repository.Put(user)
}

func (s *UserService) GetUser(id string) (*entity.User, error) {
	var user entity.User
	err := s.repository.Get(id, &user)
	return &user, err
}

func (s *UserService) DeleteUser(id string) error {
	return s.repository.Delete(id, &entity.User{})
}

