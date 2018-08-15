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
	"github.com/hidevopsio/hiboot/examples/data/gorm/entity"
	"github.com/hidevopsio/hiboot/pkg/starter/data/gorm"
)


type UserService struct {
	repository gorm.Repository
}

// will inject BoltRepository that configured in github.com/hidevopsio/hiboot/pkg/starter/data/bolt
func (s *UserService) Init(repository gorm.Repository)  {
	s.repository = repository
}

func (s *UserService) AddUser(user *entity.User) error {
	return nil
}

func (s *UserService) GetUser(id string) (user *entity.User, err error) {
	return
}

func (s *UserService) DeleteUser(id string) (err error) {
	return
}

