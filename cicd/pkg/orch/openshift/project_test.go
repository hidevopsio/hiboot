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
	"github.com/openshift/api/project/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	projectv1 "github.com/openshift/client-go/project/clientset/versioned/typed/project/v1"
	"github.com/hidevopsio/hi/cicd/pkg/orch/k8s"
	"github.com/hidevopsio/hi/boot/pkg/log"
	"os"
)

var (
	projectV1Client *projectv1.ProjectV1Client
	projects        projectv1.ProjectInterface
)

func init() {
	log.SetLevel(log.DebugLevel)
	var err error
	projectV1Client, err = projectv1.NewForConfig(k8s.Config)
	if err != nil {
		os.Exit(1)
	}
	projects = projectV1Client.Projects()
}

func TestProjectLit(t *testing.T) {
	// get all projects
	project, err := projects.List(metav1.ListOptions{})
	assert.Equal(t, nil, err)
	assert.NotEqual(t, 0, len(project.Items))
	log.Debugf("There are %d projects in the cluster", len(project.Items))

	for i, p := range project.Items {
		log.Debugf("index %d: project: %s", i, p.Name)
	}
}

func TestProjectCrud(t *testing.T) {
	projectName := "project-crud"
	ps := &v1.Project{
		ObjectMeta: metav1.ObjectMeta{
			Name: projectName,
			Labels: map[string]string{
				"project": projectName,
			},
		},
	}
	var err error

	// create project
	ps, err = projects.Create(ps)
	assert.Equal(t, nil, err)

	// read project
	p, err := projects.Get(projectName, metav1.GetOptions{})
	assert.Equal(t, nil, err)
	assert.Equal(t, projectName, p.Name)

	// delete project
	err = projects.Delete(projectName, &metav1.DeleteOptions{})
	assert.Equal(t, nil, err)

}

