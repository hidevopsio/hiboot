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

package system

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/imdario/mergo"
	"github.com/hidevopsio/hi/boot/pkg/utils"
)

type Builder struct {
	Path       string
	Name       string
	FileType   string
	Profile    string
	ConfigType interface{}
}

func (b *Builder) Build() (interface{}, error) {

	conf, err := b.Read(false)
	if err != nil {
		return nil, err
	}

	if b.Profile == "" {
		return conf, nil
	}

	confReplacer, err := b.Read(true)
	if err != nil {
		return nil, err
	}

	mergo.Merge(conf, confReplacer, mergo.WithOverride)

	return conf, nil
}

func (b *Builder) Read(override bool) (interface{}, error) {

	name := b.Name
	if override {
		name = b.Name + "-" + b.Profile
	}

	v := viper.New()
	v.AddConfigPath(b.Path)
	v.SetConfigName(name)
	v.SetConfigType(b.FileType)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Error config file: %s \n", err))
	}
	st := b.ConfigType

	cp := utils.NewReflectType(st)

	err = v.Unmarshal(cp)
	if err != nil {
		panic(fmt.Errorf("Unable to decode Config: %s \n", err))
	}
	return cp, err
}
