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

package system

import (
	"github.com/stretchr/testify/assert"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/utils/io"
	"os"
	"path/filepath"
	"testing"
)

type profiles struct {
	Include []string `json:"include"`
	Active  string   `json:"active"`
}

type app struct {
	Project        string   `json:"project"`
	Name           string   `json:"name"`
	Profiles       profiles `json:"profiles"`
	DataSourceType string   `json:"data_source_type"`
}

type server struct {
	Port int32 `json:"port"`
}

type logging struct {
	Level string `json:"level"`
}

func init() {
	io.EnsureWorkDir(1, "config/application.yml")
	log.SetLevel(log.DebugLevel)
	customProps["app.project"] = "system-test"
}

var customProps = make(map[string]interface{})

func TestBuilderBuild(t *testing.T) {
	testProj := "system-configuration-test"
	customProps["app.project"] = testProj
	b := NewBuilder(&Configuration{},
		filepath.Join(io.GetWorkDir(), "config"),
		"application",
		"yaml",
		"local",
		customProps)

	t.Run("should build configuration properly", func(t *testing.T) {
		cp, err := b.Build()
		assert.Equal(t, nil, err)
		c := cp.(*Configuration)
		assert.Equal(t, "hiboot", c.App.Name)
		assert.Equal(t, testProj, c.App.Project)
	})

	t.Run("should build configuration properly", func(t *testing.T) {
		cp, err := b.Build()
		assert.Equal(t, nil, err)
		c := cp.(*Configuration)
		assert.Equal(t, "hiboot", c.App.Name)
		log.Debugf("app: %v", b.Get("app"))
		log.Debugf("app.name: %v", b.Get("app.name"))
		log.Debugf("server: %v", b.Get("server"))
	})

}

func TestBuilderBuildWithError(t *testing.T) {

	b := &Builder{}

	_, err := b.Build()
	assert.Equal(t, nil, err)

}

func TestBuilderBuildWithProfile(t *testing.T) {
	customProps["app.project"] = "local-test"
	b := NewBuilder(&Configuration{},
		filepath.Join(io.GetWorkDir(), "config"),
		"application",
		"yaml",
		"local",
		customProps)

	cp, err := b.BuildWithProfile()
	assert.Equal(t, nil, err)

	c := cp.(*Configuration)
	assert.Equal(t, "8080", c.Server.Port)
	log.Print(c)

	b.Profile = ""
	_, err = b.BuildWithProfile()
	assert.Equal(t, nil, err)

}

func TestFileDoesNotExist(t *testing.T) {

	b := NewBuilder(&Configuration{},
		filepath.Join(io.GetWorkDir(), "config"),
		"application",
		"yaml",
		"does-not-exist",
		customProps)
	t.Run("use default profile if custom profile does not exist", func(t *testing.T) {
		_, err := b.Build()
		assert.Equal(t, nil, err)
	})
}

func TestWrongFileFormat(t *testing.T) {

	b := NewBuilder(&Configuration{},
		filepath.Join(os.TempDir(), "config"),
		"test",
		"yaml",
		"abc",
		customProps)
	io.CreateFile(os.TempDir(), "test-abc.yml")
	io.WriterFile(os.TempDir(), "test-abc.yml", []byte(": 1234"))
	t.Run("should report error: did not find expected key", func(t *testing.T) {
		_, err := b.Build()
		assert.Equal(t, nil, err)
	})
	io.WriterFile(os.TempDir(), "test-abc.yml", []byte("abc:"))
	t.Run("use default profile if custom profile does not exist", func(t *testing.T) {
		_, err := b.Build()
		assert.Equal(t, nil, err)
	})
}

func TestProfileIsEmpty(t *testing.T) {

	b := NewBuilder(Configuration{},
		filepath.Join(io.GetWorkDir(), "config"),
		"application",
		"yaml",
		"",
		customProps)

	t.Run("use default profile if custom profile does not exist", func(t *testing.T) {
		_, err := b.Build()
		assert.Equal(t, nil, err)
	})
}

func TestWithoutReplacer(t *testing.T) {

	path := filepath.Join(io.GetWorkDir(), "config")
	testProfile := "xxx"
	appConfig := "application"
	testFile := appConfig + "-" + testProfile + ".yml"
	b := NewBuilder(&Configuration{},
		path,
		"application",
		"yaml",
		"xxx",
		customProps)
	io.CreateFile(path, testFile)
	_, err := b.Build()
	os.Remove(filepath.Join(path, testFile))
	assert.Equal(t, nil, err)

}

func TestBuilderInit(t *testing.T) {
	b := NewBuilder(&Configuration{},
		filepath.Join(os.TempDir(), "config"),
		"foo",
		"yaml",
		"",
		customProps)

	err := b.Init()
	assert.Equal(t, nil, err)
}

func TestBuilderSave(t *testing.T) {
	b := NewBuilder(&Configuration{},
		filepath.Join(os.TempDir(), "config"),
		"foo",
		"yaml",
		"",
		customProps)

	err := b.Init()
	assert.Equal(t, nil, err)

	c := &Configuration{
		App: App{
			Name:    "foo",
			Project: "bar",
		},
		Server: Server{
			Port: "8080",
		},
	}

	t.Run("should save struct to file", func(t *testing.T) {
		err = b.Save(c)
		assert.Equal(t, nil, err)
	})

	t.Run("should save struct to file", func(t *testing.T) {
		err = b.Save("wrong-format")
		assert.Contains(t, err.Error(), "wrong")
	})
}
