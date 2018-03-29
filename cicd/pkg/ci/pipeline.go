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
	"github.com/hidevopsio/hi/cicd/pkg/orch/openshift"
	"os"
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
	Run(username, password string, isToken bool) error
}

type Env struct {
	Name  string
	Value string
}

type Pipeline struct {
	Name           string   `json:"name"`
	ScmUrl         string   `json:"scm_url"`
	ScmRef         string   `json:"scm_ref"`
	App            string   `json:"app"`
	Profile        string   `json:"profile"`
	Project        string   `json:"project"`
	Namespace      string   `json:"namespace"`
	Version        string   `json:"version"`
	ImageTag       string   `json:"image_tag"`
	//Type           string   `json:"type"` // ?
	//Timezone       string   `json:"timezone"` // ?
	Identifiers    []string `json:"identifiers"`
	//Targets        []string `json:"targets"` // ?
	ConfigFiles    []string `json:"config_files"`
	//FromDir        string   `json:"from_dir"` // ?
	//DeploymentFile string   `json:"deployment_file"` // ?
	ImageStream    string   `json:"image_stream"`
	//VersionFrom    string   `json:"version_from"` // ?
	//Options        string   `json:"options"` // ?
	DockerRegistry string   `json:"docker_registry"`
	RepositoryUrl  string   `json:"repository_url"`
	Env            []Env    `json:"env"`
	Script         string   `json:"script"`
}

func merge(to interface{}, from interface{}) {
	f := reflect.ValueOf(from).Elem()
	t := reflect.ValueOf(to).Elem()

	for i := 0; i < f.NumField(); i++ {
		varName := f.Type().Field(i).Name
		//varType := f.Type().Field(i).Type
		ff := f.Field(i)
		varValue := ff.Interface()
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
				if ff.Len() != 0 {
					log.Println(varName, varValue)
					tf.Set(reflect.ValueOf(varValue))
				}
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

		// read env
		mavenMirrorUrl := os.Getenv("MAVEN_MIRROR_URL")
		if mavenMirrorUrl != "" && pl.RepositoryUrl == "" {
			pl.RepositoryUrl = mavenMirrorUrl
		}

		merge(&c.Pipeline, pl)
		merge(p, &c.Pipeline)

		log.Debug(p)
	}

	if "" == p.Namespace {
		p.Namespace = p.Project + "-" + p.Profile
	}

	// TODO: replace variable inside pipeline, e.g. ${profile}

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

	buildConfig, err := openshift.NewBuildConfig(p.Namespace, p.App, p.ScmUrl, p.ScmRef, secret, p.ImageTag, p.ImageStream)
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

	// last, but not least, let's deploy the app
	err = p.Deploy()
	if err != nil {
		return fmt.Errorf("failed on Deploy! %s", err.Error())
	}

	// finally, all steps are done well, let tell the client ...
	return nil
}
