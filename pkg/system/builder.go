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
	"path/filepath"
	"gopkg.in/yaml.v2"
	"bytes"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
)

type Builder struct {
	Path       string
	Name       string
	FileType   string
	Profile    string
	ConfigType interface{}
}

// create new viper instance
func (b *Builder) New(name string) *viper.Viper {
	v := viper.New()
	v.AddConfigPath(b.Path)
	v.SetConfigName(name)
	v.SetConfigType(b.FileType)
	return v
}

// create file if it's not exist
func (b *Builder) Init() (error) {
	return io.CreateFile(b.Path, b.Name + "." + b.FileType)
}

func (b *Builder) isFileNotExist(path string) bool {
	for _, ext := range viper.SupportedExts {
		if !io.IsPathNotExist(path + ext) {
			return false
		}
	}
	return true
}

// build config file
func (b *Builder) Build() (interface{}, error) {

	conf, err := b.Read(b.Name)
	if err != nil {
		return nil, err
	}

	name := b.Name + "-" + b.Profile
	// allow the empty of the profile
	if b.Profile == "" || b.isFileNotExist(filepath.Join(b.Path, name) + ".") {
		return conf, nil
	}

	confReplacer, err := b.Read(name)
	if err != nil {
		return conf, err
	}

	mergo.Merge(conf, confReplacer, mergo.WithOverride, mergo.WithAppendSlice)

	return conf, nil
}


// build config file
func (b *Builder) BuildWithProfile() (interface{}, error) {
	name := b.Name + "-" + b.Profile
	// allow the empty of the profile
	if b.Profile == "" || b.isFileNotExist(filepath.Join(b.Path, name) + ".") {
		return reflector.NewReflectType(b.ConfigType), nil
	}

	conf, err := b.Read(name)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

// Read single file
func (b *Builder) Read(name string) (interface{}, error) {

	v := b.New(name)
	err := v.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error on config file: %s", err)
	}
	st := b.ConfigType

	cp := reflector.NewReflectType(st)

	err = v.Unmarshal(cp)
	if err != nil {
		return nil, fmt.Errorf("error on viper config unmarshal : %s", err)
	}
	return cp, err
}


// Save configurations to file
func (b *Builder) Save(p interface{}) (error) {

	v := b.New(b.Name)

	y, err := yaml.Marshal(p)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	err = v.ReadConfig(bytes.NewBuffer(y))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	return v.WriteConfig()
}

