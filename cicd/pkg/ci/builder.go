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

// TODO: app config should be generic kit

package ci

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/hidevopsio/hi/boot/pkg/utils"
	"github.com/imdario/mergo"
)

type Configuration struct {
	Pipeline Pipeline `mapstructure:"pipeline"`
}

type Builder struct {}

// TODO: we may eliminate this func, merge to system.builder instead
func (b *Builder) Build(name string) *Configuration {
	// TODO: should be refactored
	var conf Configuration

	wd := utils.GetWorkingDir("/pkg/ci/builder.go")

	viper.AddConfigPath(wd + "/config")
	viper.SetConfigName("pipeline")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Error config file: %s \n", err))
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(fmt.Errorf("Unable to decode Config: %s \n", err))
	}

	var confReplacer Configuration
	viper.SetConfigName("pipeline-" + name)
	viper.MergeInConfig()
	err = viper.Unmarshal(&confReplacer)
	if err != nil {
		panic(fmt.Errorf("Unable to decode Config: %s \n", err))
	}

	mergo.Merge(&conf, confReplacer, mergo.WithOverride)

	return &conf
}
