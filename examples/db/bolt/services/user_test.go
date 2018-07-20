package services

import (
	"testing"
	"github.com/hidevopsio/hiboot/examples/db/bolt/domain"
	"github.com/stretchr/testify/assert"
	"encoding/json"
	"errors"
)

var userService *UserService

type FakeRepository struct {}

func (r *FakeRepository) Put(key, value []byte) error  {
	return nil
}

func (r *FakeRepository) Get(key []byte) ([]byte, error)  {
	var u []byte
	var err error
	if string(key[:]) == "1" {
		user := &domain.User{Name: "John Doe", Age: 18}
		u, err = json.Marshal(user)
	} else if string(key[:]) == "" {
		err = errors.New("wrong user ID")
	} else {
		u = []byte("")
	}
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

func TestGetUser(t *testing.T) {
	u, err := userService.GetUser("1")
	assert.Equal(t, nil, err)
	assert.Equal(t, "John Doe", u.Name)
	assert.Equal(t, 18, u.Age)

	_, err = userService.GetUser("")
	assert.Equal(t, "wrong user ID", err.Error())

	_, err = userService.GetUser("2")
	assert.Equal(t, "user is not found", err.Error())
}

func TestDeleteUser(t *testing.T) {
	err := userService.DeleteUser("")
	assert.Equal(t, nil, err)
}

func TestNilRepository(t *testing.T) {
	errMsg := "repository is not injected"
	us := &UserService{}

	err := us.AddUser(nil)
	assert.Equal(t, errMsg, err.Error())

	_, err = us.GetUser("")
	assert.Equal(t, errMsg, err.Error())

	err = us.DeleteUser("")
	assert.Equal(t, errMsg, err.Error())
}