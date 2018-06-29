package services

import (
	"github.com/hidevopsio/hiboot/pkg/starter/db"
	"encoding/json"
	"github.com/hidevopsio/hiboot/examples/db/bolt/domain"
)

type UserService struct {
	Repository db.KVRepository `component:"repository" dataSourceType:"bolt"`
}

func (us *UserService) AddUser(user *domain.User) error {
	u, err := json.Marshal(user)
	if err == nil {
		us.Repository.Put([]byte("user"), []byte(user.Id), u)
	}
	return err
}

func (us *UserService) GetUser(id string) (*domain.User, error) {
	u, err := us.Repository.Get([]byte("user"), []byte(id))
	if err != nil {
		return nil, err
	}
	var user domain.User
	err = json.Unmarshal(u, &user)
	return &user, err
}

func (us *UserService) DeleteUser(id string) error {
	return us.Repository.Delete([]byte("user"), []byte(id))
}

