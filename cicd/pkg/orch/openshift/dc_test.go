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

package openshift

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hi/boot/pkg/log"
	"github.com/hidevopsio/hi/boot/pkg/system"
	"github.com/hidevopsio/hi/cicd/pkg/orch"
)

func TestDeploymentConfigCreation(t *testing.T) {
	log.Debug("TestDeploymentConfigCreation()")

	projectName := "demo"
	profile := "dev"
	namespace := projectName + "-" + profile
	app := "hello-world"

	env := []system.Env{
		{
			Name:  "SPRING_PROFILES_ACTIVE",
			Value: profile,
		},
		{
			Name:  "APP_OPTIONS",
			Value: "-Xms128m -Xmx512m -Xss512k -XX:+ExitOnOutOfMemoryError",
		},
		{
			Name:  "TZ",
			Value: "Asia/Shanghai",
		},
	}

	ports := []orch.Ports{
		{
			ContainerPort:     8080,
			Protocol: "TCP",
		},
		{
			ContainerPort:     7575,
			Protocol: "TCP",
		},
	}

	// new dc instance
	dc, err := NewDeploymentConfig(app, namespace)
	assert.Equal(t, nil, err)
	assert.Equal(t, app, dc.Name)

	// create dc
	err = dc.Create(&env, &ports, 1, false)
	assert.Equal(t, nil, err)
}

func TestDeploymentConfigInstantiation(t *testing.T) {

	log.Debug("TestDeploymentConfigInstantiation()")

	projectName := "demo"
	profile := "dev"
	namespace := projectName + "-" + profile
	app := "hello-world"

	dc, err := NewDeploymentConfig(app, namespace)
	assert.Equal(t, nil, err)
	assert.Equal(t, app, dc.Name)

	cfg, err := dc.Instantiate()
	assert.Equal(t, nil, err)
	assert.Equal(t, app, cfg.Name)
}

func TestDeploymentConfigDeletion(t *testing.T) {

	log.Debug("TestDeploymentConfigDeletion()")

	projectName := "demo"
	profile := "dev"
	namespace := projectName + "-" + profile
	app := "hello-world"

	dc, err := NewDeploymentConfig(app, namespace)
	assert.Equal(t, nil, err)
	assert.Equal(t, app, dc.Name)

	err = dc.Delete()
	assert.Equal(t, nil, err)
}
