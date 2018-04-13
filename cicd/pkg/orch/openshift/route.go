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
	routev1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/openshift/api/route/v1"
	"github.com/hidevopsio/hi/cicd/pkg/orch/k8s"
	"github.com/hidevopsio/hi/boot/pkg/log"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/api/errors"
	"fmt"
)



type Route struct{
	Name string
	Namespace string

	Interface routev1.RouteInterface
}

func NewRoute(name, namespace string) (*Route, error)  {
	log.Debug("NewRoute()")
	clientSet, err := routev1.NewForConfig(k8s.Config)
	if err != nil {
		return nil, err
	}
	return &Route{
		Name:      name,
		Namespace: namespace,
		Interface: clientSet.Routes(namespace),
	}, nil
}

func (r *Route) Create(port int32) error {
	log.Debug("Route.Create()")
	cfg := &v1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Name: r.Name,
			Labels: map[string]string{
				"app": r.Name,
			},
		},
		Spec: v1.RouteSpec{
			To: v1.RouteTargetReference{
				Kind: "Service",
				Name: r.Name,
			},
			Port: &v1.RoutePort{
				TargetPort: intstr.IntOrString{
					IntVal: port,
				},
			},
		},
	}

	result, err := r.Interface.Get(r.Name, metav1.GetOptions{})
	switch {
	case err == nil:
		cfg.ObjectMeta.ResourceVersion = result.ResourceVersion
		result, err = r.Interface.Update(cfg)
		if err == nil {
			log.Infof("Updated Route %v", result.Name)
		} else {
			return err
		}
		break
	case errors.IsNotFound(err):
		route, err := r.Interface.Create(cfg)
		if err != nil {
			return err
		}
		log.Infof("Created Route %q.\n", route.Name)
		break
	default:
		return fmt.Errorf("failed to create Route: %s", err)
	}
	return nil
}

func (r *Route) Get() (*v1.Route, error) {
	log.Debug("Route.Delete()")

	return r.Interface.Get(r.Name, metav1.GetOptions{})
}

func (r *Route) Delete() error {
	log.Debug("Route.Delete()")

	return r.Interface.Delete(r.Name, &metav1.DeleteOptions{})
}

