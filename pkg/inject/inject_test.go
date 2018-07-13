package inject

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
	"github.com/hidevopsio/hiboot/pkg/starter/db"
)

type User struct {
	Name string
}

type UserService struct{
	User           *User           `inject:"user"`
	UserRepository db.KVRepository `inject:"userRepository,dataSourceType=bolt,bucket=user"`
}

type ParentService struct{
	UserService *UserService
}

func TestNotInject(t *testing.T) {
	baz := new(UserService)
	assert.Equal(t, (*User)(nil), baz.User)
	assert.Equal(t, (db.KVRepository)(nil), baz.UserRepository)
}

func TestInject(t *testing.T) {
	dataSources := make(map[string]interface{})
	factory := new(db.DataSourceFactory)
	bolt, err := factory.New(db.DataSourceTypeBolt)
	instances := make(map[string]interface{})

	t.Run("test db factory", func(t *testing.T) {
		assert.Equal(t, nil, err)
		assert.NotEqual(t, nil, bolt)
		dataSources["bolt"] = bolt
	})

	t.Run("test inject repository", func(t *testing.T) {
		us := new(UserService)
		IntoObject(reflect.ValueOf(us), dataSources, instances)
		assert.NotEqual(t, (*User)(nil), us.User)
		assert.NotEqual(t, (db.KVRepository)(nil), us.UserRepository)
		assert.Equal(t, instances["user"].(*User), us.User)
		assert.Equal(t, instances["userRepository"].(db.KVRepository), us.UserRepository)
	})

	t.Run("test inject recursively", func(t *testing.T) {
		ps := &ParentService{UserService: new(UserService)}
		IntoObject(reflect.ValueOf(ps), dataSources, instances)
		assert.NotEqual(t, (*User)(nil), ps.UserService.User)
		assert.NotEqual(t, (db.KVRepository)(nil), ps.UserService.UserRepository)
		assert.Equal(t, instances["user"].(*User), ps.UserService.User)
		assert.Equal(t, instances["userRepository"].(db.KVRepository), ps.UserService.UserRepository)
	})
}