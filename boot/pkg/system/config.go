package system


type App struct {
	Project string `json:"project"`
	Name string `json:"name"`
}

type Server struct {
	Port int32 `json:"port"`
}

type Logging struct{
	Level string `json:"level"`
}
