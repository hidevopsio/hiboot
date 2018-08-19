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

type Foo struct {
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
	fakeFile := "application-fake.yaml"
	os.Remove(filepath.Join(configPath, fakeFile))
	fakeContent :=
		"fake:\n" +
		"  name: foo\n" +
		"  nickname: ${app.name} ${fake.name}\n" +
		"  username: ${unknown.name:bar}\n"
	n, err := io.WriterFile(configPath, fakeFile, []byte(fakeContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, n, len(fakeContent))

	f := GetFactory()
	f.Build()
	fci := f.Configuration("fake")
	assert.NotEqual(t, nil, fci)
	fc := fci.(*FakeConfiguration)

	assert.Equal(t, "hiboot foo", fc.FakeProperties.Nickname)
	assert.Equal(t, "bar", fc.FakeProperties.Username)
	assert.Equal(t, "foo", fc.FakeProperties.Name)
	assert.Equal(t, "foo", f.Instances()["foo"].(*Foo).Name)
	assert.Equal(t, "foo", f.Instance("foo").(*Foo).Name)

	// get all configs
	cfs := f.Configurations()
	assert.Equal(t, 2, len(cfs))

}