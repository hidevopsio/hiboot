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
	"github.com/stretchr/testify/assert"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/factory"
	"hidevops.io/hiboot/pkg/factory/autoconfigure"
	"hidevops.io/hiboot/pkg/factory/instantiate"
	"hidevops.io/hiboot/pkg/inject"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/utils/cmap"
	"hidevops.io/hiboot/pkg/utils/io"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

type User struct {
	Name string
	App  string
}

type fakeRepository struct {
}

type fakeProperties struct {
	DefVarSlice   []string `default:"${app.name}"`
	DefProfiles   []string `default:"${app.profiles.include}"`
	Name          string   `default:"should not inject this default value as it will inject by system.Builder"`
	Nickname      string   `default:"should not inject this default value as it will inject by system.Builder"`
	Username      string   `default:"should not inject this default value as it will inject by system.Builder"`
	Url           string   `default:"should not inject this default value as it will inject by system.Builder"`
	DefStrVal     string   `default:"this is default value"`
	DefIntPropVal int      `default:"${prop.value:123}"`
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

type fooProperties struct {
	Name string `default:"foo"`
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
	at.Tag `value:"fake"`
	inject.BaseTag
}

func newFakeTag() inject.Tag {
	return &fakeTag{}
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
	Properties fooProperties `mapstructure:"foo"`
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
	at.Path	`value:"/path/to/hiboot"`

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
	DefInt         int            `value:"${fake.value:123}"`
}

type PropTestUser struct {
	Name string
	App  string
}

type PropFooUser struct {
	Name string
}

type propTestService struct {
	PropFooUser  *PropFooUser  `inject:"name=foo"`
	PropTestUser *PropTestUser `inject:"name=${fake.name},app=${app.name}"`
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

type BxzService interface {
	GetNickname() string
}

type bazService struct {
	Name     string
	nickname string
}

func (s *bazService) GetNickname() string {
	return s.nickname
}

type buzService struct {
	Name     string
	nickname string
}

func (s *buzService) GetNickname() string {
	return s.nickname
}

type bxzServiceImpl struct {
	BxzService
	Name     string
	nickname string
}

func (s *bxzServiceImpl) GetNickname() string {
	return s.nickname
}

type dependencyInjectionTestService struct {
	BxzSvc     BxzService `inject:""`
	BazService BxzService `inject:""`
	BuzService BxzService `inject:""`
	BozService BxzService `inject:"buzService"`
	name       string
}

type buzzObj struct {
	Name string
}

type buzzService struct {
	bz *buzzObj
}

func newBuzzService(bz *buzzObj) *buzzService {
	return &buzzService{
		bz: bz,
	}
}

type CatInterface interface {
}

type animalService struct {
	Cat CatInterface `inject:""`
}

type testTag struct {
	at.Tag `value:"test"`
	inject.BaseTag
}

var (
	appName    = "hiboot"
	fakeName   = "fake"
	fooName    = "foo"
	fakeUrl    = "http://fake.com/api/foo"
	defaultUrl = "http://localhost:8080"

	cf factory.ConfigurableFactory
)

type Hello string
type HelloWorld string
type Foo struct{}
type Bar struct{}

type helloConfiguration struct {
	app.Configuration
}

func (c *helloConfiguration) HelloWorld(h Hello) HelloWorld {
	return HelloWorld(h + "World")
}

func (c *helloConfiguration) Bar(f *Foo) *Bar {
	return &Bar{}
}

type emptyInterface interface {
}

type nilConfiguration struct {
	app.Configuration
}

func (c *nilConfiguration) Bar(i emptyInterface) *Bar {
	return &Bar{}
}

type Member interface {
	SetRole(name string)
}

type member struct {
	role string
}

func newMember() Member {
	return &member{role: "admin"}
}

func (m *member) SetRole(name string) {
	m.role = name
}

type greeter struct {
	member Member
}

func newGreeter(m Member) *greeter {
	return &greeter{member: m}
}

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

func setUp(t *testing.T) factory.ConfigurableFactory {
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
	customProps := cmap.New()
	customProps.Set("app.project", "inject-test")
	cf = autoconfigure.NewConfigurableFactory(
		instantiate.NewInstantiateFactory(instances, []*factory.MetaData{}, customProps),
		configurations)
	cf.BuildSystemConfig()
	configs := []*factory.MetaData{
		factory.NewMetaData(new(fakeConfiguration)),
		factory.NewMetaData(new(fooConfiguration)),
	}
	cf.Build(configs)
	cf.BuildComponents()

	return cf
}

func TestInject(t *testing.T) {

	cf := setUp(t)
	fakeConfig := cf.Configuration("fake").(*fakeConfiguration)

	cf.SetInstance("inject_test.hello", Hello("Hello"))

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

	t.Run("should inject default int", func(t *testing.T) {
		assert.Equal(t, 123, fakeConfig.Properties.DefIntPropVal)
	})

	t.Run("should get config", func(t *testing.T) {
		fr := cf.GetInstance("inject_test.fakeRepository")
		assert.NotEqual(t, nil, fr)
	})

	injecting := inject.NewInject(cf)
	t.Run("should inject properties into sub struct", func(t *testing.T) {
		testObj := new(propTestService)
		err := injecting.IntoObject(testObj)
		assert.Equal(t, nil, err)
		assert.Equal(t, "foo", testObj.PropFooUser.Name)
	})

	t.Run("should inject through method", func(t *testing.T) {
		fu := cf.GetInstance("inject_test.fooUser")
		bu := cf.GetInstance("inject_test.user")
		fp := cf.GetInstance("inject_test.fakeRepository")
		u := new(User)

		cf.SetInstance("inject_test.barUser", u)

		svc, err := injecting.IntoFunc(newMethodInjectionService)
		assert.Equal(t, nil, err)
		s := svc.(*MethodInjectionService)
		assert.Equal(t, fu, s.fooUser)
		assert.Equal(t, bu, s.barUser)
		assert.Equal(t, fp, s.repository)
	})

	t.Run("should inject repository", func(t *testing.T) {
		us := new(userService)
		err := injecting.IntoObject(us)
		assert.Equal(t, nil, err)
		assert.NotEqual(t, (*User)(nil), us.User)
		assert.Equal(t, fooName, us.FooUser.Name)
		//assert.Equal(t, fakeName, us.FakeUser.Name)
		//assert.Equal(t, appName, us.FakeUser.App)
		assert.Equal(t, fakeUrl, us.Url)
		assert.Equal(t, defaultUrl, us.DefaultUrl)
		assert.NotEqual(t, (*fakeRepository)(nil), us.FakeRepository)
	})

	t.Run("should not inject unimplemented interface into FooBarRepository", func(t *testing.T) {
		fb := new(foobarRecursiveInject)
		err := injecting.IntoObject(fb)
		assert.Equal(t, nil, err)
	})

	t.Run("should not inject unimplemented interface into FooRepository", func(t *testing.T) {
		fs := new(fooService)
		err := injecting.IntoObject(fs)
		assert.Equal(t, "foo", fs.FooUser.Name)
		assert.Equal(t, nil, err)
	})

	t.Run("should not inject system property into object", func(t *testing.T) {
		fs := new(hibootService)
		err := injecting.IntoObject(fs)
		assert.Equal(t, nil, err)
		//assert.Equal(t, appName, fs.HibootUser.Name)
	})

	t.Run("should not inject unimplemented interface into BarRepository", func(t *testing.T) {
		bs := new(barService)
		err := injecting.IntoObject(bs)
		assert.Equal(t, nil, err)
	})

	t.Run("should inject recursively", func(t *testing.T) {
		ps := &recursiveInject{UserService: new(userService)}
		err := injecting.IntoObject(ps)
		assert.Equal(t, nil, err)
		assert.NotEqual(t, (*User)(nil), ps.UserService.User)
		assert.NotEqual(t, (*fakeRepository)(nil), ps.UserService.FakeRepository)
	})

	t.Run("should not inject slice", func(t *testing.T) {
		testSvc := struct {
			Users []FooUser `inject:""`
		}{}
		err := injecting.IntoObject(testSvc)
		assert.Equal(t, nil, err)
	})

	t.Run("should inject slice value", func(t *testing.T) {
		testSvc := new(sliceInjectionTestService)
		err := injecting.IntoObject(testSvc)
		assert.Equal(t, nil, err)
		assert.Equal(t, []string{"foo", "fake"}, testSvc.Profiles)
		assert.Equal(t, []string{"a", "b", "c", "d"}, testSvc.Options)
		log.Debug(testSvc.Profiles)
	})

	t.Run("should inject slice value", func(t *testing.T) {
		err := injecting.IntoObject((*MethodInjectionService)(nil))
		assert.Equal(t, inject.ErrInvalidObject, err)
	})

	t.Run("should inject slice value", func(t *testing.T) {
		err := injecting.IntoObject((*string)(nil))
		assert.Equal(t, inject.ErrInvalidObject, err)
	})

	t.Run("should ignore to inject with invalid struct type BazService", func(t *testing.T) {
		ts := new(dependencyInjectionTestService)
		err := injecting.IntoObject(ts)
		assert.Equal(t, nil, err)
	})

	t.Run("should failed to inject if the type of param and receiver are the same", func(t *testing.T) {
		buzz, err := injecting.IntoFunc(newBuzzService)
		assert.Equal(t, nil, err)
		assert.NotEqual(t, nil, buzz)
	})

	t.Run("should skip inject if the type of param is an unimplemented interface", func(t *testing.T) {
		catSvc := new(animalService)
		err := injecting.IntoObject(catSvc)
		assert.Equal(t, nil, err)
		assert.Equal(t, nil, catSvc.Cat)
	})

	t.Run("should inject value and it must not be singleton", func(t *testing.T) {
		a := &struct {
			TestName string `value:"test-data-from-a"`
		}{}
		b := &struct{ TestName string }{}
		err := injecting.IntoObject(a)
		assert.Equal(t, nil, err)
		injecting.IntoObject(b)
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
		err := injecting.IntoObject(a)
		assert.Equal(t, nil, err)
	})

	t.Run("should inject anonymous sturct with primitive type", func(t *testing.T) {
		a := struct {
			TestObj *int `inject:""`
		}{}
		err := injecting.IntoObject(&a)
		assert.Equal(t, nil, err)
		assert.NotEqual(t, nil, a.TestObj)
	})

	t.Run("should inject object through func", func(t *testing.T) {

		obj, err := injecting.IntoFunc(func(user *FooUser) *fooService {
			assert.NotEqual(t, nil, user)
			return &fooService{
				FooUser: user,
			}
		})
		assert.Equal(t, nil, err)
		assert.NotEqual(t, nil, obj)
	})

	t.Run("should inject object through func", func(t *testing.T) {

		obj, err := injecting.IntoFunc(func(user *FooUser) {
			assert.NotEqual(t, nil, user)
		})
		assert.Equal(t, nil, err)
		assert.Equal(t, nil, obj)
	})

	t.Run("should failed to inject object through func with empty interface", func(t *testing.T) {

		obj, err := injecting.IntoFunc(func(user interface{}) *fooService {
			return &fooService{}
		})
		assert.NotEqual(t, nil, err)
		assert.Equal(t, nil, obj)
	})

	t.Run("should failed to inject object through nil func", func(t *testing.T) {
		obj, err := injecting.IntoFunc(nil)
		assert.NotEqual(t, nil, err)
		assert.Equal(t, nil, obj)
	})

	t.Run("should failed to inject unmatched object type through func", func(t *testing.T) {
		res, err := injecting.IntoFunc(newGreeter)
		assert.NotEqual(t, nil, err)
		assert.Equal(t, nil, res)
	})

	t.Run("should inject into object by inject tag", func(t *testing.T) {
		cf.SetInstance(new(bazService))
		cf.SetInstance(new(buzService))
		cf.SetInstance(new(bxzServiceImpl))

		svc := new(dependencyInjectionTestService)
		err := injecting.IntoObject(svc)

		bazSvc := cf.GetInstance(bazService{})
		buzSvc := cf.GetInstance(buzService{})
		bxzSvc := cf.GetInstance(new(BxzService))

		assert.Equal(t, nil, err)
		assert.Equal(t, bazSvc, svc.BazService)
		assert.Equal(t, buzSvc, svc.BozService)
		assert.Equal(t, buzSvc, svc.BuzService)
		assert.Equal(t, bxzSvc, svc.BxzSvc)
	})

	t.Run("should report error when the dependency of the method is not found", func(t *testing.T) {
		conf := new(nilConfiguration)
		helloTyp := reflect.TypeOf(conf)
		numOfMethod := helloTyp.NumMethod()
		//log.Debug("methods: ", numOfMethod)
		for mi := 0; mi < numOfMethod; mi++ {
			method := helloTyp.Method(mi)
			_, err := injecting.IntoMethod(conf, method)
			assert.NotEqual(t, nil, err)
		}
	})

	t.Run("should report error when pass nil to injector that inject into method", func(t *testing.T) {

		_, err := injecting.IntoMethod(nil, nil)
		assert.NotEqual(t, nil, err)

	})

	t.Run("should inject into method", func(t *testing.T) {
		helloConfig := new(helloConfiguration)
		helloTyp := reflect.TypeOf(helloConfig)
		numOfMethod := helloTyp.NumMethod()
		//log.Debug("methods: ", numOfMethod)
		for mi := 0; mi < numOfMethod; mi++ {
			method := helloTyp.Method(mi)
			res, err := injecting.IntoMethod(helloConfig, method)
			assert.Equal(t, nil, err)
			assert.NotEqual(t, nil, res)
		}
	})

}

func TestInjectIntoFunc(t *testing.T) {
	cf := setUp(t)
	injecting := inject.NewInject(cf)
	m, err := injecting.IntoFunc(newMember)
	assert.Equal(t, nil, err)
	cf.SetInstance(m)
	t.Run("should failed to inject unmatched object type through func", func(t *testing.T) {
		res, err := injecting.IntoFunc(newGreeter)
		assert.Equal(t, nil, err)
		assert.NotEqual(t, nil, res)
	})
}

type RequestMapping string

func (m RequestMapping) Value(value string) at.StringAnnotation {
	return RequestMapping(value)
}
func (m RequestMapping) String()(value string) {
	return string(m)
}

type RequestPath struct{}

func (p *RequestPath) Value(value string) at.StringAnnotation {

	return nil
}

func (p *RequestPath) String()(value string) {
	return ""
}

type FakeUser struct {
	RequestMappingA    RequestMapping
	RequestMappingB    *RequestMapping
	RequestPathA RequestPath
	RequestPathB *RequestPath
}

func (FakeUser) Value(value string) at.StringAnnotation {return nil}
func (FakeUser) String()(value string)               {return ""}

func TestHasAnnotation(t *testing.T) {
	t.Run("should implements interface Model", func(t *testing.T) {
		m := &FakeUser{}

		s := reflect.ValueOf(m).Elem()
		typ := s.Type()
		modelType := reflect.TypeOf((*at.StringAnnotation)(nil)).Elem()

		expected := []bool {
			true,
			true,
			false,
			true,
		}

		for i := 0; i < s.NumField(); i++ {
			f := typ.Field(i)
			impl := f.Type.Implements(modelType)
			log.Debugf("%d: %s %s -> %t", i, f.Name, f.Type, impl)
			assert.Equal(t, expected[i], impl)
		}
	})
}