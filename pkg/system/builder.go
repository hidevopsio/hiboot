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
	"bytes"
	"fmt"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"path/filepath"
	"reflect"
)

// Builder is the config file (yaml, json) builder
type Builder struct {
	*viper.Viper
	Path       string
	Name       string
	FileType   string
	Profile    string
	ConfigType interface{}
}

func NewBuilder(configType interface{}, path, name, fileType, profile string) *Builder {
	return &Builder{
		Path:       path,
		Name:       name,
		FileType:   fileType,
		Profile:    profile,
		ConfigType: configType,
	}
}

// New create new viper instance
func (b *Builder) New(name string) {
	b.Viper = viper.New()
	b.AddConfigPath(b.Path)
	b.SetConfigName(name)
	b.SetConfigType(b.FileType)
}

// Init create file if it's not exist
func (b *Builder) Init() error {
	return io.CreateFile(b.Path, b.Name+"."+b.FileType)
}

func (b *Builder) isFileNotExist(path string) bool {
	for _, ext := range viper.SupportedExts {
		if !io.IsPathNotExist(path + ext) {
			return false
		}
	}
	return true
}

// Build config file
func (b *Builder) Build(profiles ...string) (interface{}, error) {

	conf, err := b.Read(b.Name)
	if err != nil {
		//log.Errorf("failed to read: %v", b.Name)
		return conf, nil
	}

	if len(profiles) == 0 && b.Profile != "" {
		profiles = append(profiles, b.Profile)
	}

	for _, profile := range profiles {
		name := b.Name + "-" + profile
		// allow the empty of the profile
		configFile := filepath.Join(b.Path, name)
		if profile == "" || b.isFileNotExist(configFile+".") {
			//log.Debugf("config file: %v does not exist", configFile)
			return conf, nil
		}

		_, err = b.Read(name)
	}

	return conf, err
}

// BuildWithProfile build config file
func (b *Builder) BuildWithProfile() (interface{}, error) {
	name := b.Name + "-" + b.Profile
	// allow the empty of the profile
	if b.Profile == "" || b.isFileNotExist(filepath.Join(b.Path, name)+".") {
		return reflector.NewReflectType(b.ConfigType), nil
	}

	return b.Read(name)
}

// Read single file
func (b *Builder) Read(name string) (interface{}, error) {
	b.New(name)
	// create new instance
	st := b.ConfigType
	val := reflect.ValueOf(st)
	//log.Debugf("value of configuration: %v, kind: %v", val, val.Kind())
	cp := b.ConfigType
	if val.Kind() == reflect.Struct {
		cp = reflector.NewReflectType(st)
	}
	// read config
	err := b.ReadInConfig()
	if err != nil {
		return cp, fmt.Errorf("error on config file: %s", err)
	}

	err = b.Unmarshal(cp)
	return cp, err
}

// Save configurations to file
func (b *Builder) Save(p interface{}) error {

	b.New(b.Name)

	y, err := yaml.Marshal(p)
	if err == nil {
		err = b.ReadConfig(bytes.NewBuffer(y))
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return err
		}
	}

	return b.WriteConfig()
}
