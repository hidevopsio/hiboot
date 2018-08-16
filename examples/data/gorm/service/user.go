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
	"github.com/hidevopsio/hiboot/pkg/utils/idgen"
	"github.com/hidevopsio/hiboot/pkg/starter/data/gorm"
)

type UserService struct {
	repository gorm.Repository
}

// will inject BoltRepository that configured in github.com/hidevopsio/hiboot/pkg/starter/data/bolt
func (s *UserService) Init(repository gorm.Repository)  {
	s.repository = repository
	repository.AutoMigrate(&entity.User{})
}

func (s *UserService) AddUser(user *entity.User) (err error) {
	if user.Id == 0 {
		user.Id, err = idgen.Next()
		if err != nil {
			return
		}
	}
	err = s.repository.Create(user).Error
	return
}

func (s *UserService) GetUser(id uint64) (user *entity.User, err error) {
	user = &entity.User{}
	err = s.repository.Where("id = ?", id).First(user).Error
	return
}

func (s *UserService) DeleteUser(id uint64) (err error) {
	err = s.repository.Where("id = ?", id).Delete(entity.User{}).Error
	return
}

