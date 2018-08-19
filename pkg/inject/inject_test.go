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

package inject_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
	"github.com/hidevopsio/hiboot/pkg/starter/data"
	"github.com/hidevopsio/hiboot/pkg/starter"
	"path/filepath"
	"os"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/inject"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
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

type fakeTag struct{
	inject.BaseTag
}


type Tag struct{
	inject.BaseTag
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

type UserService interface {
	Get() string
}

type userService struct {
	FooUser        *FooUser       `inject:"name=foo"`
	User           *user          `inject:""`
	FakeUser       *user          `inject:"name=${fake.name},app=${app.name}"`
	FakeRepository FakeRepository `inject:""`
	DefaultUrl     string         `value:"${fake.defaultUrl:http://localhost:8080}"`
	Url            string         `value:"${fake.url}"`
}

func (s *userService) Get() string {
	return "Hello, world"
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

type BazService interface {
	GetNickname() string
}

type BazImpl struct {
	BazService
	Name string `inject:""`
	nickname string
}

func (s *BazImpl) GetNickname() string  {
	return s.nickname
}

type testService struct {
	baz BazService
	name string
}

func (s *testService) Init(baz BazService)  {
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

type CatInterface interface {

}

type animalService struct {
	cat CatInterface
}

func (a *animalService) Init(cat CatInterface)  {
	a.cat = cat
}


type testTag struct{
	inject.BaseTag
}

var (
	appName    = "hiboot"
	fakeName   = "fake"
	fooName    = "foo"
	fakeUrl    = "http://fake.com/api/foo"
	defaultUrl = "http://localhost:8080"
)

func init() {
	io.EnsureWorkDir("../..")

	configPath := filepath.Join(io.GetWorkDir(), "config")
	fakeFile := "application-fake.yaml"
	os.Remove(filepath.Join(configPath, fakeFile))
	fakeContent :=
		"fake:" +
			"\n  name: " + fakeName +
			"\n  nickname: ${app.name} ${fake.name}\n" +
			"\n  username: ${unknown.name:bar}\n" +
			"\n  url: " + fakeUrl
	io.WriterFile(configPath, fakeFile, []byte(fakeContent))

	starter.AddConfig("fake", fakeConfiguration{})
	starter.AddConfig("foo", fooConfiguration{})
	starter.GetFactory().Build()
	starter.Add(new(BazImpl))
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
		err := inject.IntoObject(reflect.ValueOf(s))
		assert.Equal(t, nil, err)
		assert.NotEqual(t, (*FooUser)(nil), s.fooUser)
		assert.NotEqual(t, (*user)(nil), s.barUser)
		assert.NotEqual(t, (FakeRepository)(nil), s.repository)
	})

	t.Run("should inject repository", func(t *testing.T) {
		us := new(userService)
		err := inject.IntoObject(reflect.ValueOf(us))
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
		err := inject.IntoObject(reflect.ValueOf(fb))
		assert.Equal(t, nil, err)
	})

	t.Run("should not inject unimplemented interface into FooRepository", func(t *testing.T) {
		fs := new(fooService)
		err := inject.IntoObject(reflect.ValueOf(fs))
		assert.Equal(t, "foo", fs.FooUser.Name)
		assert.Equal(t, nil, err)
	})

	t.Run("should not inject system property into object", func(t *testing.T) {
		fs := new(hibootService)
		err := inject.IntoObject(reflect.ValueOf(fs))
		assert.Equal(t, nil, err)
		assert.Equal(t, appName, fs.HibootUser.Name)
	})

	t.Run("should not inject unimplemented interface into BarRepository", func(t *testing.T) {
		bs := new(barService)
		err := inject.IntoObject(reflect.ValueOf(bs))
		assert.Equal(t, nil, err)
	})

	t.Run("should inject recursively", func(t *testing.T) {
		ps := &recursiveInject{UserService: new(userService)}
		err := inject.IntoObject(reflect.ValueOf(ps))
		assert.Equal(t, nil, err)
		assert.NotEqual(t, (*user)(nil), ps.UserService.User)
		assert.NotEqual(t, (*fakeRepository)(nil), ps.UserService.FakeRepository)
	})

	t.Run("should not inject slice", func(t *testing.T) {
		testSvc := struct {
			Users []FooUser `inject:""`
		}{}
		err := inject.IntoObject(reflect.ValueOf(testSvc))
		assert.Equal(t, nil, err)
	})

	t.Run("should inject slice value", func(t *testing.T) {
		testSvc := new(sliceInjectionTestService)
		err := inject.IntoObject(reflect.ValueOf(testSvc))
		assert.Equal(t, nil, err)
		log.Debug(testSvc.Profiles)
	})

	t.Run("should inject slice value", func(t *testing.T) {
		err := inject.IntoObject(reflect.ValueOf((*MethodInjectionService)(nil)))
		assert.Equal(t, inject.InvalidObjectError, err)
	})

	t.Run("should inject slice value", func(t *testing.T) {
		err := inject.IntoObject(reflect.ValueOf((*string)(nil)))
		assert.Equal(t, inject.InvalidObjectError, err)
	})

	t.Run("should failed to inject with illegal struct tag", func(t *testing.T) {
		err := inject.IntoObject(reflect.ValueOf(new(testService)))
		assert.Equal(t, nil, err)
	})

	t.Run("should failed to inject if the type of param and receiver are the same", func(t *testing.T) {
		err := inject.IntoObject(reflect.ValueOf(new(buzzService)))
		assert.Equal(t, nil, err)
	})

	t.Run("should skip inject if the type of param is an unimplemented interface", func(t *testing.T) {
		catSvc := new(animalService)
		err := inject.IntoObject(reflect.ValueOf(catSvc))
		assert.Equal(t, nil, err)
		assert.Equal(t, nil, catSvc.cat)
	})

	t.Run("should inject value and it must not be singleton", func(t *testing.T) {
		a := &struct{TestName string `value:"test-data-from-a"`}{}
		b := &struct{TestName string}{}
		err := inject.IntoObject(reflect.ValueOf(a))
		assert.Equal(t, nil, err)
		inject.IntoObject(reflect.ValueOf(b))
		assert.Equal(t, nil, err)

		assert.NotEqual(t, a.TestName, b.TestName)
	})

	t.Run("should deduplicate tag", func(t *testing.T) {
		err := inject.AddTag(new(testTag))
		assert.Equal(t, nil, err)
		err = inject.AddTag(new(testTag))
		assert.Equal(t, inject.TagIsAlreadyExistError, err)
		err = inject.AddTag( nil)
		assert.Equal(t, inject.InvalidTagNameError, err)
	})

	t.Run("should not inject primitive type", func(t *testing.T) {
		a := struct{TestObj *int `inject:""`}{}
		err := inject.IntoObject(reflect.ValueOf(&a))
		assert.Equal(t, nil, err)
	})

	t.Run("should add a fake tag", func(t *testing.T) {
		err := inject.AddTag(new(fakeTag))
		assert.Equal(t, nil, err)
	})

	t.Run("should failed to add a tag with empty name", func(t *testing.T) {
		err := inject.AddTag(new(Tag))
		assert.Equal(t, inject.InvalidTagNameError, err)
	})
}
