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
