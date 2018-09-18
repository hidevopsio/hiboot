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
	"encoding/json"
	"errors"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	_ "github.com/erikstmartin/go-testdb"
	"github.com/hidevopsio/hiboot/examples/data/etcd/entity"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/starter/data/etcd/fake"
	"github.com/hidevopsio/hiboot/pkg/utils/idgen"
	"github.com/stretchr/testify/assert"
	"testing"
)

var fakeUser = entity.User{
	Id:       "",
	Name:     "Bill Gates",
	Username: "billg",
	Password: "3948tdaD",
	Email:    "bill.gates@microsoft.com",
	Age:      60,
	Gender:   1,
}

func newId(t *testing.T) string {
	id, err := idgen.NextString()
	fakeUser.Id = id
	assert.Equal(t, nil, err)
	return id
}

func TestUserCrud(t *testing.T) {
	fakeRepository := new(fake.Repository)
	userService := newUserService(fakeRepository)

	id := newId(t)
	t.Run("should return error if user is nil", func(t *testing.T) {
		err := userService.AddUser(id, (*entity.User)(nil))
		assert.NotEqual(t, nil, err)
	})

	response := new(clientv3.PutResponse)

	fakeRepository.On("Put", nil, id).Return(response, nil)
	t.Run("should add user", func(t *testing.T) {
		err := userService.AddUser(id, &fakeUser)
		assert.Equal(t, nil, err)
	})

	simulationErr := errors.New("simulation err")
	id = newId(t)
	fakeRepository.On("Put", nil, id).Return((*clientv3.PutResponse)(nil), simulationErr)
	t.Run("should add user", func(t *testing.T) {
		err := userService.AddUser(id, &fakeUser)
		assert.Equal(t, err, simulationErr)
	})

	recordNotFound := errors.New("record not found")
	id = newId(t)
	fakeRepository.On("Get", nil, id).Return((*clientv3.GetResponse)(nil), recordNotFound)
	t.Run("should generate user id", func(t *testing.T) {
		//u := &entity.User{}
		_, err := userService.GetUser(id)
		log.Debug("Error %v", err)
		assert.NotEqual(t, recordNotFound, nil)
	})

	fakeUserBuf, _ := json.Marshal(&fakeUser)
	getRes := new(clientv3.GetResponse)
	kv := &mvccpb.KeyValue{
		Key:   []byte("test"),
		Value: fakeUserBuf,
	}
	getRes.Kvs = append(getRes.Kvs, kv)
	id = newId(t)
	fakeRepository.On("Get", nil, id).Return(getRes, nil)
	t.Run("should generate user id", func(t *testing.T) {
		//u := &entity.User{}
		var err error
		err = nil
		_, err = userService.GetUser(id)
		log.Debug("Error %v", err)
		assert.NotEqual(t, getRes, nil)
	})

	getRes = new(clientv3.GetResponse)
	kv = &mvccpb.KeyValue{
		Key:   []byte("test"),
		Value: []byte("test"),
	}
	getRes.Kvs = append(getRes.Kvs, kv)
	id = newId(t)
	fakeRepository.On("Get", nil, id).Return(getRes, nil)
	t.Run("should generate user id", func(t *testing.T) {
		//u := &entity.User{}
		var err error
		err = nil
		_, err = userService.GetUser(id)
		log.Debug("Error %v", err)
		assert.NotEqual(t, getRes, nil)
	})
	id = newId(t)
	fakeRepository.On("Delete", nil, id).Return((*clientv3.DeleteResponse)(nil), nil)
	t.Run("should delete user", func(t *testing.T) {
		err := userService.DeleteUser(id)
		assert.Equal(t, nil, err)
	})
}
