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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	"github.com/openshift/api/apps/v1"
	"github.com/hidevopsio/hi/cicd/pkg/orch/k8s"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type DeploymentConfig struct{
	Name string
	Namespace string

	Interface appsv1.DeploymentConfigInterface
	Client *appsv1.AppsV1Client
}

func (dc *DeploymentConfig) Create() {

	var intervalSeconds int64
	var timeoutSeconds int64
	var updatePeriodSeconds int64

	intervalSeconds = 1
	timeoutSeconds = 600
	updatePeriodSeconds = 1

	env := make([]corev1.EnvVar, 5)

	cfg := &v1.DeploymentConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: dc.Name,
			Labels: map[string]string{
				"app": dc.Name,
			},
		},
		Spec: v1.DeploymentConfigSpec{
			Replicas: 1,

			Selector: map[string]string{
				"app": dc.Name,
				"deploymentconfig": dc.Name,
			},

			Strategy: v1.DeploymentStrategy{
				Type: v1.DeploymentStrategyTypeRolling,

				RollingParams: &v1.RollingDeploymentStrategyParams{
					IntervalSeconds: &intervalSeconds,
					TimeoutSeconds: &timeoutSeconds,
					UpdatePeriodSeconds: &updatePeriodSeconds,
					MaxSurge: &intstr.IntOrString{
						StrVal: "25%",
					},
					MaxUnavailable: &intstr.IntOrString{
						StrVal: "25%",
					},
				},
			},

			Template: &corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Env: env,
						},
					},
				},
			},
		},
	}

	dc.Interface.Create(cfg)
}

func NewDeploymentConfig(name, namespace string) (*DeploymentConfig, error)  {

	client, err := appsv1.NewForConfig(k8s.Config)
	if err != nil {
		return nil, err
	}
	return &DeploymentConfig{
		Name: name,
		Namespace: namespace,

		Interface: client.DeploymentConfigs(namespace),
		Client: client,
	}, nil
}

func TestDeploymentConfigCrud(t *testing.T)  {
	// put below configs in yaml file
	project := "demo"
	profile := "dev"
	app := "hello-world"
	namespace := project + "-" + profile

	// new dc instance
	dc, err := NewDeploymentConfig(app, namespace)
	assert.Equal(t, nil, err)
	assert.Equal(t, app, dc.Name)

	// create dc
	dc.Create()

	// deploy

}