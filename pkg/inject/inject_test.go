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
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/factory/autoconfigure"
	"github.com/hidevopsio/hiboot/pkg/factory/instantiate"
	"github.com/hidevopsio/hiboot/pkg/inject"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

type User struct {
	Name string
	App  string
}

type fakeRepository struct {
}

type fakeProperties struct {
	Name          string   `default:"should not inject this default value as it will inject by system.Builder"`
	Nickname      string   `default:"should not inject this default value as it will inject by system.Builder"`
	Username      string   `default:"should not inject this default value as it will inject by system.Builder"`
	Url           string   `default:"should not inject this default value as it will inject by system.Builder"`
	DefStrVal     string   `default:"this is default value"`
	DefIntVal     int      `default:"123"`
	DefIntVal8    int8     `default:"12"`
	DefIntVal16   int16    `default:"123"`
	DefUintVal32  int32    `default:"1234"`
	DefUintVal64  int64    `default:"12345"`
	DefIntValU    uint     `default:"123"`
	DefIntValU8   uint8    `default:"12"`
	DefIntValU16  uint16   `default:"123"`
	DefUintValU32 uint32   `default:"1234"`
	DefUintValU64 uint64   `default:"12345"`
	DefFloatVal64 float64  `default:"0.1231"`
	DefFloatVal32 float32  `default:"0.1"`
	DefBool       bool     `default:"true"`
	DefSlice      []string `default:"jupiter,mercury,mars,earth,moon"`
}

type fakeConfiguration struct {
	app.Configuration
	Properties fakeProperties `mapstructure:"fake"`
}

type fakeDataSource struct {
}

type FakeRepository interface {
}

type FooUser struct {
	Name string
}

type fakeTag struct {
	inject.BaseTag
}

type Tag struct {
	inject.BaseTag
}

func (c *fakeConfiguration) FakeRepository() FakeRepository {
	repo := new(fakeRepository)
	return repo
}

// FooUser an instance fooUser is injectable with tag `inject:""` or constructor
func (c *fakeConfiguration) FooUser() *FooUser {
	u := new(FooUser)
	u.Name = "foo"
	return u
}

// BarUser an instance fooUser is injectable with tag `inject:""` or constructor
func (c *fakeConfiguration) User() *User {
	u := new(User)
	u.Name = "bar"
	return u
}

type fooConfiguration struct {
	app.Configuration
}

type fooService struct {
	FooUser       *FooUser       `inject:"name=foo"` // TODO: should be able to change the instance name, e.g. `inject:"bazUser"`
	FooRepository FakeRepository `inject:""`
}

type hibootService struct {
	HibootUser *User `inject:"name=${app.name}"`
}

type barService struct {
	FooRepository FakeRepository `inject:""`
}

type UserService interface {
	Get() string
}

type userService struct {
	FooUser        *FooUser       `inject:"name=foo"`
	User           *User          `inject:""`
	FakeUser       *User          `inject:"name=${fake.name},app=${app.name}"`
	FakeRepository FakeRepository `inject:""`
	DefaultUrl     string         `value:"${fake.defaultUrl:http://localhost:8080}"`
	Url            string         `value:"${fake.url}"`
	DefStrVal      string         `value:"this is value"`
	DefIntVal      int            `value:"123"`
	DefIntVal8     int8           `value:"12"`
	DefIntVal16    int16          `value:"123"`
	DefUintVal32   int32          `value:"1234"`
	DefUintVal64   int64          `value:"12345"`
	DefIntValU     uint           `value:"123"`
	DefIntValU8    uint8          `value:"12"`
	DefIntValU16   uint16         `value:"123"`
	DefUintValU32  uint32         `value:"1234"`
	DefUintValU64  uint64         `value:"12345"`
	DefFloatVal64  float64        `value:"0.1231"`
	DefFloatVal32  float32        `value:"0.1"`
	DefBool        bool           `value:"true"`
	DefSlice       []string       `value:"jupiter,mercury,mars,earth,moon"`
}

func (s *userService) Get() string {
	return "Hello, world"
}

type sliceInjectionTestService struct {
	Profiles []string `value:"${app.profiles.include}"`
	Options  []string `value:"a,b,c,d"`
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
	barUser    *User
	repository FakeRepository
}

type BazService interface {
	GetNickname() string
}

type BazImpl struct {
	BazService
	Name     string `inject:""`
	nickname string
}

func (s *BazImpl) GetNickname() string {
	return s.nickname
}

type testService struct {
	baz  BazService
	name string
}

func (s *testService) Init(baz BazService) {
	s.baz = baz
}

type buzz struct {
	Name string
}

type buzzService struct {
	bz *buzz
}

func (s *buzzService) Init(bs *buzzService, bz *buzz) {
	s.bz = bz
}

type CatInterface interface {
}

type animalService struct {
	cat CatInterface
}

func (a *animalService) Init(cat CatInterface) {
	a.cat = cat
}

type testTag struct {
	inject.BaseTag
}

var (
	appName    = "hiboot"
	fakeName   = "fake"
	fooName    = "foo"
	fakeUrl    = "http://fake.com/api/foo"
	defaultUrl = "http://localhost:8080"

	cf *autoconfigure.ConfigurableFactory
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func newMethodInjectionService(fooUser *FooUser, barUser *User, repository FakeRepository) *MethodInjectionService {
	return &MethodInjectionService{
		fooUser:    fooUser,
		barUser:    barUser,
		repository: repository,
	}
}

func TestNotInject(t *testing.T) {
	baz := new(userService)
	assert.Equal(t, (*User)(nil), baz.User)
}

func TestInject(t *testing.T) {
	io.ChangeWorkDir(os.TempDir())
	configPath := filepath.Join(os.TempDir(), "config")

	fakeFile := "application.yml"
	os.Remove(filepath.Join(configPath, fakeFile))
	fakeContent :=
		"app:\n" +
			"  project: hidevopsio\n" +
			"  name: hiboot\n" +
			"  version: ${unknown.version:0.0.1}\n" +
			"  profiles:\n" +
			"    include:\n" +
			"    - foo\n" +
			"    - fake\n" +
			"\n"
	io.WriterFile(configPath, fakeFile, []byte(fakeContent))

	fakeFile = "application-fake.yml"
	os.Remove(filepath.Join(configPath, fakeFile))
	fakeContent =
		"fake:" +
			"\n  name: " + fakeName +
			"\n  nickname: ${app.name} ${fake.name}\n" +
			"\n  username: ${unknown.name:bar}\n" +
			"\n  url: " + fakeUrl
	io.WriterFile(configPath, fakeFile, []byte(fakeContent))

	instances := cmap.New()
	configurations := cmap.New()
	cf = new(autoconfigure.ConfigurableFactory)
	cf.InstantiateFactory = new(instantiate.InstantiateFactory)
	cf.InstantiateFactory.Initialize(instances, []*factory.MetaData{})
	cf.Initialize(configurations)
	cf.BuildSystemConfig()

	inject.SetFactory(cf)

	fakeConfig := new(fakeConfiguration)
	configs := []*factory.MetaData{
		factory.NewMetaData(fakeConfig),
		factory.NewMetaData(new(fooConfiguration)),
	}
	cf.Build(configs)
	cf.BuildComponents()

	t.Run("should inject default string", func(t *testing.T) {
		assert.Equal(t, "this is default value", fakeConfig.Properties.DefStrVal)
	})

	t.Run("should inject default int", func(t *testing.T) {
		assert.Equal(t, 123, fakeConfig.Properties.DefIntVal)
	})

	t.Run("should inject default uint", func(t *testing.T) {
		assert.Equal(t, uint(123), fakeConfig.Properties.DefIntValU)
	})

	t.Run("should inject default float32", func(t *testing.T) {
		assert.Equal(t, float32(0.1), fakeConfig.Properties.DefFloatVal32)
	})

	t.Run("should get config", func(t *testing.T) {
		fr := cf.GetInstance("fakeRepository")
		assert.NotEqual(t, nil, fr)
	})

	t.Run("should inject through method", func(t *testing.T) {
		fu := cf.GetInstance("fooUser")
		bu := cf.GetInstance("user")
		fp := cf.GetInstance("fakeRepository")
		u := new(User)

		cf.SetInstance("barUser", u)

		svc, err := inject.IntoFunc(newMethodInjectionService)
		assert.Equal(t, nil, err)
		s := svc.(*MethodInjectionService)
		assert.Equal(t, fu, s.fooUser)
		assert.Equal(t, bu, s.barUser)
		assert.Equal(t, fp, s.repository)
	})

	t.Run("should inject repository", func(t *testing.T) {
		us := new(userService)
		err := inject.IntoObject(us)
		assert.Equal(t, nil, err)
		assert.NotEqual(t, (*User)(nil), us.User)
		assert.Equal(t, fooName, us.FooUser.Name)
		assert.Equal(t, fakeName, us.FakeUser.Name)
		assert.Equal(t, appName, us.FakeUser.App)
		assert.Equal(t, fakeUrl, us.Url)
		assert.Equal(t, defaultUrl, us.DefaultUrl)
		assert.NotEqual(t, (*fakeRepository)(nil), us.FakeRepository)
	})

	t.Run("should not inject unimplemented interface into FooBarRepository", func(t *testing.T) {
		fb := new(foobarRecursiveInject)
		err := inject.IntoObject(fb)
		assert.Equal(t, nil, err)
	})

	t.Run("should not inject unimplemented interface into FooRepository", func(t *testing.T) {
		fs := new(fooService)
		err := inject.IntoObject(fs)
		assert.Equal(t, "foo", fs.FooUser.Name)
		assert.Equal(t, nil, err)
	})

	t.Run("should not inject system property into object", func(t *testing.T) {
		fs := new(hibootService)
		err := inject.IntoObject(fs)
		assert.Equal(t, nil, err)
		assert.Equal(t, appName, fs.HibootUser.Name)
	})

	t.Run("should not inject unimplemented interface into BarRepository", func(t *testing.T) {
		bs := new(barService)
		err := inject.IntoObject(bs)
		assert.Equal(t, nil, err)
	})

	t.Run("should inject recursively", func(t *testing.T) {
		ps := &recursiveInject{UserService: new(userService)}
		err := inject.IntoObject(ps)
		assert.Equal(t, nil, err)
		assert.NotEqual(t, (*User)(nil), ps.UserService.User)
		assert.NotEqual(t, (*fakeRepository)(nil), ps.UserService.FakeRepository)
	})

	t.Run("should not inject slice", func(t *testing.T) {
		testSvc := struct {
			Users []FooUser `inject:""`
		}{}
		err := inject.IntoObject(testSvc)
		assert.Equal(t, nil, err)
	})

	t.Run("should inject slice value", func(t *testing.T) {
		testSvc := new(sliceInjectionTestService)
		err := inject.IntoObject(testSvc)
		assert.Equal(t, nil, err)
		assert.Equal(t, []string{"foo", "fake"}, testSvc.Profiles)
		assert.Equal(t, []string{"a", "b", "c", "d"}, testSvc.Options)
		log.Debug(testSvc.Profiles)
	})

	t.Run("should inject slice value", func(t *testing.T) {
		err := inject.IntoObject((*MethodInjectionService)(nil))
		assert.Equal(t, inject.ErrInvalidObject, err)
	})

	t.Run("should inject slice value", func(t *testing.T) {
		err := inject.IntoObject((*string)(nil))
		assert.Equal(t, inject.ErrInvalidObject, err)
	})

	t.Run("should ignore to inject with invalid struct type BazService", func(t *testing.T) {
		ts := new(testService)
		err := inject.IntoObject(ts)
		assert.Equal(t, nil, err)
	})

	t.Run("should failed to inject if the type of param and receiver are the same", func(t *testing.T) {
		err := inject.IntoObject(new(buzzService))
		assert.Equal(t, nil, err)
	})

	t.Run("should skip inject if the type of param is an unimplemented interface", func(t *testing.T) {
		catSvc := new(animalService)
		err := inject.IntoObject(catSvc)
		assert.Equal(t, nil, err)
		assert.Equal(t, nil, catSvc.cat)
	})

	t.Run("should inject value and it must not be singleton", func(t *testing.T) {
		a := &struct {
			TestName string `value:"test-data-from-a"`
		}{}
		b := &struct{ TestName string }{}
		err := inject.IntoObject(a)
		assert.Equal(t, nil, err)
		inject.IntoObject(b)
		assert.Equal(t, nil, err)

		assert.NotEqual(t, a.TestName, b.TestName)
	})

	t.Run("should deduplicate tag", func(t *testing.T) {
		inject.AddTag(new(testTag))
		inject.AddTag(nil)
	})

	t.Run("should inject primitive type", func(t *testing.T) {
		type testObjTyp struct {
			TestObj *int `inject:""`
		}
		a := new(testObjTyp)
		err := inject.IntoObject(a)
		assert.Equal(t, nil, err)
	})

	t.Run("should inject anonymous sturct with primitive type", func(t *testing.T) {
		a := struct {
			TestObj *int `inject:""`
		}{}
		err := inject.IntoObject(&a)
		assert.Equal(t, nil, err)
		assert.NotEqual(t, nil, a.TestObj)
	})

	t.Run("should inject object through func", func(t *testing.T) {

		obj, err := inject.IntoFunc(func(user *FooUser) *fooService {
			assert.NotEqual(t, nil, user)
			return &fooService{
				FooUser: user,
			}
		})
		assert.Equal(t, nil, err)
		assert.NotEqual(t, nil, obj)
	})

	t.Run("should inject object through func", func(t *testing.T) {

		obj, err := inject.IntoFunc(func(user *FooUser) {
			assert.NotEqual(t, nil, user)
		})
		assert.Equal(t, nil, err)
		assert.Equal(t, nil, obj)
	})

	t.Run("should failed to inject object through func with empty interface", func(t *testing.T) {

		obj, err := inject.IntoFunc(func(user interface{}) *fooService {
			return &fooService{}
		})
		assert.NotEqual(t, nil, err)
		assert.Equal(t, nil, obj)
	})

	t.Run("should failed to inject object through nil func", func(t *testing.T) {
		obj, err := inject.IntoFunc(nil)
		assert.NotEqual(t, nil, err)
		assert.Equal(t, nil, obj)
	})
}
