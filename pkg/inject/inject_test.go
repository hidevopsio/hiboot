package inject

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
	"github.com/hidevopsio/hiboot/pkg/starter/data"
	"github.com/hidevopsio/hiboot/pkg/starter"
)

type user struct {
	Name string
}

type fakeRepository struct {
	data.BaseKVRepository
}

type fakeConfiguration struct{
}

func (c *fakeConfiguration) NewRepository(name string) data.Repository {
	repo := new(fakeRepository)
	repo.SetName(name)
	return repo
}

type fooConfiguration struct{
}

type fooService struct {
	FooRepository data.Repository `inject:"fooRepository,dataSourceType=foo,table=foo"`
}

type barService struct {
	FooRepository data.Repository `inject:"barRepository,dataSourceType=foo"`
}

type userService struct {
	User           *user           `inject:"user"`
	UserRepository data.KVRepository `inject:"userRepository,dataSourceType=fake,namespace=user"`
	Url            string          `value:"${fake.url:http://localhost:8080}"`
}

type fooBarService struct {
	FooBarRepository data.KVRepository `inject:"foobarRepository,dataSourceType=foo,namespace=foobar"`
}

type foobarRecursiveInject struct {
	Foobar *fooBarService `inject:"foobarService"`
}

type recursiveInject struct {
	UserService *userService
}

func init() {
	starter.Add("fake", fakeConfiguration{})
	starter.Add("foo", fooConfiguration{})
	starter.GetAutoConfiguration().Build()
}

func TestNotInject(t *testing.T) {
	baz := new(userService)
	assert.Equal(t, (*user)(nil), baz.User)
}

func TestInject(t *testing.T) {
	t.Run("should inject repository", func(t *testing.T) {
		us := new(userService)
		err := IntoObject(reflect.ValueOf(us))
		assert.Equal(t, nil, err)
		assert.NotEqual(t, (*user)(nil), us.User)
		assert.NotEqual(t, (*fakeRepository)(nil), us.UserRepository)
	})

	t.Run("should not inject foobar repository", func(t *testing.T) {
		fb := new(foobarRecursiveInject)
		err := IntoObject(reflect.ValueOf(fb))
		assert.Equal(t, "method NewRepository(name string) is not implemented", err.Error())
	})


	t.Run("should not inject repository with invalid configuration", func(t *testing.T) {
		fs := new(fooService)
		err := IntoObject(reflect.ValueOf(fs))
		assert.Equal(t, "method NewRepository(name string) is not implemented", err.Error())
	})

	t.Run("should not inject repository with invalid repository name", func(t *testing.T) {
		bs := new(barService)
		err := IntoObject(reflect.ValueOf(bs))
		assert.Equal(t, "namespace or table name does not specified", err.Error())
	})

	t.Run("should inject recursively", func(t *testing.T) {
		ps := &recursiveInject{UserService: new(userService)}
		err := IntoObject(reflect.ValueOf(ps))
		assert.Equal(t, nil, err)
		assert.NotEqual(t, (*user)(nil), ps.UserService.User)
		assert.NotEqual(t, (*fakeRepository)(nil), ps.UserService.UserRepository)
	})
}
