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

package ci

import (
	"testing"
	"github.com/hidevopsio/hi/boot/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hi/boot/pkg/system"
	"github.com/hidevopsio/hi/boot/pkg/utils"
	"github.com/imdario/mergo"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestPipelineBuilder(t *testing.T) {

	log.Debug("Test Pipeline Builder")

	sysCfg := system.Build()
	log.Debug(sysCfg)
	assert.Equal(t, "hi", sysCfg.App.Name)
	builder := &Builder{}
	cfg := builder.Build("java")
	log.Debug(cfg)
	assert.Equal(t, "java", cfg.Pipeline.Name)

	// replace internal variables
	utils.Replace(&cfg.Pipeline, "profile", "dev")
	log.Debug(cfg.Pipeline)
}

type Foo struct {
	Name        string
	Profile 	string
	Env         []system.Env
	ConfigFiles []string
}

type Bar struct {
	Name        string
	Age         int
	Env         []system.Env
	ConfigFiles []string
}

func TestCopy(t *testing.T) {
	t1 := &Foo{
		Name: "foo",
		Env: []system.Env{
			{
				Name:  "ENV_FOO",
				Value: "env.foo",
			},
			{
				Name:  "ENV_Bar",
				Value: "env.bar",
			},
		},
	}

	t2 := &Bar{}
	log.Println(t2)
	utils.Copy(t2, t1)
	log.Println(t2)

	assert.Equal(t, t1.Name, t2.Name)
	assert.Equal(t, t1.Env[0].Name, t2.Env[0].Name)
	assert.Equal(t, t1.Env[0].Value, t2.Env[0].Value)
}

func TestReplacement(t *testing.T) {
	t1 := &Foo{
		Name: "${foo}",
		Env: []system.Env{
			{
				Name:  "ENV_BAR",
				Value: "is_${foo}",
			},
			{
				Name:  "CONFIG",
				Value: "src/main/resources/application-${foo}.yml",
			},
		},
		ConfigFiles: []string{
			"src/main/resources/application.yml",
			"src/main/resources/application-${foo}.yml",
		},
	}

	log.Println(t1)
	utils.Replace(t1, "foo", "bar")
	log.Println(t1)

	assert.Equal(t, "bar", t1.Name)
	assert.Equal(t, "is_bar", t1.Env[0].Value)
}

func TestMerging(t *testing.T) {
	t1 := &Foo{
		Name: "foo",
		Env: []system.Env{
			{
				Name:  "ENV_FOO",
				Value: "is_foo",
			},
		},
		ConfigFiles: []string{
			"src/main/resources/application.yml",
		},
	}

	t2 := &Foo{
		Name: "bar",
		Profile: "dev",
		Env: []system.Env{
			{
				Name:  "CONFIG",
				Value: "src/main/resources/application-bar.yml",
			},
		},
		ConfigFiles: []string{
			"src/main/resources/application-bar.yml",
		},
	}

	log.Println(t1)
	assert.Equal(t, "foo", t1.Name)
	assert.Equal(t, "", t1.Profile)
	mergo.Merge(t1, t2, mergo.WithOverride)
	log.Println(t1)
	assert.Equal(t, "bar", t1.Name)
	assert.Equal(t, "dev", t1.Profile)
	assert.Equal(t, "is_foo", t1.Env[0].Value)
}
