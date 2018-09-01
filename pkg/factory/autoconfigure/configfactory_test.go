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
	"os"
	"testing"
	"path/filepath"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"github.com/hidevopsio/hiboot/pkg/system"
	"github.com/hidevopsio/hiboot/pkg/factory/instantiate"
	"github.com/hidevopsio/hiboot/pkg/factory/autoconfigure"
)

type FakeProperties struct {
	Name string			`default:"foo"`
	Nickname string		`default:"foobar"`
	Username string		`default:"fb"`
	Org string			`default:"hidevopsio"`
	Profile string		`default:"${APP_PROFILES_ACTIVE}"`
}

type FakeConfiguration struct {
	FakeProperties FakeProperties `mapstructure:"fake"`
}

type FooProperties struct {
	Name string			`default:"${fake.name}"`
	Nickname string		`default:"foobar"`
	Username string		`default:"fb"`
}

type FooConfiguration struct {
	FakeProperties FooProperties `mapstructure:"foo"`
}

func (c *FooConfiguration ) HelloWorld() string {
	return "Hello world"
}

type Foo struct {
	Name string
	Bar *Bar
}

type Bar struct {
	Name string
}

type FooBar struct {
	Name string
}

type barConfiguration struct {
	FakeProperties FooProperties `mapstructure:"bar"`
}

func (c *barConfiguration) BarBar() *Bar {
	return new(Bar)
}

func init() {
	log.SetLevel(log.DebugLevel)
	io.EnsureWorkDir(1, "config/application.yml")
}

func (c *FakeConfiguration) Foo(bar *Bar) *Foo {
	f := new(Foo)
	f.Name = c.FakeProperties.Name
	f.Bar = bar
	return f
}

func (c *FakeConfiguration) FooBar(foo *Foo) *FooBar {
	f := new(FooBar)
	f.Name = foo.Name

	return f
}

func (c *FakeConfiguration) Bar() *Bar {
	b := new(Bar)
	b.Name = "bar"

	return b
}


func TestConfigurableFactory(t *testing.T) {
	configPath := filepath.Join(os.TempDir(), "config")

	fakeFile := "application.yml"
	os.Remove(filepath.Join(configPath, fakeFile))
	fakeContent :=
		"app:\n" +
		"  project: hidevopsio\n" +
		"  name: hiboot\n" +
		"  version: ${unknown.version:0.0.1}\n"
	n, err := io.WriterFile(configPath, fakeFile, []byte(fakeContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, n, len(fakeContent))

	configContainers := cmap.New()

	f := new(autoconfigure.ConfigurableFactory)

	t.Run("should check if factory is FactoryCannotBeNilError", func(t *testing.T) {
		err := f.Initialize(nil)
		assert.Equal(t, autoconfigure.FactoryCannotBeNilError, err)
	})

	f.InstantiateFactory = new(instantiate.InstantiateFactory)
	t.Run("should check if factory is FactoryIsNotInitializedError", func(t *testing.T) {
		err := f.Initialize(nil)
		assert.Equal(t, autoconfigure.FactoryIsNotInitializedError, err)
	})

	f.InstantiateFactory.Initialize(cmap.New())
	f.Initialize(configContainers)

	t.Run("should build app config", func(t *testing.T) {
		io.ChangeWorkDir(os.TempDir())
		err := f.BuildSystemConfig(system.Configuration{})
		assert.Equal(t, nil, err)
	})

	fooFile := "application-foo.yml"
	os.Remove(filepath.Join(configPath, fooFile))
	fooContent :=
		"foo:\n" +
			"  name: foo\n" +
			"  nickname: ${app.name} ${foo.name}\n" +
			"  username: ${unknown.name:bar}\n"
	_, err = io.WriterFile(configPath, fooFile, []byte(fooContent))
	assert.Equal(t, nil, err)

	container := cmap.New()
	fooConfig := new(FooConfiguration)
	container.Set("foo", fooConfig)
	fakeCfg := new(FakeConfiguration)
	container.Set("fake", fakeCfg)
	container.Set("fakeFake", FakeConfiguration{})

	f.Build(container)

	t.Run("should instantiate by name", func(t *testing.T) {
		bc := new(barConfiguration)
		_, err := f.InstantiateByName(bc, "Bar")
		assert.Equal(t, autoconfigure.InvalidMethodError, err)
	})

	t.Run("should instantiate by name", func(t *testing.T) {
		bc := new(barConfiguration)
		bb, err := f.InstantiateByName(bc, "BarBar")
		assert.Equal(t, nil, err)
		assert.NotEqual(t, nil, bb)
	})

	t.Run("should get SystemConfiguration", func(t *testing.T) {
		sysCfg := f.SystemConfiguration()
		assert.NotEqual(t, nil, sysCfg)
	})

	t.Run("should get fake configuration", func(t *testing.T) {
		fakeCfg := f.Configuration("fake")
		assert.NotEqual(t, nil, fakeCfg)
	})

	t.Run("should not get non-existence configuration", func(t *testing.T) {
		fakeCfg := f.Configuration("non-existence")
		assert.Equal(t, nil, fakeCfg)
	})

	t.Run("should add instance to factory at runtime", func(t *testing.T) {
		fakeInstance := &struct{Name string}{Name: "fake"}
		f.SetInstance("fakeInstance", fakeInstance)
		gotFakeInstance := f.GetInstance("fakeInstance")
		assert.Equal(t, fakeInstance, gotFakeInstance)
	})

	t.Run("should get foo configuration", func(t *testing.T) {
		assert.Equal(t, "Hello world", f.GetInstance("helloWorld").(string))

		assert.Equal(t, "hiboot foo", fooConfig.FakeProperties.Nickname)
		assert.Equal(t, "bar", fooConfig.FakeProperties.Username)
		assert.Equal(t, "foo", fooConfig.FakeProperties.Name)
	})

}
