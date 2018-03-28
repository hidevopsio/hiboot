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


package k8s

import (
	"k8s.io/apimachinery/pkg/api/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"fmt"
	"github.com/hidevopsio/hi/boot/pkg/log"
)

type Service struct{
	App string
	Project string
}

func (s *Service) Create() error {
	// create service

	serviceSpec := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: s.App,
			Labels: map[string]string{
				"app": s.App,
			},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name: "http",
					Port: 8080,
				},
			},
			Selector: map[string]string{
				"app": s.App,
			},
		},
	}

	services := ClientSet.CoreV1().Services(s.Project)
	svc, err := services.Get(s.App, metav1.GetOptions{})
	switch {
	case err == nil:
		serviceSpec.ObjectMeta.ResourceVersion = svc.ObjectMeta.ResourceVersion
		serviceSpec.Spec.ClusterIP = svc.Spec.ClusterIP
		_, err = services.Update(serviceSpec)
		if err != nil {
			return fmt.Errorf("failed to update service: %s", err)
		}
		log.Info("service updated")
	case errors.IsNotFound(err):
		_, err = services.Create(serviceSpec)
		if err != nil {
			return fmt.Errorf("failed to create service")
		}
		log.Info("service created")
	default:
		return fmt.Errorf("upexected error: %s", err)
	}
	return nil
}
