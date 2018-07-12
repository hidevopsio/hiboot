package services

import (
	"testing"
	"github.com/hidevopsio/hiboot/examples/db/bolt/domain"
	"github.com/stretchr/testify/assert"
	"encoding/json"
)

var userService *UserService

type FakeRepository struct {}

func (r *FakeRepository) Put(key, value []byte) error  {
	return nil
}

func (r *FakeRepository) Get(key []byte) ([]byte, error)  {
	user := &domain.User{Name: "John Doe", Age: 18}
	u, err := json.Marshal(user)
	return u, err
}

func (r *FakeRepository) Delete(key []byte) error {
	return nil
}

func init() {
	userService = &UserService{Repository: &FakeRepository{}}
}

func TestAddUser(t *testing.T) {
	user := &domain.User{Name: "John Doe", Age: 18}
	err := userService.AddUser(user)
	assert.Equal(t, nil, err)
}

func TestGetUserUser(t *testing.T) {
	u, err := userService.GetUser("")
	assert.Equal(t, nil, err)
	assert.Equal(t, "John Doe", u.Name)
	assert.Equal(t, 18, u.Age)
}

func TestDeleteUser(t *testing.T) {
	err := userService.DeleteUser("")
	assert.Equal(t, nil, err)
}