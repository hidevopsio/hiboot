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


func TestNotInject(t *testing.T) {
	baz := new(UserService)
	assert.Equal(t, (*User)(nil), baz.User)
	assert.Equal(t, (db.KVRepository)(nil), baz.UserRepository)
}

func TestInject(t *testing.T) {
	dataSources := make(map[string]interface{})

	factory := new(db.DataSourceFactory)
	bolt, err := factory.New(db.DataSourceTypeBolt)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, bolt)

	dataSources["bolt"] = bolt

	us := new(UserService)
	instances := make(map[string]interface{})
	IntoObject(reflect.ValueOf(us), dataSources, instances)
	assert.NotEqual(t, (*User)(nil), us.User)
	assert.NotEqual(t, (db.KVRepository)(nil), us.UserRepository)
	assert.Equal(t, instances["user"].(*User), us.User)
	assert.Equal(t, instances["userRepository"].(db.KVRepository), us.UserRepository)
}