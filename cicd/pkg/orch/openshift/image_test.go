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
)

func TestImageStreamCrud(t *testing.T) {
	imageStreamName := "is-test"
	namespace := "openshift"

	imageStream := &ImageStream{
		Name: imageStreamName,
		Namespace: namespace,
	}

	// create imagestream
	is, err := imageStream.Create()
	assert.Equal(t, nil, err)
	assert.Equal(t, imageStreamName, is.ObjectMeta.Name)

	// get imagestream
	is, err = imageStream.Get()
	assert.Equal(t, nil, err)
	assert.Equal(t, imageStreamName, is.ObjectMeta.Name)

	// delete imagestream
	err = imageStream.Delete()
	assert.Equal(t, nil, err)
}


func TestImageStreamCreation(t *testing.T) {
	name := "s2i-java-test"
	namespace := "openshift"
	source := "docker.io/hidevops/s2i-java:latest"

	imageStream, err := NewImageStreamFromSource(name, namespace, source)

	// create imagestream
	is, err := imageStream.Create()
	assert.Equal(t, nil, err)
	assert.Equal(t, name, is.ObjectMeta.Name)
}


