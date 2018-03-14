package config

type Pipeline struct {
	Name string `json:"name"`
	Profile string `json:"profile"`
	Project string `json:"project"`
	App string `json:"app"`
	Version string `json:"version"`
	ImageTag string `json:"image_tag"`
	Type string `json:"type"`
	Timezone string `json:"timezone"`
	Identifiers []string `json:"identifiers"`
	Targets []string `json:"targets"`
	ConfigFiles []string `json:"config_files"`
	FromDir string `json:"from_dir"`
	DeploymentFile string `json:"deployment_file"`
	ImageStream string `json:"image_stream"`
	VersionFrom string `json:"version_from"`
	Options string `json:"options"`
	BuildCommand string `json:"build_command"`
	DockerRegistry string `json:"docker_registry"`
}

