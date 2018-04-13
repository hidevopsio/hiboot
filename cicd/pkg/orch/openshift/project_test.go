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
	"k8s.io/apimachinery/pkg/api/errors"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestProjectLit(t *testing.T) {
	// get all projects
	project, err := NewProject("", "", "")
	assert.Equal(t, nil, err)

	pl, err := project.List()
	assert.Equal(t, nil, err)
	assert.NotEqual(t, 0, len(pl.Items))
	log.Debugf("There are %d projects in the cluster", len(pl.Items))

	for i, p := range pl.Items {
		log.Debugf("index %d: project: %s", i, p.Name)
	}
}

func TestProjectCreate(t *testing.T) {
	namespace := "moses-demos-dev"
	projects, err := NewProject(namespace, "", "")
	assert.Equal(t, nil, err)

	_, err = projects.Get()
	if errors.IsNotFound(err) {
		projects.Create()
	} else {
		log.Debugf("%v is already exist", namespace)
	}
}

func TestProjectCrud(t *testing.T) {
	projectName := "project-crud"
	project, err := NewProject(projectName, projectName, "project for testing")
	assert.Equal(t, nil, err)

	// create project
	p, err := project.Create()
	assert.Equal(t, nil, err)
	assert.Equal(t, projectName, p.Name)

	// read project
	p, err = project.Get()
	assert.Equal(t, nil, err)
	assert.Equal(t, projectName, p.Name)

	// delete project
	err = project.Delete()
	assert.Equal(t, nil, err)

}

