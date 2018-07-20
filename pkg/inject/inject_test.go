package inject

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
	"github.com/hidevopsio/hiboot/pkg/starter/db"
	"github.com/hidevopsio/hiboot/pkg/starter/db/bolt"
)

type User struct {
	Name string
}

type UserService struct{
	User           *User           `inject:"user"`
	UserRepository bolt.Repository `inject:"userRepository,dataSourceType=bolt,namespace=user"`
}

type ParentService struct{
	UserService *UserService
}

func TestNotInject(t *testing.T) {
	baz := new(UserService)
	assert.Equal(t, (*User)(nil), baz.User)
}

func TestInject(t *testing.T) {
	t.Run("test inject repository", func(t *testing.T) {
		us := new(UserService)
		IntoObject(reflect.ValueOf(us))
		assert.NotEqual(t, (*User)(nil), us.User)
	})

	t.Run("test inject recursively", func(t *testing.T) {
		ps := &ParentService{UserService: new(UserService)}
		IntoObject(reflect.ValueOf(ps))
		assert.NotEqual(t, (*User)(nil), ps.UserService.User)
		assert.NotEqual(t, (db.Repository)(nil), ps.UserService.UserRepository)
	})
}