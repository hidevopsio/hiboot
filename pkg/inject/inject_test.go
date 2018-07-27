package inject

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
	"github.com/hidevopsio/hiboot/pkg/starter/data"
	"github.com/hidevopsio/hiboot/pkg/starter"
	"github.com/hidevopsio/hiboot/pkg/utils"
	"path/filepath"
	"os"
)

type user struct {
	Name string
	App  string
}

type fakeRepository struct {
	data.BaseKVRepository
}

type fakeProperties struct {
	Name     string
	Nickname string
	Username string
	Url      string
}

type fakeConfiguration struct {
	Properties fakeProperties `mapstructure:"fake"`
}

type fakeDataSource struct {
}

func (c *fakeConfiguration) FakeRepository() data.Repository {
	repo := new(fakeRepository)
	repo.SetDataSource(new(fakeDataSource))
	return repo
}

// FooUser an instance fooUser is injectable with tag `inject:"fooUser"`
func (c *fakeConfiguration) FooUser() *user {
	u := new(user)
	u.Name = "foo"
	return u
}

type fooConfiguration struct {
}

type fooService struct {
	FooUser       *user           `inject:"name=foo"`
	FooRepository data.Repository `inject:""`
}

type hibootService struct {
	HibootUser *user `inject:"name=${app.name}"`
}

type barService struct {
	FooRepository data.Repository `inject:""`
}

type userService struct {
	User           *user           `inject:""`
	FooUser        *user           `inject:"name=foo"`
	FakeUser       *user           `inject:"name=${fake.name},app=${app.name}"`
	FakeRepository data.Repository `inject:""`
	DefaultUrl     string          `value:"${fake.defaultUrl:http://localhost:8080}"`
	Url            string          `value:"${fake.url}"`
}

type fooBarService struct {
	FooBarRepository data.Repository `inject:""`
}

type foobarRecursiveInject struct {
	FoobarService *fooBarService `inject:""`
}

type recursiveInject struct {
	UserService *userService
}

var (
	appName    = "hiboot"
	fakeName   = "fake"
	fooName    = "foo"
	fakeUrl    = "http://fake.com/api/foo"
	defaultUrl = "http://localhost:8080"
)

func init() {
	utils.EnsureWorkDir("../..")

	configPath := filepath.Join(utils.GetWorkDir(), "config")
	fakeFile := "application-fake.yaml"
	os.Remove(filepath.Join(configPath, fakeFile))
	fakeContent :=
		"fake:" +
			"\n  name: " + fakeName +
			"\n  nickname: ${app.name} ${fake.name}\n" +
			"\n  username: ${unknown.name:bar}\n" +
			"\n  url: " + fakeUrl
	utils.WriterFile(configPath, fakeFile, []byte(fakeContent))

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
		assert.Equal(t, fooName, us.FooUser.Name)
		assert.Equal(t, fakeName, us.FakeUser.Name)
		assert.Equal(t, appName, us.FakeUser.App)
		assert.Equal(t, fakeUrl, us.Url)
		assert.Equal(t, defaultUrl, us.DefaultUrl)
		assert.NotEqual(t, (*fakeRepository)(nil), us.FakeRepository)
	})

	t.Run("should not inject unimplemented interface into FooBarRepository", func(t *testing.T) {
		fb := new(foobarRecursiveInject)
		err := IntoObject(reflect.ValueOf(fb))
		assert.Contains(t, err.Error(), "data.Repository is not implemented")
	})

	t.Run("should not inject unimplemented interface into FooRepository", func(t *testing.T) {
		fs := new(fooService)
		err := IntoObject(reflect.ValueOf(fs))
		assert.Equal(t, "foo", fs.FooUser.Name)
		assert.Contains(t, err.Error(), "data.Repository is not implemented")
	})

	t.Run("should not inject system property into object", func(t *testing.T) {
		fs := new(hibootService)
		err := IntoObject(reflect.ValueOf(fs))
		assert.Equal(t, nil, err)
		assert.Equal(t, appName, fs.HibootUser.Name)
	})

	t.Run("should not inject unimplemented interface into BarRepository", func(t *testing.T) {
		bs := new(barService)
		err := IntoObject(reflect.ValueOf(bs))
		assert.Contains(t, err.Error(), "data.Repository is not implemented")
	})

	t.Run("should inject recursively", func(t *testing.T) {
		ps := &recursiveInject{UserService: new(userService)}
		err := IntoObject(reflect.ValueOf(ps))
		assert.Equal(t, nil, err)
		assert.NotEqual(t, (*user)(nil), ps.UserService.User)
		assert.NotEqual(t, (*fakeRepository)(nil), ps.UserService.FakeRepository)
	})
}
