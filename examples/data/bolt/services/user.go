
package services

import (
	"fmt"
	"encoding/json"
	"github.com/hidevopsio/hiboot/examples/data/bolt/domain"
	"github.com/hidevopsio/hiboot/pkg/starter/data/bolt"
)

type UserService struct {
	Repository bolt.Repository `inject:"userRepository,dataSourceType=bolt,namespace=user"`
}

func (us *UserService) AddUser(user *domain.User) error {
	if us.Repository == nil {
		return fmt.Errorf("repository is not injected")
	}
	u, err := json.Marshal(user)
	if err == nil {
		us.Repository.Put([]byte(user.Id), u)
	}
	return err
}

func (us *UserService) GetUser(id string) (*domain.User, error) {
	if us.Repository == nil {
		return nil, fmt.Errorf("repository is not injected")
	}
	u, err := us.Repository.Get([]byte(id))
	if err != nil {
		return nil, err
	}
	if len(u) == 0 {
		return nil, fmt.Errorf("user is not found")
	}
	var user domain.User
	err = json.Unmarshal(u, &user)
	return &user, err
}

func (us *UserService) DeleteUser(id string) error {
	if us.Repository == nil {
		return fmt.Errorf("repository is not injected")
	}
	return us.Repository.Delete([]byte(id))
}

