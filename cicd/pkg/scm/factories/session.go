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

package factories


import (
	"github.com/hidevopsio/hi/boot/pkg/log"
	"github.com/hidevopsio/hi/cicd/pkg/scm"
	"github.com/hidevopsio/hi/cicd/pkg/scm/gitlab"
	"errors"
	"fmt"
)

const (
	GithubScmType = 1
	GitlabScmType = 2
)

type ScmFactory struct{}


func (s *ScmFactory) New(provider int) (scm.SessionInterface, error)  {
	log.Debug("scm.NewSession()")
	switch provider {
	case GithubScmType:
		return nil, errors.New(fmt.Sprintf("SCM of type %d not implemented\n", provider))
	case GitlabScmType:
		return new(gitlab.Session), nil
	default:
		return nil, errors.New(fmt.Sprintf("SCM of type %d not recognized\n", provider))
	}
	return nil, nil
}