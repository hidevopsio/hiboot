package services

import (
	"testing"
	"github.com/hidevopsio/hiboot/examples/data/bolt/model"
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
			u := params[1].(*model.User)
			u.Name = "John Doe"
			u.Age = 18
		}
	}

	return nil
}

func init() {
	userService = &UserService{Repository: &FakeRepository{}}
}

func TestAddUser(t *testing.T) {
	user := &model.User{Name: "John Doe", Age: 18}
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