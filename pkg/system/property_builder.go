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
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/inject/annotation"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/mapstruct"
	"github.com/hidevopsio/hiboot/pkg/utils/replacer"
	"github.com/hidevopsio/hiboot/pkg/utils/sort"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
	"github.com/hidevopsio/mapstructure"
	"github.com/hidevopsio/viper"
)

const (
	Config             = "app.config"
	ConfigDir          = "app.config.dir"
	appProfilesInclude = "app.profiles.include"
)

type ConfigFile struct {
	fd       fs.File
	path     string
	name     string
	fileType string
	profile  string
}

type propertyBuilder struct {
	at.Qualifier `value:"github.com/hidevopsio/hiboot/pkg/system.builder"`
	*viper.Viper
	ConfigFile
	configuration     interface{}
	defaultProperties map[string]interface{}
	profiles          []string
	merge             bool
	embedFS           *embed.FS
	sync.Mutex
}

// NewBuilder is the constructor of system.Builder
func NewPropertyBuilder(path string, customProperties map[string]interface{}) Builder {
	b := &propertyBuilder{
		ConfigFile:        ConfigFile{path: path},
		Viper:             viper.New(),
		defaultProperties: customProperties,
	}

	return b
}

// setCustomPropertiesFromArgs returns application config
func (b *propertyBuilder) setCustomPropertiesFromArgs() {
	log.Debug(os.Args)
	for _, val := range os.Args {
		if len(val) < 2 {
			continue
		}
		prefix := val[:2]
		if prefix == "--" {
			kv := val[2:]
			kvPair := strings.Split(kv, "=")
			// --property equal to --property=true
			if len(kvPair) == 1 {
				kvPair = append(kvPair, "true")
			}
			var v interface{}
			v = kvPair[1]
			switch v.(type) {
			case string:
				if strings.Contains(v.(string), ",") {
					v = strings.SplitN(v.(string), ",", -1)
				}
			}
			b.Set(kvPair[0], v)
		}
	}
}

// New create new viper instance
func (b *propertyBuilder) readConfigData(in io.Reader, ext string) (err error) {
	log.Debugf("reader: %v, ext:%v", in, ext)
	b.AutomaticEnv()
	viperReplacer := strings.NewReplacer(".", "_")
	b.SetEnvKeyReplacer(viperReplacer)
	b.SetConfigType(ext)
	if !b.merge {
		b.merge = true
		err = b.ReadConfig(in)
	} else {
		err = b.MergeConfig(in)
	}
	return
}

// New create new viper instance
func (b *propertyBuilder) readConfig(path, file, ext string) (err error) {
	log.Debugf("file: %v%v.%v", path, file, ext)
	b.AutomaticEnv()
	viperReplacer := strings.NewReplacer(".", "_")
	b.SetEnvKeyReplacer(viperReplacer)

	b.AddConfigPath(path)
	b.SetConfigName(file)
	b.SetConfigType(ext)
	if !b.merge {
		b.merge = true
		err = b.ReadInConfig()
	} else {
		err = b.MergeInConfig()
	}
	return
}

// deprecated
func (b *propertyBuilder) BuildWithProfile(profile string) (interface{}, error) {
	return nil, nil
}

// deprecated
func (b *propertyBuilder) SetConfiguration(in interface{}) {
}

// deprecated
// Init create file if it's not exist
func (b *propertyBuilder) Init() error {
	return nil
}

// Save configurations to file
func (b *propertyBuilder) Save(p interface{}) (err error) {
	return
}

// parseConfigFile takes a full file path and splits it into its components.
// It returns a ConfigFile struct and an error if the filename doesn't match the expected pattern.
func (b *propertyBuilder) parseConfigFile(fullPath string) (config *ConfigFile, err error) {
	config = new(ConfigFile)

	// Extract the directory path
	dir := filepath.Dir(fullPath)
	config.path = dir

	// Extract the base filename
	base := filepath.Base(fullPath)

	// Extract the file extension
	ext := filepath.Ext(base)
	if ext == "" {
		return config, fmt.Errorf("file '%s' does not have an extension", fullPath)
	}
	config.fileType = strings.TrimPrefix(ext, ".")

	// Remove the extension to get the filename without extension
	nameWithoutExt := strings.TrimSuffix(base, ext)
	config.name = nameWithoutExt

	// Define the expected prefix
	prefix := "application"

	if nameWithoutExt == prefix {
		// Case: "application.yml" with default profile
		config.profile = "default"
	} else if strings.HasPrefix(nameWithoutExt, prefix+"-") {
		// Case: "application-<profile>.yml"
		profile := strings.TrimPrefix(nameWithoutExt, prefix+"-")
		if profile == "" {
			return config, fmt.Errorf("filename '%s' does not contain a profile (with filename prefix)", fullPath)
		}
		config.profile = profile
	} else {
		return config, nil //fmt.Errorf("filename '%s' does not start with the expected prefix '%s'", fullPath, prefix)
	}

	return config, nil
}

// Build config file
func (b *propertyBuilder) Build(profiles ...string) (conf interface{}, err error) {
	// parse profiles

	// set custom properties
	for key, value := range b.defaultProperties {
		b.SetDefaultProperty(key, value)
	}

	b.setCustomPropertiesFromArgs()

	// External Config files
	var embedPaths []string
	var paths []string
	embedConfigFiles := make(map[string]map[string][]*ConfigFile)
	pp, _ := filepath.Abs(b.path)

	profile := b.GetString("app.profiles.active")
	if profile == "" {
		profile = b.GetString("profile")
	}
	if profile == "" && len(profiles) > 0 {
		profile = profiles[0]
	}

	// TODO: should combine below two process into one
	var embedActiveProfileConfigFile *ConfigFile
	var embedDefaultProfileConfigFile *ConfigFile

	// Embed Config Files
	cfg := b.Get(Config)
	if cfg != nil {
		switch cfg.(type) {
		case embed.FS:
			c := cfg.(embed.FS)
			b.embedFS = &c
		case *embed.FS:
			b.embedFS = cfg.(*embed.FS)
		}
		dir := b.GetString(ConfigDir)
		if dir == "" {
			dir = "config"
		}
		var files []fs.DirEntry
		files, err = b.embedFS.ReadDir(dir)

		for _, f := range files {
			name, _, isDir := f.Name(), f.Type(), f.IsDir()
			if isDir {
				continue
			}

			configFile, err := b.parseConfigFile(name)
			if err != nil {
				log.Warnf("skip %v, error: %v", name, err)
				continue
			}
			if str.InSlice(configFile.fileType, viper.SupportedExts) {
				if embedConfigFiles[dir] == nil {
					embedConfigFiles[dir] = make(map[string][]*ConfigFile)
				}
				fd, e := b.embedFS.Open(filepath.Join(dir, name))
				if e == nil {
					configFile.fd = fd
					if configFile.profile == "default" {
						embedDefaultProfileConfigFile = configFile
						continue
					}
					if configFile.profile != "" {
						if configFile.profile == profile {
							embedActiveProfileConfigFile = configFile
							continue
						} else {
							embedConfigFiles[dir][configFile.profile] = append(embedConfigFiles[dir][configFile.profile], configFile)
						}
					} else {
						embedConfigFiles[dir][configFile.profile] = append(embedConfigFiles[dir][configFile.profile], configFile)
					}
					foundDir := false
					for _, d := range embedPaths {
						if d == dir {
							foundDir = true
							break
						}
					}
					if !foundDir {
						embedPaths = append(embedPaths, dir)
					}
				}
			}
			//}
		}
	}

	// external files
	var activeProfileConfigFile *ConfigFile
	var defaultProfileConfigFile *ConfigFile
	configFiles := make(map[string]map[string][]string)
	err = filepath.Walk(pp, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			//*files = append(*files, path)
			if !info.IsDir() {
				//log.Debug(path)
				configFile, err := b.parseConfigFile(path)
				if err != nil {
					log.Warnf("skip config file: %v, reason: %v", path, err)
					return nil
				}

				if str.InSlice(configFile.fileType, viper.SupportedExts) {
					if configFiles[configFile.path] == nil {
						configFiles[configFile.path] = make(map[string][]string)
					}
					// do not add default profile, will be handled later
					if configFile.profile == "default" {
						defaultProfileConfigFile = configFile
						return nil
					}

					// TODO: check if profile is filter out properly

					if profile != "" {
						if activeProfileConfigFile == nil && configFile.profile == profile {
							activeProfileConfigFile = configFile
						} else {
							configFiles[configFile.path][configFile.fileType] = append(configFiles[configFile.path][configFile.fileType], configFile.name)
						}
					} else {
						configFiles[configFile.path][configFile.fileType] = append(configFiles[configFile.path][configFile.fileType], configFile.name)
					}
					foundDir := false
					for _, d := range paths {
						if d == configFile.path {
							foundDir = true
							break
						}
					}
					if !foundDir {
						paths = append(paths, configFile.path)
					}
				}
				//}
				//}
			}
		}
		return err
	})

	// sort all config files
	for _, exts := range configFiles {
		for _, files := range exts {
			sort.ByLen(files)
		}
	}
	sort.ByLen(paths)

	// read default profile first
	if embedDefaultProfileConfigFile != nil {
		err = b.readConfigData(embedDefaultProfileConfigFile.fd, embedDefaultProfileConfigFile.fileType)
		if err != nil {
			log.Error(err)
		}
		_ = embedDefaultProfileConfigFile.fd.Close()
	}

	if defaultProfileConfigFile != nil {
		err = b.readConfig(defaultProfileConfigFile.path, defaultProfileConfigFile.name, defaultProfileConfigFile.fileType)
		if err != nil {
			log.Error(err)
		}
	}

	includeProfiles := b.GetStringSlice(appProfilesInclude)

	for _, path := range embedPaths {
		ds := embedConfigFiles[path]
		for _, files := range ds {
			for _, file := range files {
				p := strings.Split(file.name, "-")
				np := len(p)
				if np > 0 && str.InSlice(p[np-1], includeProfiles) {
					err = b.readConfigData(file.fd, file.fileType)
					if err != nil {
						log.Error(err)
					}
					_ = file.fd.Close()
				}
			}
		}
	}

	// read all config files
	//log.Debug("after ...")
	for _, path := range paths {
		ds := configFiles[path]
		for ext, files := range ds {
			for _, file := range files {
				p := strings.Split(file, "-")
				np := len(p)
				if np > 0 && str.InSlice(p[np-1], includeProfiles) {
					err = b.readConfig(path, file, ext)
					if err != nil {
						log.Error(err)
					}
				}
			}
		}
	}

	// replaced with active profile
	if embedActiveProfileConfigFile != nil {
		err = b.readConfigData(embedActiveProfileConfigFile.fd, embedActiveProfileConfigFile.fileType)
		if err != nil {
			log.Error(err)
		}
		_ = embedActiveProfileConfigFile.fd.Close()
	}

	if activeProfileConfigFile != nil {
		err = b.readConfig(activeProfileConfigFile.path, activeProfileConfigFile.name, activeProfileConfigFile.fileType)
		if err != nil {
			log.Error(err)
		}
	}

	// iterate all and replace reference values or env
	allKeys := b.AllKeys()
	for _, key := range allKeys {
		val := b.GetString(key)
		if strings.Contains(val, "${") {
			newVal := b.Replace(val)
			b.SetConfig(key, newVal)
			log.Debugf(">>> replaced key: %v, value: %v, newVal: %v", key, val, newVal)
		}
	}

	log.Debugf("active profile: %v", activeProfileConfigFile)
	return
}

// Load single file
func (b *propertyBuilder) Load(properties interface{}, opts ...func(*mapstructure.DecoderConfig)) (err error) {
	ann := annotation.GetAnnotation(properties, at.ConfigurationProperties{})
	if ann != nil {
		prefix := ann.Field.StructField.Tag.Get("value")

		allSettings := b.AllSettings()
		settings := allSettings[prefix]
		if settings != nil {
			err = mapstruct.Decode(properties, settings, opts...)
		}
	}
	return
}

// Replace reference and
func (b *propertyBuilder) Replace(source string) (retVal interface{}) {
	result := source
	matches := replacer.GetMatches(source)
	if len(matches) != 0 {
		for _, m := range matches {
			varFullName, varName := m[0], m[1]
			// check if it contains default value
			var defaultValue string
			n := strings.Index(varName, ":")
			if n > 0 {
				defaultValue = varName[n+1:]
				varName = varName[:n]
				//log.Debugf("name: %v, default value: %v", varName, defaultValue)
			}
			prop := b.Get(varName)

			var newVal string
			if prop != nil {
				switch prop.(type) {
				case string:
					newVal = prop.(string)
					result = strings.Replace(result, varFullName, newVal, -1)
				default:
					retVal = prop
					return
				}
			}

			envValue := os.Getenv(varName)
			// check if  varName == strings.ToUpper(varName), the assume that varName is environment variable
			if envValue != "" || (varName == strings.ToUpper(varName) && defaultValue == "") {
				result = strings.Replace(result, varFullName, envValue, -1)
			}

			if envValue == "" && newVal == "" && defaultValue != "" {
				result = strings.Replace(result, varFullName, defaultValue, -1)
			}
			log.Debugf("replaced %v to %v", varName, result)
		}
	}
	retVal = result
	return
}

func (b *propertyBuilder) GetProperty(name string) (retVal interface{}) {
	retVal = b.Get(name)
	return
}

func (b *propertyBuilder) SetProperty(name string, val interface{}) Builder {
	b.Set(name, val)
	return b
}

func (b *propertyBuilder) SetDefaultProperty(name string, val interface{}) Builder {
	b.Lock()
	b.SetDefault(name, val)
	b.Unlock()

	return b
}
