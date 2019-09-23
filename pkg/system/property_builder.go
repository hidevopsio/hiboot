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
	"github.com/hidevopsio/mapstructure"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/inject/annotation"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/utils/mapstruct"
	"hidevops.io/hiboot/pkg/utils/replacer"
	"hidevops.io/hiboot/pkg/utils/sort"
	"hidevops.io/hiboot/pkg/utils/str"
	"hidevops.io/viper"
	"os"
	"path/filepath"
	"strings"
)

const ( appProfilesInclude  = "app.profiles.include")

type ConfigFile struct{
	path             string
	name             string
	fileType         string
}

type propertyBuilder struct {
	at.Qualifier `value:"system.builder"`
	*viper.Viper
	ConfigFile
	configuration     interface{}
	defaultProperties map[string]interface{}
	profiles          []string
	merge             bool
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
	log.Println(os.Args)
	for _, val := range os.Args {
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

// Build config file
func (b *propertyBuilder) Build(profiles ...string) (conf interface{}, err error) {
	// parse profiles

	// set custom properties
	for key, value := range b.defaultProperties {
		b.SetDefaultProperty(key, value)
	}

	b.setCustomPropertiesFromArgs()

	var paths []string
	configFiles :=  make(map[string]map[string][]string)
	pp, _ := filepath.Abs(b.path)
	var profile string

	if len(profiles) > 0 {
		profile = profiles[0]
	}

	var activeProfileConfigFile *ConfigFile
	var defaultProfileConfigFile *ConfigFile
	err = filepath.Walk(pp, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			//*files = append(*files, path)
			if !info.IsDir() {
				//log.Debug(path)
				dir, file := filepath.Split(path)
				fileAndExt := strings.Split(file, ".")
				if len(fileAndExt) == 2 {
					file, ext := fileAndExt[0], fileAndExt[1]
					if file != "" {
						if str.InSlice(ext, viper.SupportedExts) {
							if configFiles[dir] == nil {
								configFiles[dir] = make(map[string][]string)
							}
							files := configFiles[dir][ext]
							// do not add default profile, will be handled later
							if !strings.Contains(file, "-") {
								defaultProfileConfigFile = &ConfigFile{
									path:     dir,
									name:     file,
									fileType: ext,
								}
								return nil
							}

							if profile != "" {
								if strings.Contains(file, profile) {
									activeProfileConfigFile = &ConfigFile{
										path:     dir,
										name:     file,
										fileType: ext,
									}
								} else {
									files = append(files, file)
								}
							} else {
								files = append(files, file)
							}
							configFiles[dir][ext] = files
							foundDir := false
							for _, d := range paths {
								if d == dir {
									foundDir = true
									break
								}
							}
							if !foundDir {
								paths = append(paths, dir)
							}
						}
					}
				}
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
	if defaultProfileConfigFile != nil {
		err = b.readConfig(defaultProfileConfigFile.path, defaultProfileConfigFile.name, defaultProfileConfigFile.fileType)
	}

	var includeProfiles []string
	ip := b.Get(appProfilesInclude)
	if ip != nil {
		switch ip.(type) {
		case []string:
			includeProfiles = ip.([]string)
		case []interface{}:
			ipi := ip.([]interface{})
			for _, value := range ipi {
				includeProfiles = append(includeProfiles, value.(string))
			}
		case string:
			includeProfiles = append(includeProfiles, ip.(string))
			// set []string back to property appProfilesInclude if it is a string
			b.SetDefaultProperty(appProfilesInclude, includeProfiles)
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
				if  np > 0 && str.InSlice(p[np - 1], includeProfiles) {
					err = b.readConfig(path, file, ext)
				}
			}
		}
	}

	// replaced with active profile
	if activeProfileConfigFile != nil {
		err = b.readConfig(activeProfileConfigFile.path, activeProfileConfigFile.name, activeProfileConfigFile.fileType)
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

// Read single file
func (b *propertyBuilder) Load(properties interface{}, opts ...func (*mapstructure.DecoderConfig)) (err error) {
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

// Replace replace reference and
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
//
//func (b *propertyBuilder) updateProperty(name string, val interface{}) (retVal interface{})  {
//	// TODO: for debug only, TBD
//	if name == "swagger" {
//		log.Debug(name)
//	}
//	original := b.Get(name)
//	// convert struct to map
//	var dm = make(map[string]interface{})
//	sm, ok := mapstruct.DecodeStructToMap(val)
//	if ok {
//		if original != nil {
//			// copy original map to the new map
//			copier.CopyMap(dm, original.(map[string]interface{}))
//			// copy new src map to dest map
//			copier.CopyMap(dm, sm, copier.IgnoreEmptyValue)
//			// assign dest map to retVal
//			retVal = dm
//		} else {
//			retVal = sm
//		}
//	} else {
//		retVal = val
//	}
//	return
//}

func (b *propertyBuilder) SetProperty(name string, val interface{}) Builder {
	b.Set(name, val)
	return b
}

func (b *propertyBuilder) SetDefaultProperty(name string, val interface{}) Builder {
	b.SetDefault(name, val)
	return b
}
