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
	"github.com/openshift/api/project/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	projectv1 "github.com/openshift/client-go/project/clientset/versioned/typed/project/v1"
	"github.com/hidevopsio/hi/cicd/pkg/orch/k8s"
	"github.com/hidevopsio/hi/boot/pkg/log"
)

type Project struct{
	Name string
	DisplayName string
	Description string

	Interface projectv1.ProjectInterface
}

func NewProject(name, displayName, desc string) (*Project, error)  {

	clientSet, err := projectv1.NewForConfig(k8s.Config)
	if err != nil {
		return nil, err
	}

	return &Project{
		Name:        name,
		DisplayName: displayName,
		Description: desc,
		Interface:   clientSet.Projects(),
	}, nil
}

func (p *Project) Create() (*v1.Project, error)  {
	log.Debug("Project.Create()")
	ps := &v1.Project{
		ObjectMeta: metav1.ObjectMeta{
			Name: p.Name,
			Labels: map[string]string{
				"project": p.Name,
			},
		},
	}
	// create project
	return p.Interface.Create(ps)
}

func (p *Project) Get() (*v1.Project, error)  {
	log.Debug("Project.Get()")
	return p.Interface.Get(p.Name, metav1.GetOptions{})
}

func (p *Project) List() (*v1.ProjectList, error)  {
	log.Debug("Project.List()")
	return p.Interface.List(metav1.ListOptions{})
}

func (p *Project) Delete() error  {
	log.Debug("Project.Delete()")
	return p.Interface.Delete(p.Name, &metav1.DeleteOptions{})
}