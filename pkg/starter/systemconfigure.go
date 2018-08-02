package starter

type Profiles struct {
	Include []string `json:"include"`
	Active  string   `json:"active"`
}

type App struct {
	Project        string   `json:"project"`
	Name           string   `json:"name"`
	Profiles       Profiles `json:"profiles"`
	DataSourceType string   `json:"data_source_type"`
}

type Server struct {
	Port int32 `json:"port"`
}

type Logging struct {
	Level string `json:"level"`
}

type Env struct {
	Name  string
	Value string
}

type SystemConfiguration struct {
	App         App          `mapstructure:"app"`
	Server      Server       `mapstructure:"server"`
	Logging     Logging      `mapstructure:"logging"`
}

