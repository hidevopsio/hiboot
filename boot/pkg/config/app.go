package config

type Server struct {
	Port int32 `json:"port"`
}

type Logging struct{
	Level string `json:"level"`
}

type AppConfig struct {
	Project string `json:"project"`
	Name string `json:"name"`
	Server Server `json:"server"`
	Logging Logging `json:"logging"`
}

