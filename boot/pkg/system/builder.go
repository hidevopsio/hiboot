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
	"os"
	"runtime"
	"strings"
	"github.com/hidevopsio/hi/boot/pkg/log"
)

type Configuration struct {
	App     App     `mapstructure:"app"`
	Server  Server  `mapstructure:"server"`
	Logging Logging `mapstructure:"logging"`
}

var (
	conf Configuration
)

const (
	application = "application"
	config = "/config"
	yaml = "yaml"
	defaultLoggingLevel = "info"
	defaultPort = 8080
	defaultAppName = "app"
	defaultProjectName = "devops"

)

func Build() *Configuration {
	_, filename, _, _ := runtime.Caller(0)
	workdir := strings.Replace(filename, "boot/pkg/system/builder.go", "", -1)
	log.Debug(workdir)

	viper.SetDefault("app.project", defaultProjectName)
	viper.SetDefault("app.name", defaultAppName)
	viper.SetDefault("server.port", defaultPort)
	viper.SetDefault("logging.level", defaultLoggingLevel)
	viper.AddConfigPath(workdir + config)
	viper.SetConfigName(application)
	viper.SetConfigType(yaml)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Error config file: %s \n", err))
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(fmt.Errorf("Unable to decode Config: %s \n", err))
	}
	appProfile := os.Getenv("APP_PROFILES_ACTIVE")
	if appProfile != "" {
		viper.SetConfigName(application + "-" + appProfile)
		viper.MergeInConfig()
		err = viper.Unmarshal(&conf)
		if err != nil {
			panic(fmt.Errorf("Unable to decode Config: %s \n", err))
		}
	}

	return &conf
}
