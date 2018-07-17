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
	"testing"
	"github.com/hidevopsio/hiboot/pkg/utils"
	"path/filepath"
	"os"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/stretchr/testify/assert"
)

func init() {
	utils.ChangeWorkDir("../../")
}

func TestBuilderBuild(t *testing.T) {

	b := &Builder{
		Path:       filepath.Join(utils.GetWorkDir(), "config"),
		Name:       "application",
		FileType:   "yaml",
		Profile:    "local",
		ConfigType: Configuration{},
	}

	cp, err := b.Build()
	assert.Equal(t, nil, err)

	c := cp.(*Configuration)
	assert.Equal(t, "hiboot", c.App.Name)

	log.Print(c)
}

func TestFileDoesNotExist(t *testing.T) {

	b := &Builder{
		Path:       filepath.Join(utils.GetWorkDir(), "config"),
		Name:       "application",
		FileType:   "yaml",
		Profile:    "does-not-exist",
		ConfigType: Configuration{},
	}

	_, err := b.Build()
	assert.Contains(t, err.Error(), "not found")
}


func TestProfileIsEmpty(t *testing.T) {

	b := &Builder{
		Path:       filepath.Join(utils.GetWorkDir(), "config"),
		Name:       "application",
		FileType:   "yaml",
		Profile:    "",
		ConfigType: Configuration{},
	}

	_, err := b.Build()
	assert.Contains(t, err.Error(), "not found")
}



func TestWithoutReplacer(t *testing.T) {

	path := filepath.Join(utils.GetWorkDir(), "config")
	testProfile := "xxx"
	appConfig := "application"
	FileType := "yaml"
	testFile := appConfig + "-" + testProfile + ".yml"
	b := &Builder{
		Path:       path,
		Name:       appConfig,
		FileType:   FileType,
		Profile:    testProfile,
		ConfigType: Configuration{},
	}
	utils.CreateFile(path, testFile)
	_, err := b.Build()
	os.Remove(filepath.Join(path, testFile))
	assert.Equal(t, nil, err)
}


func TestBuilderInit(t *testing.T) {
	b := &Builder{
		Path:       filepath.Join(os.TempDir(), "config"),
		Name:       "foo",
		FileType:   "yaml",
		ConfigType: Configuration{},
	}

	err := b.Init()
	assert.Equal(t, nil, err)
}

func TestBuilderSave(t *testing.T) {
	b := &Builder{
		Path:       filepath.Join(os.TempDir(), "config"),
		Name:       "foo",
		FileType:   "yaml",
		ConfigType: Configuration{},
	}

	err := b.Init()
	assert.Equal(t, nil, err)

	c := &Configuration{
		App: App{
			Name: "foo",
			Project: "bar",
		},
		Server: Server{
			Port: 8080,
		},
	}
	err = b.Save(c)
	assert.Equal(t, nil, err)
}
