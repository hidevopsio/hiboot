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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hidevopsio/hiboot/examples/data/etcd/entity"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/starter/data/etcd"
	"time"
)

type UserService interface {
	AddUser(id string, user *entity.User) (err error)
	GetUser(id string) (user *entity.User, err error)
	DeleteUser(id string) (err error)
}

type UserServiceImpl struct {
	// add UserService, it means that the instance of UserServiceImpl can be found by UserService
	UserService
	repository etcd.Repository
}

func init() {
	// register UserServiceImpl
	app.Component(newUserService)
}

// will inject BoltRepository that configured in github.com/hidevopsio/hiboot/pkg/starter/data/bolt
func newUserService(repository etcd.Repository) UserService {
	return &UserServiceImpl{
		repository: repository,
	}
}

func (s *UserServiceImpl) AddUser(id string, user *entity.User) (err error) {
	if user == nil {
		return errors.New("user is not allowed nil")
	}
	userBuf, _ := json.Marshal(user)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	res, err := s.repository.Put(ctx, id, string(userBuf))
	cancel()
	if err != nil {
		fmt.Println("failed to put data to etcd, err:", err)
		return err
	}

	log.Debug(res)

	return nil
}

func (s *UserServiceImpl) GetUser(id string) (user *entity.User, err error) {
	user = &entity.User{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := s.repository.Get(ctx, id)
	cancel()
	if err != nil {
		log.Debugf("failed to get data from etcd, err: %v", err)
		return nil, err
	}

	if err = json.Unmarshal(resp.Kvs[0].Value, &user); err != nil {
		log.Debugf("failed to unmarshal data, err: %v", err)
		return nil, err
	}

	return
}

func (s *UserServiceImpl) DeleteUser(id string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = s.repository.Delete(ctx, id)
	cancel()
	return
}
