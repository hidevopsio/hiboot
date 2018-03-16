package pipeline

import "fmt"

type PipelineInterface interface {
	PullSourceCode() error
	Build() error
	RunUnitTest() error
	RunIntegrationTest() error
	Analysis() error
	CopyTarget() error
	Upload() error
	NewImage() error
	Deploy() error
}


type Pipeline struct {
	Name string `json:"name"`
	Profile string `json:"profile"`
	Project string `json:"project"`	// Project = Namespace
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


// pipeline:
// - PullSourceCode
// - Build
// - RunUnitTest
// - Analysis
// - CopyTarget
// - Upload
// - NewImage
// - Deploy

func Run(pipeline PipelineInterface) error {
	err := pipeline.PullSourceCode()
	if err != nil {
		return fmt.Errorf("failed: %s", err)
	}

	err = pipeline.Build()
	if err != nil {
		return fmt.Errorf("failed: %s", err)
	}

	err = pipeline.RunUnitTest()
	if err != nil {
		return fmt.Errorf("failed: %s", err)
	}

	return nil
}

