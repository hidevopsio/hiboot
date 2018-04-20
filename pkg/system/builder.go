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
	"github.com/hidevopsio/hiboot/pkg/utils"
	"path/filepath"
	"os"
	yaml "gopkg.in/yaml.v2"
	"bytes"
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

	_, err := os.Stat(b.Path)
	if os.IsNotExist(err) {
		err = os.Mkdir(b.Path, os.ModePerm)
	}
	if err != nil {
		return err
	}

	fn := filepath.Join(b.Path, b.Name) + "." + b.FileType
	_, err = os.Stat(fn)
	if os.IsNotExist(err) {
		f, err := os.OpenFile(fn, os.O_RDONLY | os.O_CREATE, 0666)
		f.Close()
		return err
	}

	return err
}

// build config file
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

	mergo.Merge(conf, confReplacer, mergo.WithOverride, mergo.WithAppendSlice)

	return conf, nil
}

// Read single file
func (b *Builder) Read(override bool) (interface{}, error) {

	name := b.Name
	if override {
		name = b.Name + "-" + b.Profile
	}

	v := b.New(name)
	err := v.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("Error config file: %s \n", err)
	}
	st := b.ConfigType

	cp := utils.NewReflectType(st)

	err = v.Unmarshal(cp)
	if err != nil {
		return nil, fmt.Errorf("Error on viper config unmarshal : %s \n", err)
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

