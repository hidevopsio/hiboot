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
	"github.com/openshift/api/image/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	imagev1 "github.com/openshift/client-go/image/clientset/versioned/typed/image/v1"
	"github.com/hidevopsio/hi/cicd/pkg/orch/k8s"
	"github.com/hidevopsio/hi/boot/pkg/log"
	"k8s.io/apimachinery/pkg/api/errors"
)

type ImageStreamInterface interface {
	Create() (*v1.ImageStream, error)
	Get() (*v1.ImageStream, error)
	Delete() error
}

type ImageStream struct{
	Name string
	Namespace string
	Source string

	Interface imagev1.ImageStreamInterface
}

func NewImageStream(name, namespace string) (*ImageStream, error) {
	clientSet, err := imagev1.NewForConfig(k8s.Config)

	return &ImageStream{
		Name:      name,
		Namespace: namespace,
		Interface: clientSet.ImageStreams(namespace),
	}, err
}

func NewImageStreamFromSource(name, namespace, source string) (*ImageStream, error) {
	is, err := NewImageStream(name, namespace)
	is.Source = source
	return is, err
}

// @Title Create
// @Description create imagestream
// @Param
// @Return v1.ImageStream, error
func (is *ImageStream) Create() (*v1.ImageStream, error) {
	log.Debug("ImageStream.Create()")

	imageStream := &v1.ImageStream{
		ObjectMeta: metav1.ObjectMeta{
			Name: is.Name,
			Namespace: is.Namespace,
			Labels: map[string]string{
				"app": is.Name,
			},
		},
	}

	if is.Source != "" {
		imageStream.Spec = v1.ImageStreamSpec{
			Tags: []v1.TagReference{
				{
					Name: "latest",
					From: &corev1.ObjectReference{
						Kind: "DockerImage",
						Name: is.Source,
					},
				},
			},
		}
	}

	result, err := is.Get()
	message := "create ImageStream"
	switch {
	case err == nil:
		imageStream.ObjectMeta.ResourceVersion = result.ResourceVersion
		result, err = is.Interface.Update(imageStream)
		message = "update ImageStream"

	case errors.IsNotFound(err):
		result, err = is.Interface.Create(imageStream)
	}

	if err != nil {
		log.Errorf("Failed to %v %v.", message, result.Name)
		return nil, err
	}
	log.Infof("Succeed to %v %v.", message, result.Name)
	return result, err
}

// @Title Get
// @Description get imagestream
// @Param
// @Return v1.ImageStream, error
func (is *ImageStream) Get() (*v1.ImageStream, error) {
	log.Debug("ImageStream.Get()")
	return is.Interface.Get(is.Name, metav1.GetOptions{})
}


// @Title Delete
// @Description delete imagestream
// @Param
// @Return error
func (is *ImageStream) Delete() error {
	log.Debug("ImageStream.Delete()")
	return is.Interface.Delete(is.Name, &metav1.DeleteOptions{})
}
