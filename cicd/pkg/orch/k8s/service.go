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
	svcv1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"fmt"
	"github.com/hidevopsio/hi/boot/pkg/log"
	"github.com/jinzhu/copier"
)

type Service struct{
	Name      string
	Namespace string
	Interface  svcv1.ServiceInterface
}

func NewService(name, namespace string) *Service {

	return &Service{
		Name: name,
		Namespace: namespace,
		Interface: ClientSet.CoreV1().Services(namespace),
	}
}

func (s *Service) Create(ports interface{}) error {

	p := make([]corev1.ServicePort, 0)
	copier.Copy(&p, ports)

	// create service
	serviceSpec := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: s.Name,
			Labels: map[string]string{
				"app": s.Name,
			},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: p,
			Selector: map[string]string{
				"app": s.Name,
			},
		},
	}

	svc, err := s.Interface.Get(s.Name, metav1.GetOptions{})
	switch {
	case err == nil:
		serviceSpec.ObjectMeta.ResourceVersion = svc.ObjectMeta.ResourceVersion
		serviceSpec.Spec.ClusterIP = svc.Spec.ClusterIP
		_, err = s.Interface.Update(serviceSpec)
		if err != nil {
			return fmt.Errorf("failed to update service: %s", err)
		}
		log.Info("service updated")
	case errors.IsNotFound(err):
		_, err = s.Interface.Create(serviceSpec)
		if err != nil {
			return fmt.Errorf("failed to create service")
		}
		log.Info("service created")
	default:
		return fmt.Errorf("upexected error: %s", err)
	}
	return nil
}

func (s *Service) Delete() error {
	return s.Interface.Delete(s.Name, &metav1.DeleteOptions{})
}
