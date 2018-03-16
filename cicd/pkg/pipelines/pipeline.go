package pipelines

import (
	"fmt"
	"github.com/hidevopsio/hi/boot/pkg/log"
)

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
	Run() error
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


// @Title EnsurePipline
// @Description set default value
// @Param pipeline
// @Return error
func (p *Pipeline) EnsurePipeline()  {
	if "" == p.ImageTag {
		p.ImageTag = "latest"
	}
	if "" == p.DockerRegistry {
		p.DockerRegistry = "docker-registry.default.svc:5000"
	}
	if "" == p.Profile {
		p.Profile = "dev"
	}
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

func (p *Pipeline) PullSourceCode() error {
	log.Debug("Pipeline.PullSourceCode()")
	return nil
}

func (p *Pipeline) Build() error {
	return nil
}

func (p *Pipeline) RunUnitTest() error {
	return nil
}

func (p *Pipeline) RunIntegrationTest() error {
	return nil
}

func (p *Pipeline) CopyTarget() error {
	return nil
}

func (p *Pipeline) Upload() error {
	return nil
}

func (p *Pipeline) NewImage() error {
	return nil
}

func (p *Pipeline) Deploy() error {
	log.Debug("Pipeline.Deploy()")
	return nil
}

func (p *Pipeline) Run() error {
	err := p.PullSourceCode()
	if err != nil {
		return fmt.Errorf("failed: %s", err)
	}

	err = p.Build()
	if err != nil {
		return fmt.Errorf("failed: %s", err)
	}

	err = p.RunUnitTest()
	if err != nil {
		return fmt.Errorf("failed: %s", err)
	}

	return nil
}

