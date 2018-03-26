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


package auth

import (
	"testing"
	"github.com/hidevopsio/hi/cicd/pkg/scm/factories"
	"github.com/stretchr/testify/assert"
	"os"
	"github.com/hidevopsio/hi/boot/pkg/log"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestUserGetSession(t *testing.T) {
	baseUrl :=  os.Getenv("SCM_URL") + os.Getenv("SCM_API_VER")
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")

	scmFactory := new(factories.ScmFactory)
	scm, err := scmFactory.New(factories.GitlabScmType)
	assert.Equal(t, nil, err)

	scm.GetSession(baseUrl, username, password)
	log.Debug(scm)
}
