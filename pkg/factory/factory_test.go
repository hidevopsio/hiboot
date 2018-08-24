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

package factory

import (
	"os"
	"testing"
	"path/filepath"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"github.com/hidevopsio/hiboot/pkg/system"
)

type FakeProperties struct {
	Name string			`value:"foo"`
	Nickname string		`value:"foobar"`
	Username string		`value:"fb"`
	Org string			`value:"hidevopsio"`
	Profile string		`value:"${APP_PROFILES_ACTIVE}"`
}

type FakeConfiguration struct {
	FakeProperties FakeProperties `mapstructure:"fake"`
}

type FooProperties struct {
	Name string			`value:"${fake.name}"`
	Nickname string		`value:"foobar"`
	Username string		`value:"fb"`
}

type FooConfiguration struct {
	FakeProperties FooProperties `mapstructure:"foo"`
}

func (c *FooConfiguration ) HelloWorld() string {
	return "Hello world"
}

type Foo struct {
	Name string
}
type FooBar struct {
	Name string
}

func init() {
	log.SetLevel(log.DebugLevel)
	io.EnsureWorkDir("../../")
	//app.AddConfig("fake", FakeConfiguration{})
}

func (c *FakeConfiguration) Foo() *Foo {
	f := new(Foo)
	f.Name = c.FakeProperties.Name

	return f
}

func TestBuild(t *testing.T) {
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

	f := new(ConfigurableFactory)
	f.InstanceFactory = new(InstanceFactory)
	f.InstanceFactory.Init(cmap.New())
	f.Init(configContainers)

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
	container.Set("foo", FooConfiguration{})
	container.Set("fake", FakeConfiguration{})

	f.Build(container)

	t.Run("should add instance to factory at runtime", func(t *testing.T) {
		fakeInstance := &struct{Name string}{Name: "fake"}
		f.SetInstance("fakeInstance", fakeInstance)
		gotFakeInstance := f.GetInstance("fakeInstance")
		assert.Equal(t, fakeInstance, gotFakeInstance)
	})

	t.Run("should get foo configuration", func(t *testing.T) {
		fci, ok := f.configurations.Get("foo")
		assert.Equal(t, true, ok)
		assert.NotEqual(t, nil, fci)

		assert.Equal(t, "Hello world", f.GetInstance("helloWorld").(string))

		fc := fci.(*FooConfiguration)
		assert.Equal(t, "hiboot foo", fc.FakeProperties.Nickname)
		assert.Equal(t, "bar", fc.FakeProperties.Username)
		assert.Equal(t, "foo", fc.FakeProperties.Name)
	})
}
