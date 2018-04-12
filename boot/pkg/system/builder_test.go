package system


type Configuration struct {
	App     App     `mapstructure:"app"`
	Server  Server  `mapstructure:"server"`
	Logging Logging `mapstructure:"logging"`
}


const (
	application = "application"
	config = "/config"
	yaml = "yaml"
)

//wd := utils.GetWorkingDir("boot/pkg/system/builder.go")
//appProfile := os.Getenv("APP_PROFILES_ACTIVE")
//if appProfile != "" {
//	viper.SetConfigName(application + "-" + appProfile)
//	viper.MergeInConfig()
//	err = viper.Unmarshal(&conf)
//	if err != nil {
//		panic(fmt.Errorf("Unable to decode Config: %s \n", err))
//	}
//}
