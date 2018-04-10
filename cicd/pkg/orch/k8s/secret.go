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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/typed/core/v1"
	"github.com/hidevopsio/hi/boot/pkg/log"
	"k8s.io/apimachinery/pkg/api/errors"
	"strings"
)

type Secret struct {
	Name      string
	Username  string
	Password  string
	Namespace string

	secrets *corev1.Secret

	Interface v1.SecretInterface
}

// Create new instance of type Secret
func NewSecret(name, username, password, namespace string, isToken bool) (*Secret) {
	log.Debugf("NewSecret(%v, %v, %v)", username, strings.Repeat("*", len(password)), namespace)

	var s *Secret
	s = &Secret{
		Name:      name,
		Password:  password,
		Namespace: namespace,
		Interface: ClientSet.CoreV1().Secrets(namespace),
	}
	if !isToken {
		s.Username = username
	}

	return s
}

// Create takes the representation of a secret and creates it.  Returns the server's representation of the secret, and an error, if there is any.
func (s *Secret) Create() error {
	log.Debug("Secret.Create()")
	var data map[string][]byte
	if s.Username != "" {
		data = map[string][]byte{
			corev1.BasicAuthUsernameKey: []byte(s.Username),
			corev1.BasicAuthPasswordKey: []byte(s.Password),
		}
	} else {
		data = map[string][]byte{
			corev1.BasicAuthPasswordKey: []byte(s.Password),
		}
	}
	// k8s.io/api/core/v1/types.go
	coreSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: s.Name,
		},
		Data: data,
		Type: corev1.SecretTypeBasicAuth,
	}
	var err error

	_, err = s.Get()
	if errors.IsNotFound(err) {
		s.secrets, err = s.Interface.Create(coreSecret)
	} else {
		s.secrets, err = s.Interface.Update(coreSecret)
	}

	return err
}

// Get takes name of the secret, and returns the corresponding secret object, and an error if there is any.
func (s *Secret) Get() (*corev1.Secret, error) {
	log.Debug("Secret.Get()")
	var err error
	s.secrets, err = s.Interface.Get(s.Name, metav1.GetOptions{})

	return s.secrets, err
}

// Delete takes name of the secret and deletes it. Returns an error if one occurs.
func (s *Secret) Delete() error {
	log.Debug("Secret.Delete()")
	var err error
	err = s.Interface.Delete(s.Name, &metav1.DeleteOptions{})

	return err
}
