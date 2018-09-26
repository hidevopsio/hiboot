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

package depends

import (
	"fmt"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
)

type ByDependency []interface{}

func (s ByDependency) Len() int {
	return len(s)
}

func (s ByDependency) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByDependency) Less(i, j int) bool {
	item := s[j]
	// find the index of its dependency
	depPath := make([]string, 0)
	depPath, err := s.findDependencies(item, depPath)
	if err != nil {
		// report cycling dependencies
		log.Error(err)
	}
	if len(depPath) > 1 {
		return true
	}
	return false
}

func (s ByDependency) findDependencyIndex(depName string) int {
	for i, item := range s {
		name, _ := reflector.GetName(item)
		if name == depName {
			return i
		}

		pkgName := reflector.ParseObjectPkgName(item)
		fullName := pkgName + "." + name
		if fullName == depName || pkgName == depName {
			return i
		}
	}
	return -1
}

func (s ByDependency) findDependencies(item interface{}, p []string) (path []string, err error) {
	var name string
	path = p
	name, err = reflector.GetName(item)
	if err == nil && len(path) == 0 {
		path = append(path, name)
	}
	depTag, ok := reflector.FindEmbeddedFieldTag(item, "depends")
	if ok {
		if str.InSlice(depTag, path) {
			path = append(path, depTag)
			err = fmt.Errorf("cycling dependencies: %v", path)
			return
		}
		path = append(path, depTag)
		depIdx := s.findDependencyIndex(depTag)
		if depIdx > 0 {
			return s.findDependencies(s[depIdx], path)
		}
	}
	return
}
