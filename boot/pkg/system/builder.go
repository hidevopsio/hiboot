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
	"reflect"
)

type Builder struct {
	Path     string
	Name     string
	FileType string
	Profile  string
}

func (b *Builder) Build(c interface{}) (interface{}, error) {
	t := reflect.TypeOf(c)
	conf := reflect.New(t)
	confReplacer := reflect.New(t)
	err := b.Read(&conf)
	if err != nil {
		return nil, err
	}
	err = b.Read(&confReplacer)
	if err != nil {
		return nil, err
	}

	mergo.Merge(&conf, confReplacer, mergo.WithOverride)

	return &conf, nil
}

func (b *Builder) Read(conf interface{}) error {
	v := viper.New()

	v.AddConfigPath(b.Path)
	v.SetConfigName(b.Name)
	v.SetConfigType(b.FileType)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Error config file: %s \n", err))
	}
	err = v.Unmarshal(conf)
	if err != nil {
		panic(fmt.Errorf("Unable to decode Config: %s \n", err))
	}
	return err
}
