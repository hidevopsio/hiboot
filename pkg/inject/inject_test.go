// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	"github.com/hidevopsio/hiboot/pkg/log"
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

type FakeRepository data.Repository

type FooUser struct {
	Name string
}

func (c *fakeConfiguration) FakeRepository() FakeRepository {
	repo := new(fakeRepository)
	repo.SetDataSource(new(fakeDataSource))
	return repo
}

// FooUser an instance fooUser is injectable with tag `inject:"fooUser"`
func (c *fakeConfiguration) FooUser() *FooUser {
	u := new(FooUser)
	u.Name = "foo"
	return u
}

type fooConfiguration struct {
}

type fooService struct {
	FooUser       *FooUser       `inject:"name=foo"` // TODO: should be able to change the instance name, e.g. `inject:"bazUser"`
	FooRepository FakeRepository `inject:""`
}

type hibootService struct {
	HibootUser *user `inject:"name=${app.name}"`
}

type barService struct {
	FooRepository FakeRepository `inject:""`
}

type userService struct {
	FooUser        *FooUser       `inject:"name=foo"`
	User           *user          `inject:""`
	FakeUser       *user          `inject:"name=${fake.name},app=${app.name}"`
	FakeRepository FakeRepository `inject:""`
	DefaultUrl     string         `value:"${fake.defaultUrl:http://localhost:8080}"`
	Url            string         `value:"${fake.url}"`
}

type sliceInjectionTestService struct {
	Profiles         []string       `value:"${app.profiles.include}"`
}

type fooBarService struct {
	FooBarRepository FakeRepository `inject:""`
}

type foobarRecursiveInject struct {
	FoobarService *fooBarService `inject:""`
}

type recursiveInject struct {
	UserService *userService
}

type MethodInjectionService struct {
	fooUser    *FooUser
	barUser    *user
	repository FakeRepository
}

type Baz struct {
	Name string `inject:""`
}

type testService struct {
	baz *Baz
}

func (s *testService) Init(baz *Baz)  {
	s.baz = baz
}

type buzz struct {
	Name string
}

type buzzService struct {
	bz *buzz
}

func (s *buzzService) Init(bs *buzzService, bz *buzz)  {
	s.bz = bz
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
	starter.GetFactory().Build()
}

// Init automatically inject FooUser and FakeRepository that instantiated in fakeConfiguration
func (s *MethodInjectionService) Init(fooUser *FooUser, barUser *user, repository FakeRepository) {
	s.fooUser = fooUser
	s.barUser = barUser
	s.repository = repository
}

func TestNotInject(t *testing.T) {
	baz := new(userService)
	assert.Equal(t, (*user)(nil), baz.User)
}

func TestInject(t *testing.T) {
	t.Run("should inject through method", func(t *testing.T) {
		s := new(MethodInjectionService)
		err := IntoObject(reflect.ValueOf(s))
		assert.Equal(t, nil, err)
		assert.NotEqual(t, (*FooUser)(nil), s.fooUser)
		assert.NotEqual(t, (*user)(nil), s.barUser)
		assert.NotEqual(t, (FakeRepository)(nil), s.repository)
	})

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
		assert.Equal(t, nil, err)
	})

	t.Run("should not inject unimplemented interface into FooRepository", func(t *testing.T) {
		fs := new(fooService)
		err := IntoObject(reflect.ValueOf(fs))
		assert.Equal(t, "foo", fs.FooUser.Name)
		assert.Equal(t, nil, err)
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
		assert.Equal(t, nil, err)
	})

	t.Run("should inject recursively", func(t *testing.T) {
		ps := &recursiveInject{UserService: new(userService)}
		err := IntoObject(reflect.ValueOf(ps))
		assert.Equal(t, nil, err)
		assert.NotEqual(t, (*user)(nil), ps.UserService.User)
		assert.NotEqual(t, (*fakeRepository)(nil), ps.UserService.FakeRepository)
	})

	t.Run("should not inject slice", func(t *testing.T) {
		testSvc := struct {
			Users []FooUser `inject:""`
		}{}
		err := IntoObject(reflect.ValueOf(testSvc))
		assert.Equal(t, nil, err)
	})

	t.Run("should inject slice value", func(t *testing.T) {
		testSvc := new(sliceInjectionTestService)
		err := IntoObject(reflect.ValueOf(testSvc))
		assert.Equal(t, nil, err)
		log.Debug(testSvc.Profiles)
	})

	t.Run("should inject slice value", func(t *testing.T) {
		err := IntoObject(reflect.ValueOf((*MethodInjectionService)(nil)))
		assert.Equal(t, InvalidObjectError, err)
	})

	t.Run("should inject slice value", func(t *testing.T) {
		err := IntoObject(reflect.ValueOf((*string)(nil)))
		assert.Equal(t, InvalidObjectError, err)
	})

	t.Run("should failed to inject with illegal struct tag", func(t *testing.T) {
		err := IntoObject(reflect.ValueOf(new(testService)))
		assert.Equal(t, nil, err)
	})

	t.Run("should failed to inject if the type of param and receiver are the same", func(t *testing.T) {
		err := IntoObject(reflect.ValueOf(new(buzzService)))
		assert.Equal(t, nil, err)
	})

	t.Run("should inject value and it must not be singleton", func(t *testing.T) {
		a := &struct{TestName string `value:"test-data-from-a"`}{}
		b := &struct{TestName string}{}
		err := IntoObject(reflect.ValueOf(a))
		assert.Equal(t, nil, err)
		IntoObject(reflect.ValueOf(b))
		assert.Equal(t, nil, err)

		assert.NotEqual(t, a.TestName, b.TestName)
	})

	t.Run("should deduplicate tag", func(t *testing.T) {
		err := AddTag("test", new(BaseTag))
		assert.Equal(t, nil, err)
		err = AddTag("test", new(BaseTag))
		assert.Equal(t, TagIsAlreadyExistError, err)
		err = AddTag("nil", nil)
		assert.Equal(t, TagIsNilError, err)
	})

	t.Run("should not inject primitive type", func(t *testing.T) {
		a := struct{TestObj *int `inject:""`}{}
		err := IntoObject(reflect.ValueOf(&a))
		assert.Equal(t, nil, err)
	})
}
