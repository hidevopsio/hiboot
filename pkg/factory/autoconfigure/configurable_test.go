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

package autoconfigure_test

import (
	"github.com/stretchr/testify/assert"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/app/web/context"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/factory"
	"hidevops.io/hiboot/pkg/factory/autoconfigure"
	"hidevops.io/hiboot/pkg/factory/instantiate"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/utils/cmap"
	"hidevops.io/hiboot/pkg/utils/io"
	"os"
	"path/filepath"
	"testing"
)

type FakeProperties struct {
	Name     string `default:"foo"`
	Nickname string `default:"foobar"`
	Username string `default:"fb"`
	Org      string `default:"hidevopsio"`
	Profile  string `default:"${APP_PROFILES_ACTIVE:dev}"`
}

type Connection struct {
	context context.Context
}

type ContextAwareConfiguration struct {
	at.AutoConfiguration
	at.ContextAware
}

func newContextAwareConfiguration() *ContextAwareConfiguration {
	return &ContextAwareConfiguration{}
}

func (c *ContextAwareConfiguration) Connection(context context.Context) *Connection {
	return &Connection{context: context}
}

type unknownConfiguration struct {
	FakeProperties FakeProperties `mapstructure:"fake"`
}

type FooInterface interface {
}

type unsupportedConfiguration struct {
	FooInterface
	FakeProperties FakeProperties `mapstructure:"fake"`
}

type Properties struct {
	Name     string `default:"${fake.name}"`
	Nickname string `default:"foobar"`
	Username string `default:"fb"`
}

type Configuration struct {
	at.AutoConfiguration
}

type emptyConfiguration struct {
}

type FooConfiguration struct {
	at.AutoConfiguration
	FakeProperties Properties `mapstructure:"foo"`
}

type BarConfiguration struct {
	app.Configuration
	FakeProperties Properties `mapstructure:"bar"`
}

type FooBarConfiguration struct {
	app.Configuration
	FakeProperties Properties `mapstructure:"foobar"`
	foobar         *FooBar
}

type mercury struct {
	jupiter *jupiter
}

type mars struct {
	mercury *mercury
}

type jupiter struct {
}

type marsConfiguration struct {
	app.Configuration `depends:"mercuryConfiguration"`
}

func (c *marsConfiguration) Mars(m *mercury) *mars {
	return &mars{mercury: m}
}

type mercuryConfiguration struct {
	app.Configuration `depends:"jupiterConfiguration"`
}

func (c *mercuryConfiguration) Mercury(j *jupiter) *mercury {
	return &mercury{jupiter: j}
}

type jupiterConfiguration struct {
	app.Configuration
}

func (c *jupiterConfiguration) Jupiter() *jupiter {
	return new(jupiter)
}

func newFooBarConfiguration(foobar *FooBar) *FooBarConfiguration {
	return &FooBarConfiguration{
		foobar: foobar,
	}
}

type HelloWorld string
type Hello string

func (c *FooConfiguration) HelloWorld(foo Hello) HelloWorld {
	return HelloWorld(foo + " world")
}

func (c *FooConfiguration) Hello() Hello {
	return Hello("Hello")
}

type Foo struct {
	Name string
	Bar  *Bar
}

type Bar struct {
	Name string
}

type FooBar struct {
	Name string
}

type barConfiguration struct {
	FakeProperties Properties `mapstructure:"bar"`
}

func (c *barConfiguration) BarBar() *Bar {
	return new(Bar)
}

func init() {
	log.SetLevel(log.DebugLevel)
	io.EnsureWorkDir(1, "config/application.yml")
}

type fakeConfiguration struct {
	app.Configuration
	FakeProperties FakeProperties `mapstructure:"fake"`
}

func newFakeConfiguration() *fakeConfiguration {
	return &fakeConfiguration{}
}

type foobarConfiguration struct {
	app.Configuration
	FakeProperties FakeProperties `mapstructure:"foobar"`
}

func newFoobarConfiguration() *foobarConfiguration {
	return &foobarConfiguration{}
}

func (c *foobarConfiguration) Foo(bar *Bar) *Foo {
	f := new(Foo)
	f.Name = c.FakeProperties.Name
	f.Bar = bar
	return f
}

func (c *foobarConfiguration) FooBar(foo *Foo) *FooBar {
	f := new(FooBar)
	f.Name = foo.Name

	return f
}

func (c *foobarConfiguration) Bar() *Bar {
	b := new(Bar)
	b.Name = "bar"

	return b
}

type EarthConfiguration struct {
	app.Configuration
	instantiateFactory factory.InstantiateFactory
}

func newEarthConfiguration(instantiateFactory factory.InstantiateFactory) *EarthConfiguration {
	c := &EarthConfiguration{
		instantiateFactory: instantiateFactory,
	}
	c.RuntimeDeps.Set(c.RuntimeTree, []string{"autoconfigure_test.leaf", "autoconfigure_test.branch"})
	return c
}

type Land struct {
	Mountain *Mountain
}

type RuntimeTree struct {
	leaf   *Leaf
	branch *Branch
}

type Tree interface {
	Family() string
}

type oakTree struct {
	Tree
}

func (t *oakTree) Family() string {
	return "Beech"
}

type mapleTree struct {
}

func (t *mapleTree) Family() string {
	return "Soapberry"
}

type Branch struct {
}

type Leaf struct {
}

type Mountain struct {
	oakTree Tree
}

func (c *EarthConfiguration) Land(mountain *Mountain) *Land {
	return &Land{Mountain: mountain}
}

func (c *EarthConfiguration) Mountain(tree Tree) *Mountain {
	return &Mountain{oakTree: tree}
}

func (c *EarthConfiguration) Branch() *Branch {
	return &Branch{}
}

func (c *EarthConfiguration) Leaf() *Leaf {
	return &Leaf{}
}

func (c *EarthConfiguration) RuntimeTree() *RuntimeTree {
	leaf := c.instantiateFactory.GetInstance("autoconfigure_test.leaf").(*Leaf)
	branch := c.instantiateFactory.GetInstance("autoconfigure_test.branch").(*Branch)
	return &RuntimeTree{leaf: leaf, branch: branch}
}

func (c *EarthConfiguration) OakTree() Tree {
	return &oakTree{}
}

func (c *EarthConfiguration) MapleTree() Tree {
	return &mapleTree{}
}

type helloService struct{ foo *Foo }

func newHelloService(foo *Foo) *helloService {
	return &helloService{foo: foo}
}

func setFactory(t *testing.T, customProperties cmap.ConcurrentMap) factory.ConfigurableFactory {
	io.ChangeWorkDir(os.TempDir())

	configPath := filepath.Join(os.TempDir(), "config")

	fakeFile := "application.yml"
	os.Remove(filepath.Join(configPath, fakeFile))
	fakeContent :=
		"app:\n" +
			"  project: hidevopsio\n" +
			"  name: hiboot-test\n" +
			"  version: 0.0.1\n"
	n, err := io.WriterFile(configPath, fakeFile, []byte(fakeContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, n, len(fakeContent))

	fooFile := "application-foo.yml"
	os.Remove(filepath.Join(configPath, fooFile))
	fooContent :=
		"foo:\n" +
			"  name: foo\n" +
			"  nickname: ${app.name} ${foo.name}\n" +
			"  username: ${unknown.name:bar}\n"
	_, err = io.WriterFile(configPath, fooFile, []byte(fooContent))
	assert.Equal(t, nil, err)

	configContainers := cmap.New()
	customProperties.Set("logging.level", "debug")
	f := autoconfigure.NewConfigurableFactory(
		instantiate.NewInstantiateFactory(cmap.New(), []*factory.MetaData{}, customProperties),
		configContainers,
	)

	return f
}

func TestConfigurableFactory(t *testing.T) {
	customProperties := cmap.New()
	f := setFactory(t, customProperties)

	var err error
	t.Run("should set factory instance", func(t *testing.T) {
		err = f.SetInstance(factory.InstantiateFactoryName, f)
		assert.Equal(t, nil, err)

		f.SetInstance(factory.ConfigurableFactoryName, f)
		assert.Equal(t, nil, err)
	})
	//inject.SetFactory(f)
	// backup profile
	profile := os.Getenv(autoconfigure.EnvAppProfilesActive)
	t.Run("should build app config", func(t *testing.T) {
		os.Setenv(autoconfigure.EnvAppProfilesActive, "")
		sc, err := f.BuildSystemConfig()
		assert.Equal(t, nil, err)
		assert.Equal(t, "default", sc.App.Profiles.Active)
		// restore profile
		os.Setenv(autoconfigure.EnvAppProfilesActive, profile)
	})

	os.Setenv(autoconfigure.EnvAppProfilesActive, profile)
	customProperties.Set(autoconfigure.PropAppProfilesActive, "dev")
	t.Run("should build app config", func(t *testing.T) {
		sc, err := f.BuildSystemConfig()
		assert.Equal(t, nil, err)
		assert.Equal(t, "dev", sc.App.Profiles.Active)
		// restore profile
		os.Setenv(autoconfigure.EnvAppProfilesActive, profile)
	})

	t.Run("should build app config", func(t *testing.T) {
		_, err = f.BuildSystemConfig()
		assert.Equal(t, nil, err)
	})

	fooConfig := new(FooConfiguration)

	f.Build([]*factory.MetaData{
		factory.NewMetaData(new(emptyConfiguration)),
		factory.NewMetaData(newFakeConfiguration),
		factory.NewMetaData("foo", fooConfig),
		factory.NewMetaData(new(Configuration)),
		factory.NewMetaData(new(BarConfiguration)),
		factory.NewMetaData(new(marsConfiguration)),
		factory.NewMetaData(new(jupiterConfiguration)),
		factory.NewMetaData(new(mercuryConfiguration)),
		factory.NewMetaData(new(unknownConfiguration)),
		factory.NewMetaData(new(unsupportedConfiguration)),
		factory.NewMetaData(newFoobarConfiguration),
		factory.NewMetaData(newEarthConfiguration),
		factory.NewMetaData(newContextAwareConfiguration),
	})

	f.AppendComponent(newHelloService)
	ctx := web.NewContext(nil)
	f.AppendComponent("context.context", ctx)

	err = f.BuildComponents()

	t.Run("should get SystemConfiguration", func(t *testing.T) {
		sysCfg := f.SystemConfiguration()
		assert.NotEqual(t, nil, sysCfg)
	})

	t.Run("should get fake configuration", func(t *testing.T) {
		fake := f.Configuration("fake")
		assert.NotEqual(t, nil, fake)
	})

	t.Run("should not get non-existence configuration", func(t *testing.T) {
		nonExist := f.Configuration("non-existence")
		assert.Equal(t, nil, nonExist)
	})

	t.Run("should add instance to factory at runtime", func(t *testing.T) {
		fakeInstance := &struct{ Name string }{Name: "fake"}
		f.SetInstance("autoconfigure_test.fakeInstance", fakeInstance)
		gotFakeInstance := f.GetInstance("autoconfigure_test.fakeInstance")
		assert.Equal(t, fakeInstance, gotFakeInstance)
	})

	t.Run("should get foo configuration", func(t *testing.T) {
		helloWorld := f.GetInstance("autoconfigure_test.helloWorld")
		assert.NotEqual(t, nil, helloWorld)
		assert.Equal(t, HelloWorld("Hello world"), helloWorld)
	})

	t.Run("should get runtime created instances", func(t *testing.T) {
		runtimeTree := f.GetInstance(RuntimeTree{})
		leaf := f.GetInstance(Leaf{})
		branch := f.GetInstance(Branch{})
		assert.NotEqual(t, nil, runtimeTree)
		rt := runtimeTree.(*RuntimeTree)
		assert.Equal(t, branch, rt.branch)
		assert.Equal(t, leaf, rt.leaf)
	})

	t.Run("should report error on nil configuration", func(t *testing.T) {

	})
}

func TestReplacer(t *testing.T) {
	customProperties := cmap.New()
	customProperties.Set("app.profiles.filter", true)
	customProperties.Set("app.profiles.include", []string{"foo", "fake"})
	f := setFactory(t, customProperties)
	var err error

	t.Run("should build app config", func(t *testing.T) {
		_, err = f.BuildSystemConfig()
		assert.Equal(t, nil, err)
	})
	fooConfig := new(FooConfiguration)

	type outConfiguration struct {
		at.AutoConfiguration
		Properties Properties `mapstructure:"out"`
	}

	f.Build([]*factory.MetaData{
		factory.NewMetaData(new(outConfiguration)),
		factory.NewMetaData(newFakeConfiguration),
		factory.NewMetaData("foo", fooConfig),
	})

	t.Run("should get foo configuration", func(t *testing.T) {
		assert.Equal(t, "hiboot-test foo", fooConfig.FakeProperties.Nickname)
		assert.Equal(t, "bar", fooConfig.FakeProperties.Username)
		assert.Equal(t, "foo", fooConfig.FakeProperties.Name)
	})
}
