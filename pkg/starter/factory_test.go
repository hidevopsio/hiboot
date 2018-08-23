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

package starter

import (
	"testing"
	"path/filepath"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
	"os"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
)

type FakeProperties struct {
	Name string
	Nickname string
	Username string
}

type FakeConfiguration struct {
	FakeProperties FakeProperties `mapstructure:"fake"`
}

type FooConfiguration struct {
	FakeProperties FakeProperties `mapstructure:"fake"`
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
	AddConfig("fake", FakeConfiguration{})
}

func (c *FakeConfiguration) Foo() *Foo {
	f := new(Foo)
	f.Name = c.FakeProperties.Name

	return f
}

func TestBuild(t *testing.T) {
	configPath := filepath.Join(io.GetWorkDir(), "config")
	fakeFile := "application-fake.yml"
	os.Remove(filepath.Join(configPath, fakeFile))
	fakeContent :=
		"fake:\n" +
		"  name: foo\n" +
		"  nickname: ${app.name} ${fake.name}\n" +
		"  username: ${unknown.name:bar}\n"
	n, err := io.WriterFile(configPath, fakeFile, []byte(fakeContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, n, len(fakeContent))

	AddConfig("foo", FooConfiguration{})
	factory := GetFactory()
	factory.Build()

	t.Run("should add instance to factory at runtime", func(t *testing.T) {
		fakeInstance := &struct{Name string}{Name: "fake"}
		factory.AddInstance("fakeInstance", fakeInstance)
		gotFakeInstance := factory.Instance("fakeInstance")
		assert.Equal(t, fakeInstance, gotFakeInstance)
	})

	f := factory.Configuration("fake")
	assert.NotEqual(t, nil, f)
	fc := f.(*FakeConfiguration)

	fooCfg := factory.Configuration("foo")
	assert.NotEqual(t, nil, fooCfg)

	assert.Equal(t, "hiboot fake", fc.FakeProperties.Nickname)
	assert.Equal(t, "bar", fc.FakeProperties.Username)
	assert.Equal(t, "fake", fc.FakeProperties.Name)
	foo, ok := factory.Instances().Get("foo")
	assert.Equal(t, true, ok)
	assert.Equal(t, "fake", foo.(*Foo).Name)
	assert.Equal(t, "fake", factory.Instance("foo").(*Foo).Name)

	// get all configs
	cfs := factory.Configurations()
	assert.Equal(t, 3, cfs.Count())

}

func TestAdd(t *testing.T) {
	fooBar := new(FooBar)
	Add(fooBar)
	factory := GetFactory()
	f := factory.Instance("fooBar")
	assert.Equal(t, f, fooBar)
}

func TestAddConfigs(t *testing.T) {
	preCfg := struct{Name string}{}
	AddPreConfig("pre", preCfg)
	postCfg := struct{Name string}{}
	AddPostConfig("post", postCfg)
}