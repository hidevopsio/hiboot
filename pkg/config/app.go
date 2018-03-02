package config


type App struct {
	Profile string
	Project string
	Name string
	Version string
	Type string
}

type AppConfig struct {
	Profile string
	Project string
	Name string
	Version string
	Type string
	Timezone string
	Identifiers []string
	Targets []string
	ConfigFiles []string
	FromDir string
	DeploymentFile string
	ImageStream string
	VersionFrom string
	Options string
	BuildCommand string
}

