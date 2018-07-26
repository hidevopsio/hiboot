
package services

import (
	"fmt"
	"github.com/hidevopsio/hiboot/examples/data/bolt/model"
	"github.com/hidevopsio/hiboot/pkg/starter/data/bolt"
)

type UserService struct {
	Repository bolt.Repository `inject:"boltRepository"`
}

func (us *UserService) AddUser(user *model.User) error {
	if us.Repository == nil {
		return fmt.Errorf("repository is not injected")
	}

	return us.Repository.Put(user)
}

func (us *UserService) GetUser(id string) (*model.User, error) {
	if us.Repository == nil {
		return nil, fmt.Errorf("repository is not injected")
	}
	var user model.User
	err := us.Repository.Get(id, &user)
	return &user, err
}

func (us *UserService) DeleteUser(id string) error {
	if us.Repository == nil {
		return fmt.Errorf("repository is not injected")
	}
	return us.Repository.Delete(id, &model.User{})
}

