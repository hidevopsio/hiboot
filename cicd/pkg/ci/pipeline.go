// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ci

import (
	"github.com/hidevopsio/hi/boot/pkg/log"
	"fmt"
	"github.com/hidevopsio/hi/cicd/pkg/orch/k8s"
	"github.com/hidevopsio/hi/cicd/pkg/orch/openshift"
	"os"
	"github.com/hidevopsio/hi/boot/pkg/system"
	"github.com/hidevopsio/hi/cicd/pkg/orch"
	"github.com/imdario/mergo"
	"github.com/hidevopsio/hi/boot/pkg/utils"
)

// config file
// pipeline:
// - PullSourceCode
// - Build
// - RunUnitTest
// - Analysis
// - CopyTarget
// - Upload
// - NewImage
// - Deploy

type PipelineInterface interface {
	Init(pl *Pipeline)
	CreateSecret(username, password string, isToken bool) (string, error)
	Build(secret string) error
	RunUnitTest() error
	RunIntegrationTest() error
	Analysis() error
	CreateDeploymentConfig() error
	InjectSideCar() error
	Deploy() error
	CreateService() error
	CreateRoute() error
	Run(username, password string, isToken bool) error
}

//type Ports struct {
//	Name          string `json:"name"`
//	Port          int32  `json:"port" `
//	ContainerPort int32  `json:"container_port"`
//	Protocol      string `json:"protocol,omitempty"`
//}

type Pipeline struct {
	Name           string       `json:"name"`
	ScmType        string       `json:"scm_type"`
	ScmUrl         string       `json:"scm_url"`
	ScmRef         string       `json:"scm_ref"`
	App            string       `json:"app"`
	Profile        string       `json:"profile"`
	Project        string       `json:"project"`
	Namespace      string       `json:"namespace"`
	Replicas       int32        `json:"replicas"`
	Version        string       `json:"version"`
	ImageTag       string       `json:"image_tag"`
	ImageStream    string       `json:"image_stream"`
	DockerRegistry string       `json:"docker_registry"`
	RepositoryUrl  string       `json:"repository_url"`
	Identifiers    []string     `json:"identifiers"`
	ConfigFiles    []string     `json:"config_files"`
	Env            []system.Env `json:"env"`
	Ports          []orch.Ports `json:"ports"`
	Script         string       `json:"script"`
}

// @Title Init
// @Description set default value
// @Param pipeline
// @Return error
func (p *Pipeline) Init(pl *Pipeline) {
	log.Debug("Pipeline.EnsureParam()")

	// load config file
	if pl != nil {
		c := Build(pl.Name)

		// read env
		mavenMirrorUrl := os.Getenv("MAVEN_MIRROR_URL")
		if mavenMirrorUrl != "" && pl.RepositoryUrl == "" {
			pl.RepositoryUrl = mavenMirrorUrl
		}

		if pl.Profile == "" {
			p.Profile = "dev"
		}

		mergo.Merge(&c.Pipeline, pl, mergo.WithOverride)
		mergo.Merge(p, c.Pipeline, mergo.WithOverride)

	}
	// TODO: replace variable inside pipeline, e.g. ${profile}
	utils.Replace(p, "profile", pl.Profile)

	if "" == p.Namespace {
		if "" == pl.Profile {
			p.Namespace = p.Project
		} else {

			p.Namespace = p.Project + "-" + p.Profile
		}
	}

	log.Debug(p)
}

func (p *Pipeline) CreateSecret(username, password string, isToken bool) (string, error) {
	log.Debug("Pipeline.CreateSecret()")
	if username == "" {
		return "", fmt.Errorf("unkown username")
	}
	// Create secret
	secretName := username + "-secret"
	secret := k8s.NewSecret(secretName, username, password, p.Namespace, isToken)
	err := secret.Create()

	return secretName, err
}

func (p *Pipeline) Build(secret string) error {
	log.Debug("Pipeline.Build()")

	scmUrl := p.CombineScmUrl()
	buildConfig, err := openshift.NewBuildConfig(p.Namespace, p.App, scmUrl, p.ScmRef, secret, p.ImageTag, p.ImageStream)
	if err != nil {
		return err
	}
	_, err = buildConfig.Create()
	if err != nil {
		return err
	}
	// Build image stream
	_, err = buildConfig.Build(p.RepositoryUrl, p.Script)
	return err
}

func (p *Pipeline) CombineScmUrl() string {
	scmUrl := p.ScmUrl + "/" + p.Project + "/" + p.App + "." + p.ScmType
	return scmUrl
}

func (p *Pipeline) RunUnitTest() error {
	log.Debug("Pipeline.RunUnitTest()")
	return nil
}

func (p *Pipeline) RunIntegrationTest() error {
	log.Debug("Pipeline.RunIntegrationTest()")
	return nil
}

func (p *Pipeline) Analysis() error {
	log.Debug("Pipeline.Analysis()")
	return nil
}

func (p *Pipeline) CreateDeploymentConfig() error {
	log.Debug("Pipeline.CreateDeploymentConfig()")

	// new dc instance
	dc, err := openshift.NewDeploymentConfig(p.App, p.Namespace)
	if err != nil {
		return err
	}

	// create dc
	err = dc.Create(&p.Env, &p.Ports, p.Replicas)
	if err != nil {
		return err
	}
	return nil
}

func (p *Pipeline) InjectSideCar() error {
	log.Debug("Pipeline.InjectSideCar()")
	return nil
}

func (p *Pipeline) Deploy() error {
	log.Debug("Pipeline.Deploy()")

	// new dc instance
	dc, err := openshift.NewDeploymentConfig(p.App, p.Namespace)
	if err != nil {
		return err
	}

	_, err = dc.Instantiate()
	if err != nil {
		return err
	}

	return nil
}

func (p *Pipeline) CreateService() error {
	log.Debug("Pipeline.CreateService()")

	// new dc instance
	svc := k8s.NewService(p.App, p.Namespace)

	err := svc.Create(&p.Ports)
	if err != nil {
		return err
	}

	return nil
}

func (p *Pipeline) CreateRoute() error {
	log.Debug("Pipeline.CreateRoute()")

	route, err := openshift.NewRoute(p.App, p.Namespace)
	if err != nil {
		return err
	}

	err = route.Create(8080)
	return nil
}

func (p *Pipeline) Run(username, password string, isToken bool) error {
	// first, let's check if namespace is exist or not

	// create secret for building image
	secret, err := p.CreateSecret(username, password, isToken)
	if err != nil {
		return fmt.Errorf("failed on CreateSecret! %s", err.Error())
	}

	// build image
	err = p.Build(secret)
	if err != nil {
		return fmt.Errorf("failed on Build! %s", err.Error())
	}

	// create dc - deployment config
	err = p.CreateDeploymentConfig()
	if err != nil {
		return fmt.Errorf("failed on CreateDeploymentConfig! %s", err.Error())
	}

	// inject side car
	err = p.InjectSideCar()
	if err != nil {
		return fmt.Errorf("failed on InjectSideCar! %s", err.Error())
	}

	// deploy
	err = p.Deploy()
	if err != nil {
		return fmt.Errorf("failed on Deploy! %s", err.Error())
	}

	// create service
	err = p.CreateService()
	if err != nil {
		return fmt.Errorf("failed on CreateService! %s", err.Error())
	}

	// create route
	err = p.CreateRoute()
	if err != nil {
		return fmt.Errorf("failed on CreateRoute! %s", err.Error())
	}

	// finally, all steps are done well, let tell the client ...
	return nil
}
