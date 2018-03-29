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
	"github.com/hidevopsio/hi/boot/pkg/log"
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
)

func init()  {
	log.SetLevel(log.DebugLevel)
}


func TestSecretCreate(t *testing.T) {
	log.Debug("TestSecretCrud()")
	username := os.Getenv("SCM_USERNAME")
	password :=  os.Getenv("SCM_PASSWORD")
	secretName := username + "-secret"
	namespace := "demo-dev"

	secret := NewSecret(secretName, username, password, namespace, false)

	// Create secret
	err := secret.Create()
	assert.Equal(t, nil, err)
}

func TestSecretCrud(t *testing.T) {
	log.Debug("TestSecretCrud()")

	secretName := "the-test-secret"
	username := "test"
	password := "test-pwd"
	namespace := "demo-dev"

	secret := NewSecret(secretName, username, password, namespace, false)

	// Create secret
	err := secret.Create()
	assert.Equal(t, nil, err)

	// Get secret
	s, err := secret.Get()
	assert.Equal(t, nil, err)
	assert.Equal(t, s.Name, secretName)

	// Delete secret
	err = secret.Delete()
	assert.Equal(t, nil, err)

}