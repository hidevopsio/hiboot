// TODO: app config should be generic kit

package config

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/hidevopsio/hi/cicd/pkg/pipelines"
)

type Configuration struct {
	PipelineConf pipelines.Pipeline `mapstructure:"pipeline"`
}

func BuildPipelineConfig(profile, project, app, version, appType string) *Pipeline {
	viper.SetDefault("pipeline.profile", profile)
	viper.SetDefault("pipeline.project", project)
	viper.SetDefault("pipeline.app", app)
	viper.SetDefault("pipeline.version", version)
	viper.SetDefault("author", "John Deng <john.deng@outlook.com>")
	viper.SetDefault("license", "Apache 2.0")
	viper.AddConfigPath("./config")
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Error config file: %s \n", err))
	}

	var conf Configuration

	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(fmt.Errorf("Unable to decode Config: %s \n", err))
	}

	viper.SetConfigName("pipeline-" + appType)
	viper.MergeInConfig()
	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(fmt.Errorf("Unable to decode Config: %s \n", err))
	}

	return &conf.PipelineConf
}
