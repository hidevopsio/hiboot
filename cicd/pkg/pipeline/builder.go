// TODO: app config should be generic kit

package pipeline

import (
	"fmt"
	"github.com/spf13/viper"
	"runtime"
	"github.com/hidevopsio/hi/boot/pkg/log"
	"strings"
)

type Configuration struct {
	Pipeline Pipeline `mapstructure:"pipeline"`
}

var conf Configuration


// TODO: we may eliminate this func, merge to system.builder instead
func Build(name string) *Configuration {

	_, filename, _, _ := runtime.Caller(0)
	workdir := strings.Replace(filename, "/pkg/pipeline/builder.go", "", -1)
	log.Debug(workdir)

	viper.AddConfigPath(workdir + "/config")
	viper.SetConfigName("pipeline")
	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer()
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Error config file: %s \n", err))
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(fmt.Errorf("Unable to decode Config: %s \n", err))
	}

	viper.SetConfigName("pipeline-" + name)
	viper.MergeInConfig()
	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(fmt.Errorf("Unable to decode Config: %s \n", err))
	}

	return &conf
}
