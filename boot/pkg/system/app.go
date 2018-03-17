package system

type Server struct {
	Port int32 `json:"port"`
}

type Logging struct{
	Level string `json:"level"`
}

type App struct {
	Project string `json:"project"`
	Name string `json:"name"`
	Server Server `json:"server"`
	Logging Logging `json:"logging"`
}

