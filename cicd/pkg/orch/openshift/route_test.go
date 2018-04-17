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

func TestRouteCreation(t *testing.T)  {
	projectName := "demo"
	profile := "dev"
	namespace := projectName + "-" + profile
	app := "hello-world"

	route, err := NewRoute(app, namespace)
	assert.Equal(t, nil, err)

	err = route.Create(8080)
	assert.Equal(t, nil, err)
}

func TestRouteDeletion(t *testing.T)  {
	projectName := "demo"
	profile := "dev"
	namespace := projectName + "-" + profile
	app := "hello-world"

	route, err := NewRoute(app, namespace)
	assert.Equal(t, nil, err)

	err = route.Delete()
	assert.Equal(t, nil, err)
}