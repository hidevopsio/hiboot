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
	"reflect"
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
	CreateSecret(username, password string, isToken bool) error
	Build() error
	RunUnitTest() error
	RunIntegrationTest() error
	Analysis() error
	CreateDeploymentConfig() error
	InjectSideCar() error
	Deploy() error
	Run(username, password string, isToken bool) error
}

type Pipeline struct {
	Name           string   `json:"name"`
	GitUrl         string   `json:"git_url"`
	App            string   `json:"app"`
	Profile        string   `json:"profile"`
	Project        string   `json:"project"`
	Namespace      string   `json:"namespace"`
	Version        string   `json:"version"`
	ImageTag       string   `json:"image_tag"`
	Type           string   `json:"type"`
	Timezone       string   `json:"timezone"`
	Identifiers    []string `json:"identifiers"`
	Targets        []string `json:"targets"`
	ConfigFiles    []string `json:"config_files"`
	FromDir        string   `json:"from_dir"`
	DeploymentFile string   `json:"deployment_file"`
	ImageStream    string   `json:"image_stream"`
	VersionFrom    string   `json:"version_from"`
	Options        string   `json:"options"`
	BuildCommand   string   `json:"build_command"`
	DockerRegistry string   `json:"docker_registry"`
}

func copyPipeline(to *Pipeline, from *Pipeline) {
	f := reflect.ValueOf(from).Elem()
	t := reflect.ValueOf(to).Elem()

	for i := 0; i < f.NumField(); i++ {
		varName := f.Type().Field(i).Name
		//varType := f.Type().Field(i).Type
		varValue := f.Field(i).Interface()
		//log.Debugf("%v %v %v\n", varName, varType, varValue)
		tf := t.FieldByName(varName)

		if tf.IsValid() && tf.CanSet() {
			kind := tf.Kind()
			switch kind {
			case reflect.String:
				fv := fmt.Sprintf("%v", varValue)
				if fv != "" {
					tf.SetString(fmt.Sprintf("%v", varValue))
				}
				break
			case reflect.Slice:
				break
			default:
				break
			}

		}
	}
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
		log.Debug(c.Pipeline)
		log.Debug(p)

		log.Debug("Pipeline Input Param ...")

		copyPipeline(&c.Pipeline, pl)
		copyPipeline(p, &c.Pipeline)

		log.Debug(p)
	}

	if "" == p.ImageTag {
		p.ImageTag = "latest"
	}
	if "" == p.DockerRegistry {
		p.DockerRegistry = "docker-registry.default.svc:5000"
	}
	if "" == p.Profile {
		p.Profile = "dev"
	}

	p.Namespace = p.Project + "-" + p.Profile
}

func (p *Pipeline) CreateSecret(username, password string, isToken bool) error {
	log.Debug("Pipeline.CreateSecret()")

	// Create secret
	secret := k8s.NewSecret(username + "-secret", username, password, p.Namespace, isToken)
	err := secret.Create()

	return err
}

func (p *Pipeline) Build() error {
	log.Debug("Pipeline.Build()")
	return nil
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
	return nil
}

func (p *Pipeline) InjectSideCar() error {
	log.Debug("Pipeline.InjectSideCar()")
	return nil
}

func (p *Pipeline) Deploy() error {
	log.Debug("Pipeline.Deploy()")
	return nil
}

func (p *Pipeline) Run(username, password string, isToken bool) error {
	err := p.CreateSecret(username, password, isToken)
	if err != nil {
		return fmt.Errorf("failed! %s", err.Error())
	}

	err = p.Build()
	if err != nil {
		return fmt.Errorf("failed! %s", err.Error())
	}

	err = p.CreateDeploymentConfig()
	if err != nil {
		return fmt.Errorf("failed! %s", err.Error())
	}

	err = p.InjectSideCar()
	if err != nil {
		return fmt.Errorf("failed! %s", err.Error())
	}

	err = p.Deploy()
	if err != nil {
		return fmt.Errorf("failed! %s", err.Error())
	}

	return nil
}
