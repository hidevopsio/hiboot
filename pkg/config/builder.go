package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type Configuration struct {
	AppConf AppConfig `mapstructure:"app"`
}

func BuildConfig(profile, project, appName, appVersion, appType string) *AppConfig {
	os.Setenv("APP_PROFILE", profile)
	os.Setenv("APP_PROJECT", project)
	os.Setenv("APP_NAME", appName)
	os.Setenv("APP_VERSION", appVersion)
	viper.AddConfigPath("./config")
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Error config file: %s \n", err))
	}

	var conf Configuration

	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(fmt.Errorf("Unable to decode Config: %s \n", err))
	}

	viper.SetConfigName("app-" + appType)
	viper.MergeInConfig()
	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(fmt.Errorf("Unable to decode Config: %s \n", err))
	}

	return &conf.AppConf
}
