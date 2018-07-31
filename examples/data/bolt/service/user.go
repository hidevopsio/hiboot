
package service

import (
	"github.com/hidevopsio/hiboot/examples/data/bolt/entity"
	"github.com/hidevopsio/hiboot/pkg/starter/data/bolt"
)

type BoltRepository bolt.Repository

type UserService struct {
	repository BoltRepository
}

// will inject BoltRepository that configured in github.com/hidevopsio/hiboot/pkg/starter/data/bolt
func (s *UserService) Init(repository BoltRepository)  {
	s.repository = repository
}

func (s *UserService) AddUser(user *entity.User) error {
	return s.repository.Put(user)
}

func (s *UserService) GetUser(id string) (*entity.User, error) {
	var user entity.User
	err := s.repository.Get(id, &user)
	return &user, err
}

func (s *UserService) DeleteUser(id string) error {
	return s.repository.Delete(id, &entity.User{})
}

